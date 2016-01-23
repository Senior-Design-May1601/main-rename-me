package main

type PluginManager struct {
	manager ProcessManager
}

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

func (x *PluginManager) RestartPlugins() error {
	// TODO
	return nil
}

func (x *PluginManager) StopPlugins() error {
	return x.manager.KillProcesses()
}
