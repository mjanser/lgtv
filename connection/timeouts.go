package connection

import (
	"time"
)

const (
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 1 * time.Second
	defaultPingTimeout     = 5 * time.Second
	defaultResponseTimeout = 2 * time.Second
	defaultRegisterTimeout = 1 * time.Minute
)

// Timeouts defines various timeouts for the connection
type Timeouts struct {
	Read     time.Duration
	Write    time.Duration
	Ping     time.Duration
	Response time.Duration
	Register time.Duration
}

// DefaultTimeouts returns the default timeouts
func DefaultTimeouts() Timeouts {
	return Timeouts{
		Read:     defaultReadTimeout,
		Write:    defaultWriteTimeout,
		Ping:     defaultPingTimeout,
		Response: defaultResponseTimeout,
		Register: defaultRegisterTimeout,
	}
}
