package main
import("net/rpc"
	"log"
	"bytes"
	"encoding/gob"	
)

type Args struct{
	A, B int
}
type Client struct{
	conn *rpc.Client
}
type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error{
	*reply = args.A * args.B
	return nil
}


func main(){
	var reply int
	var buf []byte
	b:= bytes.NewBuffer(buf)
	enc := gob.NewEncoder(b)
	err := enc.Encode(new(Arith))

	if err != nil{
		log.Fatal("Encoding",err)
	}

	client, err := rpc.DialHTTP("tcp","127.0.0.1:1234")
	if err != nil{
		log.Fatal("dialing:",err)
	}
	
	err = client.Call("RpcServer.RegisterType",b.Bytes(),&reply)
	if err != nil{
		log.Fatal("arith err:",err)
	}
	


}
