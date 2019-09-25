package main

import (
	"drill/config"
	"drill/realtime"
	"drill/restapi"
	"drill/turn"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "configuration file location")
	flag.Parse()

	if configFile == "" {
		fmt.Println("cannot load file at ", configFile)
		return
	}

	jsonFile, err := os.Open(configFile)
	if err != nil {
		fmt.Println("load config file content fail", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("load config file content fail", err)
		return
	}

	json.Unmarshal(byteValue, &config.GlobalConfig)

	roomBackendChoice := config.GlobalConfig.RoomBackend
	if roomBackendChoice != "memory" && roomBackendChoice != "redis" {
		fmt.Println("room backend only support (memory | redis)")
		return
	}
	if roomBackendChoice == "redis" {
		if config.GlobalConfig.RedisRoomBackendConfig == nil {
			fmt.Println("miss redis config")
			return
		}
	}

	door := realtime.NewWsDoor(config.GlobalConfig.WsPort)
	api := restapi.NewRestApi(config.GlobalConfig.HttpPort)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		door.Start()
	}()
	go func() {
		api.Start()
	}()

	var turnServer *turn.TurnServer
	if config.GlobalConfig.EnableTurn {
		turnServer = turn.NewTurnServer(
			config.GlobalConfig.TurnConfig.User,
			config.GlobalConfig.TurnConfig.Password,
			config.GlobalConfig.TurnConfig.UdpPort,
			config.GlobalConfig.TurnConfig.Realm,
		)
		go func() {
			turnServer.Start()
		}()
	}

	<-c

	if config.GlobalConfig.EnableTurn {
		if turnServer != nil {
			turnServer.Stop()
		}
	}
}
