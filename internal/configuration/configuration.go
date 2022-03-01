package configuration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Configuration struct {
	ServiceConfigs []ServiceConfiguration
	ServerConfigs  []ServerConfiguration
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
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))

	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.ObjectMeta.Name)
	}

	pods, err := clientset.CoreV1().Pods("qdak").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Println(pod.ObjectMeta.Name)
	}

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	namespace := "qdak"
	podName := "qdak-vendor-admin-frontend-84d4c9b696-4t2vm"
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", podName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			podName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", podName, namespace)
	}
	services, err := clientset.CoreV1().Services("qdak").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d services in the cluster\n", len(services.Items))

	if len(services.Items) > 0 {
		fmt.Println("%s\n\n", services.Items[0])
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
