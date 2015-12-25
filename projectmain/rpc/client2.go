package main
import("net/rpc"
	"log"
	"../plugin/"
	"encoding/gob"
//	"bytes"
//	"net/rpc/jsonrpc"
	"fmt"
)

func main(){
	var reply int
	client, err := rpc.DialHTTP("tcp",":8080")
	if err != nil{
		log.Fatal("dialing:",err)
	}
	
	var hp plugin.Honeypot
	hp = plugin.Arith{}
	gob.RegisterName("Honeypot",new(plugin.Honeypot))
	gob.RegisterName("Arith",plugin.Arith{});
	err = client.Call("Handler.RegisterType",&hp,&reply)
	if err != nil{
		log.Fatal("call err ",err)
	}
	fmt.Printf("%d \n",reply)
	client.Close()
}
