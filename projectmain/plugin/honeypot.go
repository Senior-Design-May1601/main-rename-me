package plugin

// Interface declaration for a projectmain servce.

type Honeypot interface {
    Start( a *struct{}, b *struct{})     error
    Stop()      error
    Restart()   error
}
