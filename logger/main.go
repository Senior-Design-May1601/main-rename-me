package logger

import (
	"log"
	"net/rpc"

	"github.com/Senior-Design-May1601/projectmain/control"
)

type LogWriter struct {
	client *rpc.Client
    responseChan chan *rpc.Call
}

func NewLogger(prefix string, flag int) *log.Logger {
	client, err := rpc.DialHTTP("tcp", control.CONTROL_PORT_CORE)
	if err != nil {
		log.Fatal(err)
	}

	writer := &LogWriter{
        client: client,
        responseChan: make(chan *rpc.Call, 100)}

    go writer.handleResponses()

	return log.New(writer, prefix, flag)
}

func (x *LogWriter) Write(p []byte) (n int, err error) {
	var r int
	x.client.Go("LogManager.Log", p, &r, x.responseChan)
	return len(p), nil
}

func (x *LogWriter) handleResponses() {
    for response := range x.responseChan {
        if response.Error != nil {
            log.Fatal("Log manager error:", response.Error)
        }
    }
}
