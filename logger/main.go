package logger

import (
    "log"
    "net/rpc"

    "github.com/Senior-Design-May1601/projectmain/control"
)

type LogWriter struct {
    client *rpc.Client
}

func (x *LogWriter) Write(p []byte) (n int, err error) {
    // TODO: should this be async?
    var r int
    e := x.client.Call("LogManager.Log", p, &r)
    if e != nil {
        return 0, e
    }
    return len(p), nil
}

func NewLogger(prefix string, flag int) *log.Logger {
    client, err := rpc.DialHTTP("tcp", control.CONTROL_PORT_CORE)
    if err != nil {
        log.Fatal(err)
    }
    writer := &LogWriter{client}
    return log.New(writer, prefix, flag)
}
