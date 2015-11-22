package main

import(
"net/rpc"
"log"
"net"
"net/http"
"fmt"
)

type RpcServer int
var port int
func (t *RpcServer) RegisterType(args interface{}, reply *int) error{
	rpc.Register(args)
	
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
	arith := new(RpcServer)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp","127.0.0.1:1234")
	if e != nil{
		log.Fatal("listen error:",e)
	}
	port = 1235
        http.Serve(l,nil)

}
