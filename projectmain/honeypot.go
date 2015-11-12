package projectmain

// Interface declaration for a projectmain servce.

type Servce interface {
    Start()     error
    Stop()      error
    Restart()   error
}
