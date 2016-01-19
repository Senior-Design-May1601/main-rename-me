package main

import (
    "flag"
    "log"
    "sync"

    "github.com/BurntSushi/toml"
)

const CONTROL_PORT = "localhost:10100"

type pluginConfig struct {
    Path string
}

type PluginConfigs struct {
    PluginConfig []pluginConfig
}

var pluginManager *PluginManager

func main() {
    var configPath = flag.String("config", "", "projectmain config file")
    flag.Parse()

    var configs PluginConfigs
    if _, err := toml.DecodeFile(*configPath, &configs); err != nil {
        log.Fatal(err)
    }

    var wg sync.WaitGroup
    pluginManager = NewPluginManager(&wg)
    err := startPlugins(configs)
    if err != nil {
        log.Fatal(err)
    }
    // TODO: something other than wait here?
    wg.Wait()
}
