package plugin

import (
    "github.com/Senior-Design-May1601/projectmain/projectmain"
)

// TODO: actually setup a Service interface
func cmdHoneypot struct {
    honeypot projectmain.Honeypot
}

func (c *cmdHoneypot) Start() error {
    return c.honeypot.Start()
}

func (c *cmdHoneypot) Stop() error {
    return c.honeypot.Stop()
}

func (c *cmdHoneypot) Restart() error {
    return c.honeypot.Restart()
}
