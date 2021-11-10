package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (cfg *Configuration) GetConfiguration() { // read config from json in
	pwd, _ := os.Getwd()
	dat, err := os.ReadFile(pwd + "/mnt/config.json")
	check(err)

	json.Unmarshal([]byte(dat), &cfg)
	for _, config := range cfg.ServiceConfigs {
		fmt.Print(config.String())
	}
	//fmt.Print(string(dat))
}

var configInstance *Configuration

func GetInstance() *Configuration {
	if configInstance == nil {
		configInstance = new(Configuration)
		configInstance.GetConfiguration()
	}
	return configInstance
}

type ServiceConfiguration struct {
	PingConfigs []PingConfiguration
}
type Configuration struct {
	ServiceConfigs []ServiceConfiguration
}

type PingConfiguration struct {
	Periodicity int
	Method      interface{}
	Target      string
}

func (pc PingConfiguration) String() string {
	return fmt.Sprintf("periodicity\t\ttarget\n%d\t\t%s\n", pc.Periodicity, pc.Target)
}

func (sc ServiceConfiguration) String() string {
	ret := ""
	for _, curPingConfig := range sc.PingConfigs {
		ret += curPingConfig.String()
	}
	return ret
}

func (c Configuration) String() string {
	var ret string

	for _, curSvcConfig := range c.ServiceConfigs {
		ret = ret + curSvcConfig.String()
	}

	return ret

}
