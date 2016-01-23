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

type processList struct {
	sync.RWMutex
	values []*processInfo
}

type ProcessManager struct {
	processes processList
}

func NewProcessManager(paths []string) *ProcessManager {
	p := ProcessManager{}
	p.processes.Lock()
	defer p.processes.Unlock()
	for _, path := range paths {
		pi := &processInfo{exec.Command(path)}
		pi.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		// TODO: chroot here too
		p.processes.values = append(p.processes.values, pi)
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
	for _, pi := range x.processes.values {
		if err := pi.cmd.Start(); err != nil {
			return err
		}
		go x.MonitorProcess(pi)
		log.Println("Started process:", pi.cmd.Path)
	}
	return nil
}

func (x *ProcessManager) MonitorProcess(pi *processInfo) {
	err := pi.cmd.Wait()
	if err != nil {
		// TODO:
		//  - remove process from map
		//  - kill all other processes
		//  - die
		//  - don't do anything if *we* killed the process
		log.Fatal("Process crashed:", err)
	}
}

func (x *ProcessManager) RestartProcesses() error {
	// TODO
	//  - send kill (hup?) signal
	//  - wait for processes to exit
	//  - restart processes
	return nil
}

func (x *ProcessManager) KillProcesses() error {
	x.processes.Lock()
	defer x.processes.Unlock()
	for _, pi := range x.processes.values {
		pgid, err := syscall.Getpgid(pi.cmd.Process.Pid)
		if err != nil {
			return err
		}
		syscall.Kill(-pgid, 15)
	}
	return nil
}
