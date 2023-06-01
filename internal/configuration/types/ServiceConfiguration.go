package types

type ServiceConfiguration struct {
	Name            string
	ServiceLabel    string
	Annotations     map[string]string
	Pods            []PodConfiguration
	ProcessorConfig ProcessorConfiguration
}

func (sc ServiceConfiguration) String() string {

	ret := ""
	for _, curPingConfig := range sc.Pods {
		ret += curPingConfig.String()
	}
	return ret
}
