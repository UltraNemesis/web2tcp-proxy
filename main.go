// main.go
package main

import (
	"log"
	"github.com/UltraNemesis/web2tcp-proxy/web2tcp"

	"github.com/spf13/viper"
)

var conf web2tcp.Configuration

func main() {
	loadConfig(".", "config", &conf)

	server := web2tcp.NewServer(conf)

	log.Println("Starting web2tcp services...")

	server.Start()
}

func loadConfig(configPath string, configName string, conf interface{}) {

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	viper.ReadInConfig()

	viper.Unmarshal(conf)
}
