package loggerplugin

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/Senior-Design-May1601/projectmain/control"
)

type LoggerPluginServer struct {
	Port      int
	Name      string
	listener  net.Listener
	readyChan chan *rpc.Call
	client    *rpc.Client
}

type LoggerPlugin interface {
	Log(msg []byte, _ *int) error
}

type ReadyArg struct {
	Port int
	Name string
}

func NewLoggerPlugin(logger LoggerPlugin) (*LoggerPluginServer, error) {
	client, err := rpc.DialHTTP("tcp", control.CONTROL_PORT_CORE)
	if err != nil {
		return nil, err
	}

	name := id()

	rpc.RegisterName(name, logger)
	rpc.HandleHTTP()
	listener, port, err := getListener()
	if err != nil {
		return nil, err
	}

	s := &LoggerPluginServer{
		Port:      port,
		Name:      name,
		listener:  listener,
		readyChan: make(chan *rpc.Call, 1),
		client:    client}

	// handle exactly one manager ready reply
	go func() {
		for call := range s.readyChan {
			if call.Error != nil {
				log.Fatal(call.Error)
			}
			close(s.readyChan)
		}
	}()

	return s, nil
}

func (x *LoggerPluginServer) Run() error {
	var r int
	x.client.Go("LogManager.Ready", ReadyArg{x.Port, x.Name}, &r, x.readyChan)
	return http.Serve(x.listener, nil)
}

func getListener() (net.Listener, int, error) {
	for port := control.CONTROL_PORT_MIN; port <= control.CONTROL_PORT_MAX; port++ {
		l, e := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
		if e == nil {
			return l, port, nil
		}
	}
	return nil, 0, errors.New("No available control ports.")
}
