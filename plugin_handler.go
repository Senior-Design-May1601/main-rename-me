package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
    "strconv"
    "sync"

    "github.com/Senior-Design-May1601/projectmain/plugin"
)

type clientMap struct {
    sync.RWMutex
    values map[int]*rpc.Client
}

type PluginManager struct {
    callChan chan *rpc.Call
    clients clientMap
    listener net.Listener
    wg *sync.WaitGroup
}

func (x *PluginManager) Ready(port int, _ *struct{}) error {
    go x.connectAndStart(port)
    log.Println("Starting plugin on port", port)
    return nil
}

func (x *PluginManager) connectAndStart(port int) error {
    err := x.connect(port)
    if err != nil {
        return err
    }
    var reply plugin.Reply
    x.clients.RLock()
    x.clients.values[port].Go(
        "Plugin.Start",
        &plugin.Args{},
        &reply,
        x.callChan)
    x.clients.RUnlock()
    return nil
}

// TODO: check that we don't have collisions in map?
func (x *PluginManager) connect(port int) error {
    x.clients.Lock()
    defer x.clients.Unlock()
    client, err := rpc.DialHTTP("tcp", "localhost:" + strconv.Itoa(port))
    if err != nil {
        return err
    }
    x.clients.values[port] = client
    log.Println("Plugin started on port", port)
    return nil
}

func (x *PluginManager) handleCallReplies() {
    for call := range x.callChan {
        if call.Error != nil {
            log.Fatal("Error from client: ", call.Error)
        } else {
            // TODO: log port here
            log.Println("Client start: OK.")
        }
    }
}

func NewPluginManager(wg *sync.WaitGroup) *PluginManager {
    manager := &PluginManager{
        callChan: make(chan *rpc.Call, 100),
        clients: clientMap{
            values: make(map[int]*rpc.Client),
        },
        listener: nil,
        wg: wg,
    }
    (*wg).Add(1)

    rpc.Register(manager)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", CONTROL_PORT)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    (*manager).listener = l
    go http.Serve(l, nil)
    go (*manager).handleCallReplies()

    return manager
}
