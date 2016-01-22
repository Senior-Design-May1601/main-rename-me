package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
    "os"
    "strconv"
    "sync"

    "github.com/Senior-Design-May1601/projectmain/control"
    "github.com/Senior-Design-May1601/projectmain/loggerplugin"
)

type ConnectionKey struct {
    Port int
    Name string
}

type loggerConnectionMap struct {
    sync.RWMutex
    values map[ConnectionKey]*rpc.Client
}

type LogManager struct {
    callChan chan *rpc.Call
    manager ProcessManager
    loggerConnections loggerConnectionMap
    listener net.Listener
}

// TODO: update configs to already have this list
func NewLogManager(configs []loggerConfig) *LogManager {
    paths := make([]string, len(configs))
    for i, v := range configs {
        paths[i] = v.Path
    }

    manager := LogManager{
        callChan: make(chan *rpc.Call, 100),
        manager: *NewProcessManager(paths),
        loggerConnections: loggerConnectionMap{
            values: make(map[ConnectionKey]*rpc.Client),
        },
        listener: nil,
    }

    rpc.Register(&manager)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", control.CONTROL_PORT_CORE)
    if e != nil {
        log.Fatal("listen error:", e)
    }

    manager.listener = l
    go http.Serve(l, nil)
    go manager.handleCallReplies()

    log.Println("Log manager listening on:", control.CONTROL_PORT_CORE)

    return &manager
}

// TODO: timeout and fail if a logger takes too long to get setup
func (x *LogManager) StartLoggers() error {
    err := x.manager.StartProcesses()
    if err != nil {
        return err
    }
    // TODO: somehow wait here until all plugins are both ready and connected

    return nil
}

func (x *LogManager) StopLoggers(signal os.Signal) error {
    return x.manager.KillProcesses(signal)
}

// TODO
func (x *LogManager) RestartLoggers() error {
    return nil
}

// called when a logger is ready to actually start logging
func (x *LogManager) Ready(arg loggerplugin.ReadyArg, _ *int) error {
    log.Println("Ready() called")
    go x.connect(ConnectionKey{arg.Port, arg.Name})
    log.Println("Starting plugin on port", arg.Port)
    return nil
}

func (x *LogManager) Log(p []byte, _ *int) error {
    log.Println("Got log event:", string(p))
    x.loggerConnections.RLock()
    var r int
    for key, client := range x.loggerConnections.values {
        // TODO: handle reply/err
        client.Go(key.Name + ".Log", p, &r, nil)
    }
    x.loggerConnections.RUnlock()

    return nil
}

func (x *LogManager) connect(key ConnectionKey) error {
    client, err := rpc.DialHTTP("tcp", "localhost:" + strconv.Itoa(key.Port))
    if err != nil {
        return err
    }

    x.loggerConnections.Lock()
    x.loggerConnections.values[key] = client
    x.loggerConnections.Unlock()
    log.Println("Logger started on port", key.Port)

    return nil
}

func (x *LogManager) handleCallReplies() {
    for call := range x.callChan {
        if call.Error != nil {
            log.Fatal("Error from client: ", call.Error)
        } else {
            // TODO: log port here
            log.Println("Logger start: OK.")
        }
    }
}
