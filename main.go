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
	MasterConfig masterConfig
	PluginConfig []pluginConfig
	LoggerConfig []loggerConfig
}

type PluginConfig interface {
	GetName() string
	GetPath() string
	GetArgs() []string
}

type masterConfig struct {
	Logfile string
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
var logfile *os.File

func errExit(reason string) {
	log.Println("Fatal error:", reason)
	err := pluginManager.StopPlugins()
	if err != nil {
		log.Println("Error stopping plugins:", err)
	}
	err = logManager.StopLoggers()
	if err != nil {
		log.Println("Error stopping loggers:", err)
	}
	logfile.Close()
	os.Exit(1)
}

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
		errExit(err.Error())
	}

	// our internal logfile
	logfile, err := os.OpenFile(config.MasterConfig.Logfile,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644)
	if err != nil {
		errExit(err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)

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
	err = logManager.StartLoggers()
	if err != nil {
		errExit(err.Error())
	}

	pluginManager = NewPluginManager(pluginConfigs)
	err = pluginManager.StartPlugins()
	if err != nil {
		errExit(err.Error())
	}

	for signal := range signalHandler {
        log.Println(signal, "received. Shutting down.")

        err = pluginManager.StopPlugins()
        if err != nil {
            log.Println("Error stopping plugins:", err)
        }

        err = logManager.StopLoggers()
        if err != nil {
            log.Println("Error stopping plugins:", err)
        }

        os.Exit(0)
    }
}
