package configuration

type IConfiguration interface {
	GetFileConfiguration()
	GetKubernetesConfiguration()
	String() string
	GetInstance()
}
