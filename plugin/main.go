package plugin

import (
    "errors"
    "net"
    "net/http"
    "net/rpc"
    "strconv"
)

const (
    CONTROL_PORT_MIN = 10000
    CONTROL_PORT_MAX = 10099
    // TODO: make this importable from projectmain
    CORE_CONTROL_PORT = "localhost:10100"
)

type Server struct {
    Port int
    listener net.Listener
}

func (x *Server) Serve() error {
    client, err := rpc.DialHTTP("tcp", CORE_CONTROL_PORT)
    if err != nil {
        return err
    }
    // TODO: handle errors here somehow?
    // TODO: possible race condition?
    // TODO: handle plugin manager response
    client.Go("PluginManager.Ready", x.Port, struct{}{}, nil)
    return http.Serve(x.listener, nil)
}

func NewPluginServer(plugin Plugin) (*Server, error) {
    // XXX can we handle having multiple "Plugin" names registered?
    rpc.RegisterName("Plugin", plugin)
    rpc.HandleHTTP()
    listener, port, err := getListener()
    if err != nil {
        return nil, err
    }
    return &Server{port, listener}, nil
}

func getListener() (net.Listener, int, error) {
    for port := CONTROL_PORT_MIN; port <= CONTROL_PORT_MAX; port++ {
        l, e := net.Listen("tcp", "localhost:" + strconv.Itoa(port))
        if e == nil {
            return l, port, nil
        }
    }
    return nil, 0, errors.New("No available control ports.")
}

type Plugin interface {
    Start(args *Args, reply *Reply) error
    Stop(args *Args, reply *Reply) error
    Restart(args *Args, reply *Reply) error
}

type Args struct {
    // TODO
}

type Reply struct {
    // TODO
}
