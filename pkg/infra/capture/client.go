package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/m-mizutani/devourer/pkg/domain/interfaces"
	"github.com/m-mizutani/goerr"
)

func NewDevice(iface string) (interfaces.Capture, error) {
	const (
		snapshotLen int32 = 1024 * 1024
		promiscuous bool  = true
	)
	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, pcap.BlockForever)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to open device")
	}
	source := gopacket.NewPacketSource(handle, handle.LinkType())

	return &client{
		source: source,
	}, nil
}

func NewFile(fpath string) (interfaces.Capture, error) {
	handle, err := pcap.OpenOffline(fpath)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to open file")
	}
	source := gopacket.NewPacketSource(handle, handle.LinkType())

	return &client{
		source: source,
	}, nil
}

type client struct {
	source *gopacket.PacketSource
}

func (x *client) Read() chan gopacket.Packet {
	return x.source.Packets()
}
