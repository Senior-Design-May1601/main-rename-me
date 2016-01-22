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

    pluginManager = NewPluginManager(config.PluginConfig)
    err = pluginManager.StartPlugins()
    if err != nil {
        log.Fatal(err)
    }
    // TODO: make this signal startup less racy? what happens if we get
    //       a signal after we've started processes, but before we start the
    //       signal handler? :(
	signalHandler := make(chan os.Signal, 2)
	signal.Notify(signalHandler,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

    for signal := range signalHandler {
        if signal == syscall.SIGHUP {
            pluginManager.RestartPlugins()
        } else {
            pluginManager.StopPlugins(signal)
            // TODO: exit with exit status
            os.Exit(0)
        }
    }
}
