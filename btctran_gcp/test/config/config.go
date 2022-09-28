package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Version     string `yaml:"version"`
	WIFPrivKey  string `yaml:"wifPrivKey""`
	Blockcypher struct {
		Token   string `yaml:"token"`
		Coin    string `yaml:"coin"`
		Chain   string `yaml:"chain"`
		FeeRate int    `yaml:"feeRate"`
	} `yaml:"blockcypher"`
}

// ParseConfig from config.yml
func ParseConfig() Config {
	c := Config{}

	data, _err := ioutil.ReadFile("./config/config.yml")
	//data, _err := ioutil.ReadFile("./config/config.yml")
	if _err != nil {
		log.Printf("config.Get err   #%v ", _err)
	}

	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return c
}
