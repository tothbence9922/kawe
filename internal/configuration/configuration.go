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
	ServiceConfigs   []ServiceConfiguration
	KubernetesConfig KubernetesConfiguration
	ServerConfigs    []ServerConfiguration
}

var configInstance *Configuration

func GetInstance() *Configuration {

	if configInstance == nil {
		configInstance = new(Configuration)
		configInstance.GetFileConfiguration()
		configInstance.GetKubernetesConfiguration()
	}
	return configInstance
}

func (cfg *Configuration) GetKubernetesConfiguration() {

	fmt.Println("Kubernetes configuration started...")

	// Authentication - from outside of the cluster
	var kubeconfig *string
	if pwd, _ := os.Getwd(); pwd != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(pwd, "/mnt/kubeConfig"), "(optional) absolute path to the kubeconfig file")
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
	namespaceConfigs := []KubernetesNamespaceConfiguration{}

	// Iterating through all namespaces, adding services and pods behind services to the configuration
	for _, namespace := range namespaces.Items {
		// Getting all pods - including the ones behind services and separate ones aswell.
		allPods, err := clientset.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		podConfigs := []KubernetesPodConfiguartion{}
		for _, pod := range allPods.Items {
			podIp := pod.Status.PodIP
			podName := pod.ObjectMeta.Name
			podConfigs = append(podConfigs, KubernetesPodConfiguartion{Address: podIp, Name: podName, Enabled: true, Periodicity: 5, Timeout: 5000})
		}

		// Using endpoints as "services"
		endpoints, err := clientset.CoreV1().Endpoints(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		serviceConfigurations := []KubernetesServiceConfiguration{}
		for _, endpoint := range endpoints.Items {
			serviceName := endpoint.ObjectMeta.Name
			currentPodConfigs := []KubernetesPodConfiguartion{}

			for _, subset := range endpoint.Subsets {
				for _, address := range subset.Addresses {
					for _, podConfig := range podConfigs {
						if address.IP == podConfig.Address {
							currentPodConfigs = append(currentPodConfigs, podConfig)
						}
					}
				}
			}
			serviceConfigurations = append(serviceConfigurations, KubernetesServiceConfiguration{Name: serviceName, Pods: currentPodConfigs})
		}

		separatePodConfigs := []KubernetesPodConfiguartion{}
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

		namespaceConfigs = append(namespaceConfigs, KubernetesNamespaceConfiguration{Name: namespace.ObjectMeta.Name, Services: serviceConfigurations, Pods: separatePodConfigs})
	}

	namespacesJSON, err := json.Marshal(namespaceConfigs)

	if err == nil {
		fmt.Println(string(namespacesJSON))
	}

	newCfg := Configuration{ServiceConfigs: cfg.ServiceConfigs, ServerConfigs: cfg.ServerConfigs, KubernetesConfig: KubernetesConfiguration{Namespaces: namespaceConfigs}}

	file, err := json.MarshalIndent(newCfg, "", " ")

	pwd, _ := os.Getwd()
	err = os.WriteFile(pwd+"/mnt/config.json", file, 0644)

	fmt.Println("Kubernetes configuration finished!")
}

func (cfg *Configuration) GetFileConfiguration() {

	pwd, _ := os.Getwd()
	dat, err := os.ReadFile(pwd + "/mnt/config.json")
	check(err)

	json.Unmarshal([]byte(dat), &cfg)
}

func (c Configuration) String() string {

	var ret string

	for _, curSvcConfig := range c.ServiceConfigs {
		ret = ret + curSvcConfig.String()
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

type KubernetesConfiguration struct {
	Namespaces []KubernetesNamespaceConfiguration
}

type KubernetesNamespaceConfiguration struct {
	Name     string
	Services []KubernetesServiceConfiguration
	Pods     []KubernetesPodConfiguartion
}

type KubernetesServiceConfiguration struct {
	Name string
	Pods []KubernetesPodConfiguartion
}

type KubernetesPodConfiguartion struct {
	Name        string
	Address     string
	Enabled     bool
	Periodicity int
	Timeout     int
}

type ServiceConfiguration struct {
	Name            string
	PingConfigs     []PingConfiguration
	ProcessorConfig ProcessorConfiguration
}

func (sc ServiceConfiguration) String() string {

	ret := ""
	for _, curPingConfig := range sc.PingConfigs {
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

func check(e error) {

	if e != nil {
		panic(e)
	}
}
