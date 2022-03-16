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

/*
	- Automatic discovery of cluster-resources (namespaces, services, pods)
		+ Authenticating
		+ GET-ing the resources
	- Creating a configuration entry
		+ Mapping the acquired data into a configuration entry
	- Appending the configuration entry to the already existing (file based) configuration
		+ or creating the whole configuration
*/

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
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))

	// Iterating through all namespaces, adding services and pods behind services to the configuration
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.ObjectMeta.Name)
		services, err := clientset.CoreV1().Services(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Services:")
		for _, service := range services.Items {
			fmt.Println("\t", service.ObjectMeta.Name)
			fmt.Println("\t\t", service.ObjectMeta.String())
			fmt.Println("\t\t", service.TypeMeta.String())
			fmt.Println("\t\t", service.Status.String())
			fmt.Println("\t\t", service.Status.Conditions)
		}

		endpoints, err := clientset.CoreV1().Endpoints(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		fmt.Println("Endpoints:")
		for _, endpoint := range endpoints.Items {
			fmt.Println("\t", endpoint.ObjectMeta.Name)
			for _, subset := range endpoint.Subsets {
				for _, address := range subset.Addresses {
					fmt.Println("\t\t", address.IP)
				}
			}
		}

		pods, err := clientset.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Pods:")
		for _, pod := range pods.Items {
			podIp := pod.Status.PodIP
			fmt.Println("----->", podIp)
			annotations := pod.Annotations

			portsString, ok := annotations["field.cattle.io/ports"]

			if !ok {
				fmt.Println("No ports found")
				continue
			}
			fmt.Println("PORTS STRING: ", portsString)
			portsByteArray := []byte(portsString)

			var ports [][]map[string]interface{}

			if err := json.Unmarshal(portsByteArray, &ports); err != nil {
				panic(err)
			}

			for i := range ports {
				for _, port := range ports[i] {
					fmt.Println("PROTOCOL: ", port["protocol"])
					fmt.Println("PORT: ", port["containerPort"])
				}
			}
		}
	}
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
	Namespaces []NamespaceConfiguartion
}

type NamespaceConfiguartion struct {
	Services []ServiceConfiguartion
	Pods     []PodConfiguartion
}
type ServiceConfiguartion struct {
	Pods []PodConfiguartion
}
type PodConfiguartion struct {
	Name    string
	Address string
	Method  string
	Enabled bool
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
