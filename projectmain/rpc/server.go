package main 

import (
	"net/rpc"
	"log"
	"fmt"
	"net"
)

var port int;


type Server struct{
	server *rpc.Server
}

func NewServer() (*Server){
	return &Server{server: rpc.NewServer()}	
}

func (s *Server) Register(rcvr interface{},reply *int) error{
	fmt.Printf("I'm through the rabbit hole\n")
	s.server.Register(rcvr)
	address := fmt.Sprintf("127.0.0.1:%d",port)
	l,e := net.Listen("tcp",address)
	if e != nil{
		log.Fatal("listen error:",e)
	}
	fmt.Printf("listening on %d",port)
	go s.server.Accept(l)
	port++
	return nil 
}
func main(){
	s:=NewServer()
	s.server.Register(s)
	s.server.HandleHTTP("/_goRPC_","/debug/rpc")
	l, e := net.Listen("tcp","127.0.0.1:1234")
	if e != nil{
		log.Fatal("listen error:",e)
	}
	port = 1235
        s.server.Accept(l)
}

