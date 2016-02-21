package main

type PluginManager struct {
	manager ProcessManager
}

func NewPluginManager(configs []PluginConfig) *PluginManager {
	//cmds := make([]string, len(configs))
	//for i, v := range configs {
	//		cmds[i] = v.Exec
	//	}

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
