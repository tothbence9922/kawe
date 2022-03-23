package configuration

type IConfiguration interface {
	GetConfiguration()
	String() string
	GetInstance()
}
