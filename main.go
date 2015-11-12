package main

import (
    "os"
    "runtime"

    "github.com/Senior-Design-May1601/projectmain/projectmain"
)

// TODO: make main package for projectmain application
func main() {
    if os.Getenv("GOMAXPROCS") == "" {
        runtime.GOMAXPROCS(runtime.NumCPU())
    }
}
