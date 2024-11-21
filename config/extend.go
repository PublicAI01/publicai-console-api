package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	DataHubIp string `yaml:"dataHubIp"`
	Token     string `yaml:"token"`
}

type Demo struct {
	Name string `yaml:"name" json:"name"`
}
