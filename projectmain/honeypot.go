package projectmain

// Interface declaration for a projectmain servce.

type Honeypot interface {
    Start()     error
    Stop()      error
    Restart()   error
}
