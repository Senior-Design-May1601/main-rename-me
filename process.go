package main

import (
	"log"
	"os/exec"
	"sync"
	"syscall"
)

type processInfo struct {
	cmd *exec.Cmd
}

type processMap struct {
	sync.RWMutex
	values map[int]*processInfo
}

type ProcessManager struct {
	processes processMap
}

func NewProcessManager(cmds []PluginConfig) *ProcessManager {
	p := ProcessManager{processes: processMap{values: make(map[int]*processInfo)}}
	p.processes.Lock()
	defer p.processes.Unlock()
	for idx, cmd := range cmds {
		pi := &processInfo{exec.Command(cmd.GetPath(), cmd.GetArgs()...)}
		pi.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		// TODO: chroot here too
		p.processes.values[idx] = pi
	}

	return &p
}

func (x *ProcessManager) NumProcs() int {
	x.processes.RLock()
	defer x.processes.RUnlock()
	return len(x.processes.values)
}

func (x *ProcessManager) StartProcesses() error {
	x.processes.RLock()
	defer x.processes.RUnlock()
	for idx, pi := range x.processes.values {
		if err := pi.cmd.Start(); err != nil {
			return err
		}
		go x.MonitorProcess(idx, pi)
		log.Println("Started process:", pi.cmd.Path)
	}
	return nil
}

func (x *ProcessManager) MonitorProcess(idx int, pi *processInfo) {
	err := pi.cmd.Wait()
	if err != nil {
		// TODO: don't do anything if *we* killed the process
		x.processes.Lock()
		delete(x.processes.values, idx)
		x.processes.Unlock()
		errExit(err.Error())
	}
}

func (x *ProcessManager) KillProcesses() error {
	x.processes.Lock()
	defer x.processes.Unlock()
	for _, pi := range x.processes.values {
		// we could potentially fail before the processes are correctly started
		// in that case, do nothing
		if pi.cmd.Process != nil {
			pgid, err := syscall.Getpgid(pi.cmd.Process.Pid)
			if err != nil {
				return err
			}
			syscall.Kill(-pgid, 15)
		}
	}
	return nil
}
