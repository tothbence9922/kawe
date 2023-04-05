package utils

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
)

var (
	kubeconfig *string
)

func initFlags() {
	if pwd, err := os.Getwd(); err == nil && pwd != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(pwd, "kubeConfig.yaml"), "(optional) absolute path to the kubeconfig file")
	} else {
		panic("Failed to define kubeconfig flag.")
	}
}

func GetClientSet() *kubernetes.Clientset {
	var clientSet *kubernetes.Clientset

	if kubeconfig == nil {
		initFlags()
	}

	// Authentication - from outside of the cluster
	if pwd, err := os.Getwd(); err == nil && pwd != "" {
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
	} else {
		// Authentication - from inside of the cluster
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
	}
	return clientSet
}

func GetAnnotationsForEndpointByName(name string, services *v1.ServiceList) map[string]string {
	for _, service := range services.Items {
		if service.ObjectMeta.Name == name {
			return service.Annotations
		}
	}
	return make(map[string]string)
}

func GetNameSpaceConfigs(clientSet *kubernetes.Clientset) []configTypes.NamespaceConfiguration {
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	namespaceConfigs := []configTypes.NamespaceConfiguration{}

	// Iterating through all namespaces, adding services and pods behind services to the configuration
	for _, namespace := range namespaces.Items {
		if err != nil {
			panic(err.Error())
		}

		services, err := clientSet.CoreV1().Services(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		serviceConfigurations := []configTypes.ServiceConfiguration{}
		for _, service := range services.Items {
			podConfigs := []configTypes.PodConfiguration{}

			serviceName := service.Name
			serviceAnnotations := service.ObjectMeta.Annotations
			labelSet := labels.Set(service.Spec.Selector)

			if pods, err := clientSet.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSet.AsSelector().String()}); err == nil {
				for _, pod := range pods.Items {
					podIp := pod.Status.PodIP
					podName := pod.ObjectMeta.Name

					for _, port := range service.Spec.Ports {
						podConfigs = append(podConfigs, configTypes.PodConfiguration{Address: podIp, Port: strconv.Itoa(port.TargetPort.IntValue()), Name: podName, Enabled: true, Periodicity: 5, Timeout: 5000})
					}
				}
			}
			serviceConfigurations = append(serviceConfigurations, configTypes.ServiceConfiguration{Name: serviceName, Annotations: serviceAnnotations, Pods: podConfigs})

		}
		namespaceConfigs = append(namespaceConfigs, configTypes.NamespaceConfiguration{Name: namespace.ObjectMeta.Name, Services: serviceConfigurations})
	}

	return namespaceConfigs
}

func GetServerConfigurations() []configTypes.ServerConfiguration {
	var serverConfigs []configTypes.ServerConfiguration

	httpPortString := os.Getenv("KAWE_HTTP_PORT")
	if httpPortString != "" {
		httpPort, err := strconv.Atoi(httpPortString)
		if err == nil {
			serverConfigs = append(serverConfigs, configTypes.ServerConfiguration{Type: "HTTP", Port: httpPort})
		}
	} else {
		fmt.Println("KAWE_HTTP_PORT not found, using default port")
		serverConfigs = append(serverConfigs, configTypes.ServerConfiguration{Type: "HTTP", Port: 80})
	}

	prometheusPortString := os.Getenv("KAWE_PROMETHEUS_PORT")
	if prometheusPortString != "" {
		prometheusPort, err := strconv.Atoi(prometheusPortString)
		if err == nil {
			serverConfigs = append(serverConfigs, configTypes.ServerConfiguration{Type: "PROMETHEUS", Port: prometheusPort})
		}
	} else {
		fmt.Println("KAWE_PROMETHEUS_PORT not found, using default port")
		serverConfigs = append(serverConfigs, configTypes.ServerConfiguration{Type: "PROMETHEUS", Port: 80})
	}

	return serverConfigs
}
