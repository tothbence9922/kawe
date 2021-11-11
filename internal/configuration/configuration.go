package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ServiceConfigs []ServiceConfiguration
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

type ServiceConfiguration struct {
	Name        string
	PingConfigs []PingConfiguration
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
}

func (pc PingConfiguration) String() string {

	return fmt.Sprintf("periodicity\t\ttarget\n%d\t\t%s\n", pc.Periodicity, pc.Target)
}

func check(e error) {

	if e != nil {
		panic(e)
	}
}
