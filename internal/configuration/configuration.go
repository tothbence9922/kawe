package configuration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Configuration struct {
	ServerConfigs   []ServerConfiguration
	EndpointConfigs EndpointConfiguration
}

var configInstance *Configuration

func GetInstance() *Configuration {

	if configInstance == nil {
		configInstance = new(Configuration)
		configInstance.GetConfiguration()
	}
	return configInstance
}

func getClientSet(inCluster bool) *kubernetes.Clientset {
	var clientSet *kubernetes.Clientset

	if inCluster {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()

		if err != nil {
			panic(err.Error())
		}

		// creates the clientSet
		clientSet, err = kubernetes.NewForConfig(config)

		if err != nil {
			panic(err.Error())
		}
	} else {
		// Authentication - from outside of the cluster
		var kubeconfig *string

		if pwd, _ := os.Getwd(); pwd != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(pwd, "/mnt/kubeConfig.yaml"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

		if err != nil {
			panic(err.Error())
		}

		// create the clientSet
		clientSet, err = kubernetes.NewForConfig(config)

		if err != nil {
			panic(err.Error())
		}

	}
	return clientSet
}

func getNameSpaceConfigs(clientSet *kubernetes.Clientset) []NamespaceConfiguration {
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	namespaceConfigs := []NamespaceConfiguration{}

	// Iterating through all namespaces, adding services and pods behind services to the configuration
	for _, namespace := range namespaces.Items {

		// Getting all pods - including the ones behind services and separate ones aswell.
		allPods, err := clientSet.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		podConfigs := []PodConfiguration{}

		for _, pod := range allPods.Items {
			podIp := pod.Status.PodIP
			podName := pod.ObjectMeta.Name
			podConfigs = append(podConfigs, PodConfiguration{Address: podIp, Name: podName, Enabled: true, Periodicity: 5, Timeout: 5000})
		}

		// Using endpoints as "services"
		endpoints, err := clientSet.CoreV1().Endpoints(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		serviceConfigurations := []ServiceConfiguration{}

		for _, endpoint := range endpoints.Items {

			serviceName := endpoint.ObjectMeta.Name
			currentPodConfigs := []PodConfiguration{}

			for _, subset := range endpoint.Subsets {
				for _, address := range subset.Addresses {
					for _, podConfig := range podConfigs {
						if address.IP == podConfig.Address {
							currentPodConfigs = append(currentPodConfigs, podConfig)
						}
					}
				}
			}

			serviceConfigurations = append(serviceConfigurations, ServiceConfiguration{Name: serviceName, Pods: currentPodConfigs})

		}

		separatePodConfigs := []PodConfiguration{}

		for _, podConfig := range podConfigs {
			found := false

			for _, serviceConfig := range serviceConfigurations {
				for _, svcPodConfig := range serviceConfig.Pods {
					if svcPodConfig.Address == podConfig.Address {
						found = true
					}
				}
			}

			if !found {
				separatePodConfigs = append(separatePodConfigs, podConfig)
			}
		}

		namespaceConfigs = append(namespaceConfigs, NamespaceConfiguration{Name: namespace.ObjectMeta.Name, Services: serviceConfigurations, Pods: separatePodConfigs})
	}

	return namespaceConfigs
}

func getServerConfigurations() []ServerConfiguration {
	var serverConfigs []ServerConfiguration

	httpPortString, found := os.LookupEnv("KAWE_HTTP_PORT")
	if found {
		httpPort, err := strconv.Atoi(httpPortString)
		if err == nil {
			serverConfigs = append(serverConfigs, ServerConfiguration{Type: "HTTP", Port: httpPort})
		}
	}

	prometheusPortString, found := os.LookupEnv("KAWE_PROMETHEUS_PORT")
	if found {
		prometheusPort, err := strconv.Atoi(prometheusPortString)
		if err == nil {
			serverConfigs = append(serverConfigs, ServerConfiguration{Type: "PROMETHEUS", Port: prometheusPort})
		}
	}
	return serverConfigs
}

func (cfg *Configuration) GetConfiguration() {
	isInCluster := false
	if len(os.Args) > 1 {
		inCluster := flag.Bool("inCluster", false, "a bool")
		flag.Parse()
		isInCluster = *inCluster
	}

	clientSet := getClientSet(isInCluster)

	GetInstance().EndpointConfigs = EndpointConfiguration{Namespaces: getNameSpaceConfigs(clientSet)}
	GetInstance().ServerConfigs = getServerConfigurations()
}

func (c Configuration) String() string {
	byteFormat, err := json.MarshalIndent(c, "", " ")
	var ret string
	if err != nil {
		return "Failed to marshal configuration."
	} else {
		ret = string(byteFormat)
	}

	return ret
}

type ServerConfiguration struct {
	Type string
	Port int
}

func (sc ServerConfiguration) String() string {

	portString := fmt.Sprintf(":%d", sc.Port)

	return sc.Type + portString
}

type EndpointConfiguration struct {
	Namespaces []NamespaceConfiguration
}

type NamespaceConfiguration struct {
	Name     string
	Services []ServiceConfiguration
	Pods     []PodConfiguration
}

type ServiceConfiguration struct {
	Name            string
	Pods            []PodConfiguration
	ProcessorConfig ProcessorConfiguration
}

type PodConfiguration struct {
	Name        string
	Address     string
	Enabled     bool
	Periodicity int
	Timeout     int
}

func (p PodConfiguration) String() string {
	enabledString := "DISABLED"
	if p.Enabled {
		enabledString = "ENABLED"
	}
	ret := p.Name + " " + p.Address + " " + enabledString + " " + fmt.Sprint(p.Periodicity) + " " + fmt.Sprint(p.Timeout)
	return ret
}

func (sc ServiceConfiguration) String() string {

	ret := ""
	for _, curPingConfig := range sc.Pods {
		ret += curPingConfig.String()
	}
	return ret
}

type PingConfiguration struct {
	Periodicity int
	Method      interface{}
	Target      string
	Timeout     int
}

func (pc PingConfiguration) String() string {

	return fmt.Sprintf("periodicity\t\ttarget\n%d\t\t%s\n", pc.Periodicity, pc.Target)
}

type ProcessorConfiguration struct {
	Type   string
	Params interface{}
}

func (pc ProcessorConfiguration) String() string {

	return fmt.Sprintf("\tType\n%s\n", pc.Type)
}
