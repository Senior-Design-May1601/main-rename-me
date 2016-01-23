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

// TODO: support plugin arguments
type pluginConfig struct {
    Name string
    Path string
}

type loggerConfig struct {
    Name string
    Path string
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

    logManager = NewLogManager(config.LoggerConfig)
    err := logManager.StartLoggers()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Out of log manager!!")

    pluginManager = NewPluginManager(config.PluginConfig)
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
