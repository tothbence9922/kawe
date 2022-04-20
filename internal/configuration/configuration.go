package configuration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Configuration struct {
	EndpointConfigs EndpointConfiguration
	ServerConfigs   []ServerConfiguration
}

var configInstance *Configuration

func GetInstance() *Configuration {

	if configInstance == nil {
		configInstance = new(Configuration)
		configInstance.GetConfiguration()
	}
	return configInstance
}

func (cfg *Configuration) GetConfiguration() {

	pwd, _ := os.Getwd()
	dat, err := os.ReadFile(pwd + "/mnt/config.json")

	if err == nil {
		json.Unmarshal([]byte(dat), &cfg)
	}

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

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// Authentication done, client is available

	// Querying the required resources
	// 	1) namespaces
	//  2) for each namespace, the services and the pods
	//  3) for each pod, the annotations so we can see if it is behind a service
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	namespaceConfigs := []NamespaceConfiguration{}

	// Iterating through all namespaces, adding services and pods behind services to the configuration
	for _, namespace := range namespaces.Items {
		// Getting all pods - including the ones behind services and separate ones aswell.
		allPods, err := clientset.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
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
		endpoints, err := clientset.CoreV1().Endpoints(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

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

	newCfg := Configuration{ServerConfigs: cfg.ServerConfigs, EndpointConfigs: EndpointConfiguration{Namespaces: namespaceConfigs}}

	file, err := json.MarshalIndent(newCfg, "", " ")

	pwd, _ = os.Getwd()
	err = os.WriteFile(pwd+"/mnt/config.json", file, 0644)

	json.Unmarshal([]byte(file), &cfg)
}

// TODO fix
func (c Configuration) String() string {

	var ret string

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

	//	{
	//       "Name": "calico-kube-controllers-748bcb7bb-pk565",
	//       "Address": "10.42.219.56",
	//       "Enabled": true,
	//       "Periodicity": 5,
	//       "Timeout": 5000
	//      }
	enabledString := "DISABLED"
	if p.Enabled {
		enabledString = "ENABLED"
	}
	ret := p.Name + " " + p.Address + " " + enabledString + " " + string(p.Periodicity) + " " + string(p.Timeout)
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
