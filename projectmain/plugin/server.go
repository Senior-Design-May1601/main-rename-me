package plugin

import (
    projectmainrpc "github.com/Senior-Design-May1601/projectmain/projectmain/rpc"
)

// wait for a connection to this plugin and return a projectmain RPC server
// that you can use to register components and serve them
func Server() (*projectmainrpc.Server, error) {
    // TODO: verify we are starting this ourselves (e.g. magic cookie)

    listener, err := net.Listen("unix", 
}
