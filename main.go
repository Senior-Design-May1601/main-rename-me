package main

import (
    "flag"
    "log"
    "sync"

    "github.com/BurntSushi/toml"
)

type pluginConfig struct {
    Path string
    Port int
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
    log.Println("Created plugin manager.")
    // start plugin subprocesses
    err := startPlugins(configs)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Plugin subprocesses started.")
    // TODO: something other than wait here?
    wg.Wait()
}
