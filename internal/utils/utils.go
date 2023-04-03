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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
)

func GetClientSet() *kubernetes.Clientset {
	var clientSet *kubernetes.Clientset

	// Authentication - from outside of the cluster
	if pwd, err := os.Getwd(); err == nil && pwd != "" {
		kubeconfig := flag.String("kubeconfig", filepath.Join(pwd, "kubeConfig.yaml"), "(optional) absolute path to the kubeconfig file")
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

		// Getting all pods - including the ones behind services and separate ones aswell.
		allPods, err := clientSet.CoreV1().Pods(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		podConfigs := []configTypes.PodConfiguration{}

		for _, pod := range allPods.Items {
			podIp := pod.Status.PodIP
			podName := pod.ObjectMeta.Name
			podLabels := pod.ObjectMeta.Labels
			podAnnotations := pod.GetAnnotations()
			// TODO use target port from service.Spec.Ports
			podConfigs = append(podConfigs, configTypes.PodConfiguration{Address: podIp, Name: podName, Labels: podLabels, Annotations: podAnnotations, Enabled: true, Periodicity: 5, Timeout: 5000})
		}

		// Using endpoints as "services" to see IP addresses
		endpoints, err := clientSet.CoreV1().Endpoints(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
		services, err := clientSet.CoreV1().Services(namespace.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		serviceConfigurations := []configTypes.ServiceConfiguration{}

		for _, endpoint := range endpoints.Items {

			serviceName := endpoint.ObjectMeta.Name
			serviceAnnotations := GetAnnotationsForEndpointByName(endpoint.ObjectMeta.Name, services)
			currentPodConfigs := []configTypes.PodConfiguration{}

			for _, subset := range endpoint.Subsets {
				for _, address := range subset.Addresses {
					for _, podConfig := range podConfigs {
						if address.IP == podConfig.Address {
							currentPodConfigs = append(currentPodConfigs, podConfig)
						}
					}
				}
			}

			serviceConfigurations = append(serviceConfigurations, configTypes.ServiceConfiguration{Name: serviceName, Annotations: serviceAnnotations, Pods: currentPodConfigs})

		}

		separatePodConfigs := []configTypes.PodConfiguration{}

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

		namespaceConfigs = append(namespaceConfigs, configTypes.NamespaceConfiguration{Name: namespace.ObjectMeta.Name, Services: serviceConfigurations, Pods: separatePodConfigs})
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
