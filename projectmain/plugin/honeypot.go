package plugin

// Interface declaration for a projectmain servce.

type Honeypot interface {
	Start(args int, reply *int)     error
   	Stop(args int, reply *int)      error
	Restart(args int, reply *int)   error

   } 
