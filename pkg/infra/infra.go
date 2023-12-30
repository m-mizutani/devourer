package infra

import (
	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
)

type Clients struct {
	capture interfaces.Capture
	dumper  interfaces.Dumper
}

type Option func(*Clients)

func New(opts ...Option) *Clients {
	x := &Clients{}
	for _, opt := range opts {
		opt(x)
	}

	return x
}

func WithCapture(capture interfaces.Capture) Option {
	return func(x *Clients) {
		x.capture = capture
	}
}

func (x *Clients) Capture() interfaces.Capture {
	return x.capture
}

func WithDumper(dumper interfaces.Dumper) Option {
	return func(x *Clients) {
		x.dumper = dumper
	}
}

func (x *Clients) Dumper() interfaces.Dumper {
	return x.dumper
}
