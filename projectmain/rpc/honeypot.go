package rpc

import (
    "github.com/Senior-Design-May1601/projectmain/projectmain"
    "net/rpc"
)

// TODO: Implement projectmain.Honeypot over RPC
type honeypot struct {
    client *rpc.Client
}

// wrap a projectmain.Honeypot impl and make it exportable as part of an
// RPC server
type HoneypotServer struct {
    honeypot projectmain.Honeypot
}

func (h *honeypot) Start() error {
    err := h.client.Call("Honeypot.Start", ...)
    if err != nil {
        panic("Error: " + err.Error())
    }

    return nil
}

func (h *honeypot) Stop() error {
    err := h.client.Call("Honeypot.Stop", ...)
    if err != nil {
        panic("Error: " + err.Error())
    }

    return nil
}

func (h *honeypot) Restart() error {
    err := h.client.Call("Honeypot.Restart", ...)
    if err != nil {
        panic("Error: " + err.Error())
    }

    return nil
}

func (h *HoneypotServer) Start() error {
    err := h.honeypot.Start()
    if err != nil {
        panic("Run error: " + err.Error())
    }

    return nil
}

func (h *HoneypotServer) Stop() error {
    err := h.honeypot.Stop()
    if err != nil {
        panic("Run error: " + err.Error())
    }

    return nil
}

func (h *HoneypotServer) Restart() error {
    err := h.honeypot.Restart()
    if err != nil {
        panic("Run error: " + err.Error())
    }

    return nil
}
