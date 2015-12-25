package main 

import (
	"net/http"
	"net/rpc"
	"log"
	"net"
	"../plugin/"
	"encoding/gob"	
	"fmt"

)
func main(){
	rpc.Register(&Handler{8081})
	rpc.HandleHTTP()
	gob.RegisterName("HoneyPot",new(plugin.Honeypot))
        gob.RegisterName("Arith",new(plugin.Arith))

	l, e:= net.Listen("tcp",":8080")
	if e != nil{
		log.Fatal("Listen error:", e)
	}
	http.Serve(l,nil)
}

type Handler struct{
	port int
}
func (t *Handler) RegisterType(args plugin.Honeypot,reply *int) error{
	rpc.Register(args)
	port := fmt.Sprintf(":%d",t.port)
	*reply = t.port
	t.port++
	l, e:= net.Listen("tcp",port)
	if e != nil{
		log.Fatal("Listen error:", e)
	}
	go http.Serve(l,nil)
	return nil
}


