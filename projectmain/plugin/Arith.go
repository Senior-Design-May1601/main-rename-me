package plugin
import(
	"fmt"
)

type Arith struct{
}

func (t Arith) Start(args int,reply *int) error{
	fmt.Printf("I have started")
	return nil
}
func (t Arith) Stop(args int,reply *int) error{
	fmt.Printf("I have Stopped")
	return nil
}
func (t Arith) Restart(args int, reply *int) error{
	fmt.Printf("I have Restarted")
	return nil
}



