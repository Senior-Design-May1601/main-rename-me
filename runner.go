package main

import (
    "os/exec"
)

func startPlugins(configs PluginConfigs) error {
    for _, config := range configs.PluginConfig {
        err := startPlugin(config.Path)
        if err != nil {
            return err
        }
    }
    return nil
}

func startPlugin(path string) error {
    cmd := exec.Command(path)
    return cmd.Start()
}
