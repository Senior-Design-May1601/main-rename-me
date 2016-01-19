package plugin

import (
    "net"
    "net/http"
    "net/rpc"
    "strconv"
)

type Server struct {
    Port int
    listener net.Listener
}

func (x *Server) Serve() error {
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
    if err != nil {
        return err
    }
    client.Go("PluginManager.Ready", x.Port, struct{}{}, nil)
    return http.Serve(x.listener, nil)
}

func NewPlugin(plugin Plugin) (*Server, error) {
    rpc.RegisterName("Plugin", plugin)
    rpc.HandleHTTP()
    listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(plugin.Port()))
    if err != nil {
        return nil, err
    }
    return &Server{plugin.Port(), listener}, nil
}

type Plugin interface {
    // RPC functions
    Start(args *Args, reply *Reply) error
    Stop(args *Args, reply *Reply) error
    Restart(args *Args, reply *Reply) error

    // TODO: should these be RPC functions too?
    Port() int
}

type Args struct {
    Port int `json:"port"`
}

type Reply struct {
    Code int `json:"code"`
    Msg string `json:"msg"`
}
