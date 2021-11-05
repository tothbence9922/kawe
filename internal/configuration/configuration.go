package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration2 interface {
	GetConfiguration()
}

type configuration struct {
	PingConfigs []PingConfiguration
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (cfg *configuration) GetConfiguration() { // read config from json in
	pwd, _ := os.Getwd()
	dat, err := os.ReadFile(pwd + "/mnt/config.json")
	check(err)

	json.Unmarshal([]byte(dat), &cfg)
	for _, config := range cfg.PingConfigs {
		fmt.Print(config.String())
	}
	//fmt.Print(string(dat))
}

var configInstance *configuration

func GetInstance() *configuration {
	if configInstance == nil {
		configInstance = new(configuration)
		configInstance.GetConfiguration()
	}
	return configInstance
}

type PingConfiguration struct {
	Periodicity int
	Method      interface{}
	Target      string
}

func (pc PingConfiguration) String() string {
	return fmt.Sprintf("periodicity\t\ttarget\n%d\t\t%s\n", pc.Periodicity, pc.Target)
}

func (c configuration) String() string {
	var ret string

	for _, curPingConfig := range c.PingConfigs {
		ret = ret + curPingConfig.String()
	}

	return ret

}
