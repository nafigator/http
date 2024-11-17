package debug

import "context"

type logger interface {
	Debug(args ...interface{})
}

type Debug struct {
	log logger
}

// New creates Debug instance.
func New(log logger) *Debug {
	return &Debug{
		log: log,
	}
}

// Flush sends message into debug logger.
func (d *Debug) Flush(_ context.Context, msg string) {
	d.log.Debug(msg)
}
