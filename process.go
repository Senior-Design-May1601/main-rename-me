package main

import (
    "log"
    "os"
    "os/exec"
    "sync"
    "syscall"
)

type ProcessMap struct {
    sync.RWMutex
    values []*exec.Cmd
}

type ProcessManager struct {
    processes ProcessMap
}

func NewProcessManager(paths []string) *ProcessManager {
    p := ProcessManager{}
    p.processes.Lock()
    defer p.processes.Unlock()
    for _, path := range paths {
        p.processes.values = append(p.processes.values, exec.Command(path))
    }

    return &p
}

func (x *ProcessManager) StartProcesses() error {
    x.processes.RLock()
    defer x.processes.RUnlock()
    for _, proc := range x.processes.values {
        if err := (*proc).Start(); err != nil {
            return err
        }
        go x.MonitorProcess(proc)
        // TODO: log to our configured loggers?
        log.Println("Started process:", proc.Path)
    }
    return nil
}

// TODO: make an issue about whether or not we should crash everything here
func (x *ProcessManager) MonitorProcess(p *exec.Cmd) {
    err := (*p).Wait()
    if err != nil {
        x.KillProcesses(syscall.SIGKILL)
        log.Fatal("Plugin crashed:", err)
    }
}

// TODO:
//  - send kill (hup?) signal
//  - wait for processes to exit
//  - restart processes
//
// NOTE: it's fine to block here

func (x *ProcessManager) RestartProcesses() error {
    return nil
}

// TODO:
//  - the "nice" thing to do here would be to send the actual signal we get
//    to the plugins, set a timer, and then send a sigkill if they're not
//    done when the timer runs out
//  - return exit code???
//
// NOTE: fine to block here

func (x *ProcessManager) KillProcesses(signal os.Signal) error {
    x.processes.Lock()
    x.processes.Unlock()
    var wg sync.WaitGroup
    for _, proc := range x.processes.values {
        wg.Add(1)
        go func() {
            // TODO: handle this error somehow?
            _ = (*proc).Wait()
            wg.Done()
        }()
        // TODO: handle this error somehow?
        // TODO: how to check the process status?
        _ = (*proc).Process.Kill()
    }
    wg.Wait()
    return nil
}
