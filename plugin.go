package main

import (
    "os"
)

type PluginManager struct {
    manager ProcessManager
}

// TODO: do this nicer...maybe change structure of config file?
func NewPluginManager(configs []pluginConfig) *PluginManager {
    paths := make([]string, len(configs))
    for i, v := range configs {
        paths[i] = v.Path
    }

    return &PluginManager{*NewProcessManager(paths)}
}

func (x *PluginManager) StartPlugins() error {
    return x.manager.StartProcesses()
}

// TODO:
//  - send kill (hup?) signal
//  - wait for processes to exit
//  - restart processes
//
// NOTE: it's fine to block here

func (x *PluginManager) RestartPlugins() error {
    return nil
}

// TODO:
//  - the "nice" thing to do here would be to send the actual signal we get
//    to the plugins, set a timer, and then send a sigkill if they're not
//    done when the timer runs out
//  - return exit code???
//
// NOTE: fine to block here

func (x *PluginManager) StopPlugins(signal os.Signal) error {
    return x.manager.KillProcesses(signal)
}
