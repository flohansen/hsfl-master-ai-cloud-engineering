package proxyutils

import (
	"github.com/spf13/viper"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/http-proxy-service/proxy"
	"log"
)

func DefaultProxyManagerConfigurationReader(path string) *proxy.Config {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while loading proxy configuration: %v", err)
	}
	viper.AutomaticEnv()

	proxyConfig := &proxy.Config{}

	err = viper.UnmarshalKey("proxy", proxyConfig)
	if err != nil {
		panic(err)
	}

	return proxyConfig
}
