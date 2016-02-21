package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
)

type Config struct {
	PluginConfig []pluginConfig
	LoggerConfig []loggerConfig
}

type PluginConfig interface {
	GetName() string
	GetPath() string
	GetArgs() []string
}

type pluginConfig struct {
	Name string
	Path string
	Args []string
}

type loggerConfig struct {
	Name string
	Path string
	Args []string
}

func (x *pluginConfig) GetName() string {
	return x.Name
}

func (x *pluginConfig) GetPath() string {
	return x.Path
}

func (x *pluginConfig) GetArgs() []string {
	return x.Args
}

func (x *loggerConfig) GetName() string {
	return x.Name
}

func (x *loggerConfig) GetPath() string {
	return x.Path
}

func (x *loggerConfig) GetArgs() []string {
	return x.Args
}

var pluginManager *PluginManager
var logManager *LogManager

func main() {
	signalHandler := make(chan os.Signal, 2)
	signal.Notify(signalHandler,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var configPath = flag.String("config", "", "projectmain config file")
	flag.Parse()

	var config Config
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		log.Fatal(err)
	}

	// this is messy...done for easier toml parsing
	// there's almost certainly a better way to do this
	loggerConfigs := make([]PluginConfig, len(config.LoggerConfig))
	for i := range config.LoggerConfig {
		loggerConfigs[i] = &config.LoggerConfig[i]
	}

	pluginConfigs := make([]PluginConfig, len(config.PluginConfig))
	for i := range config.PluginConfig {
		pluginConfigs[i] = &config.PluginConfig[i]
	}

	logManager = NewLogManager(loggerConfigs)
	err := logManager.StartLoggers()
	if err != nil {
		log.Fatal(err)
	}

	pluginManager = NewPluginManager(pluginConfigs)
	err = pluginManager.StartPlugins()
	if err != nil {
		log.Fatal(err)
	}

	for signal := range signalHandler {
		if signal == syscall.SIGHUP {
			// TODO
			pluginManager.StopPlugins()
			logManager.StopLoggers()
			os.Exit(0)
		} else {
			pluginManager.StopPlugins()
			logManager.StopLoggers()
			// TODO: exit with useful exit status
			os.Exit(0)
		}
	}
}
