package main

import (
	"strings"

	"github.com/AlertFlow/runner/config"
	"github.com/AlertFlow/runner/internal/common"
	payloadendpoints "github.com/AlertFlow/runner/internal/payload_endpoints"
	"github.com/AlertFlow/runner/internal/plugins"
	"github.com/AlertFlow/runner/internal/runner"

	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
)

const version string = "0.21.0"

var (
	configFile = kingpin.Flag("config", "Config File").Short('c').Default("config.yaml").String()
)

func logging(logLevel string) {
	logLevel = strings.ToLower(logLevel)

	if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "warn" {
		log.SetLevel(log.WarnLevel)
	} else if logLevel == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Info("Starting AlertFlow Runner. Version: ", version)

	log.Info("Loading config")
	config, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	logging(config.LogLevel)

	plugins, pluginsMap, actions, payloadEndpoints := plugins.Init()

	common.RegisterActions(actions)
	go payloadendpoints.InitPayloadRouter(config.PayloadEndpoints.Port, plugins, payloadEndpoints)

	runner.RegisterAtAPI(version, pluginsMap, actions, payloadEndpoints)
	go runner.SendHeartbeat()

	Init()

	<-make(chan struct{})
}

func Init() {
	switch strings.ToLower(config.Config.Mode) {
	case "master":
		log.Info("Runner is in Master Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker()
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	case "worker":
		log.Info("Runner is in Worker Mode")
		log.Info("Starting Execution Checker")
		go common.StartWorker()
	case "listener":
		log.Info("Runner is in Listener Mode")
		log.Info("Starting Payload Listener")
		// go payloadhandler.InitPayloadRouter(config.Config.Payloads.Port, config.Config.Payloads.Managers)
	}
}
