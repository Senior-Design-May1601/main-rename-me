package main

import(
	"net/rpc"
	"log"
	"net"
	"net/http"
	"fmt"
	"encoding/gob"
	"../plugin/"
)

type RpcServer struct{

}
var port int

func (t *RpcServer) RegisterType(args plugin.HoneyPot, reply *int) error{
	rpc.Register(args.Instance)
	address:= fmt.Sprintf("127.0.0.1:%d",port)
	l, e := net.Listen("tcp",address)
	if e != nil{
		log.Fatal("listen error:",e)
	}
	fmt.Printf("Listening on %d",port)
	go http.Serve(l,nil)
	port++
	*reply = 1 
	return nil
}


func main(){
	Server := new(RpcServer)
	rpc.Register(Server)
	gob.Register(plugin.Arith{})
	rpc.HandleHTTP()
	l, e := net.Listen("tcp","127.0.0.1:1234")
	if e != nil{
		log.Fatal("listen error:",e)
	}
	port = 1235
        http.Serve(l,nil)
}
