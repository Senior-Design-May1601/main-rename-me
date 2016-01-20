package plugin

import (
    "errors"
    "log"
    "net"
    "net/http"
    "net/rpc"
    "strconv"
    "sync"

    "github.com/Senior-Design-May1601/projectmain/control"
)

type Server struct {
    Port int
    listener net.Listener
    readyChan chan *rpc.Call
}

func (x *Server) Serve() {
    client, err := rpc.DialHTTP("tcp", control.CONTROL_PORT_CORE)
    if err != nil {
        log.Fatal(err)
    }
    // TODO: do this better (wait groups feel weird here)
    var wg sync.WaitGroup
    wg.Add(1)
    // Serve needs to run before the client connection to avoid race condition
    go http.Serve(x.listener, nil)
    go client.Go("PluginManager.Ready", x.Port, struct{}{}, x.readyChan)
    wg.Wait()
}

func NewPluginServer(plugin Plugin) (*Server, error) {
    // XXX can we actually handle having multiple "Plugin" names registered?
    rpc.RegisterName("Plugin", plugin)
    rpc.HandleHTTP()
    listener, port, err := getListener()
    if err != nil {
        return nil, err
    }

    s := &Server{port, listener, make(chan *rpc.Call)}

    // handle exactly one manager ready reply
    go func() {
        for call := range s.readyChan {
            if call.Error != nil {
                log.Fatal(call.Error)
            }
            break
        }
        close(s.readyChan)
    }()

    return s, nil
}

func getListener() (net.Listener, int, error) {
    for port := control.CONTROL_PORT_MIN; port <= control.CONTROL_PORT_MAX; port++ {
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
