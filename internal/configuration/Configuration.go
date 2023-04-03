package configuration

import (
	"encoding/json"
	"time"

	"github.com/go-co-op/gocron"
	"k8s.io/client-go/kubernetes"

	configTypes "github.com/tothbence9922/kawe/internal/configuration/types"
	"github.com/tothbence9922/kawe/internal/utils"
)

type Configuration struct {
	ServerConfigs   []configTypes.ServerConfiguration
	EndpointConfigs configTypes.EndpointConfiguration
	scheduler       *gocron.Scheduler
}

var (
	configInstance *Configuration
	kubeconfig     *string
)

func GetInstance() *Configuration {

	if configInstance == nil {
		configInstance = new(Configuration)
		configInstance.scheduler = gocron.NewScheduler(time.UTC)
		configInstance.GetConfiguration()
	}
	return configInstance
}

func (cfg *Configuration) getTargets(clientSet *kubernetes.Clientset) {
	cfg.EndpointConfigs = configTypes.EndpointConfiguration{Namespaces: utils.GetNameSpaceConfigs(clientSet)}
}

func (cfg *Configuration) GetConfiguration() {

	clientSet := utils.GetClientSet()

	cfg.ServerConfigs = utils.GetServerConfigurations()
	cfg.getTargets(clientSet)
	cfg.scheduler.Every("1m").Do(cfg.getTargets, clientSet)
	cfg.scheduler.StartAsync()
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
