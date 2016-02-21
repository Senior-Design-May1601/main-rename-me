package main

type PluginManager struct {
	manager ProcessManager
}

func NewPluginManager(configs []PluginConfig) *PluginManager {
	return &PluginManager{*NewProcessManager(configs)}
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
