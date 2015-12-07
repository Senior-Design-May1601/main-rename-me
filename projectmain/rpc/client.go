package main
import("net/rpc"
	"log"
	"encoding/gob"
	"bytes"
)

type Args struct{
	A, B int
}
type Client struct{
	conn *rpc.Client
}
type Arith struct{}
type Walker interface{
	Multiply(args *Args, reply *int) error
}
func (t Arith) Multiply(args *Args, reply *int) error{
	*reply = args.A * args.B
	return nil
}


func main(){
	var reply int
	gob.Register(Arith{})
	arith := Arith{}
	var b bytes.Buffer
	var w Walker = arith
	e:= gob.NewEncoder(&b)
	
	if err:= e.Encode(&w); err != nil{
		panic(err)
	}
	
//	jArith,err := json.Marshal(arith);
/*	if err != nil{
		log.Fatal("Encoding",err)
	}*/

	client, err := rpc.DialHTTP("tcp","127.0.0.1:1234")
	if err != nil{
		log.Fatal("dialing:",err)
	}
	
	err = client.Call("RpcServer.RegisterMe",b,&reply)
	if err != nil{
		log.Fatal("call err:",err)
	}
	


}
