package interfaces

import (
	"context"

	"github.com/google/gopacket"
	"github.com/m-mizutani/devourer/pkg/domain/model"
)

type Capture interface {
	Read() chan gopacket.Packet
}

type Dumper interface {
	Dump(ctx context.Context, record *model.Record) error
	Close()
}
