package logic

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/m-mizutani/devourer/pkg/domain/model"
)

type Engine struct {
	timeout time.Duration
	flowMap *FlowMap
}

func NewEngine() *Engine {
	return &Engine{
		timeout: 120 * time.Second,
		flowMap: NewFlowMap(),
	}
}

type Option func(*Engine)

func WithTimeout(d time.Duration) Option {
	return func(x *Engine) {
		x.timeout = d
	}
}

func (x *Engine) InputPacket(pkt gopacket.Packet) (*model.Record, error) {
	if pkt.NetworkLayer() == nil || pkt.TransportLayer() == nil {
		// not supported
		return nil, nil
	}

	netLayer := pkt.NetworkLayer().NetworkFlow()

	var proto string
	var srcPort, dstPort uint32
	switch pkt.TransportLayer().LayerType() {
	case layers.LayerTypeTCP:
		tcpLayer := pkt.TransportLayer().(*layers.TCP)
		srcPort = uint32(tcpLayer.SrcPort)
		dstPort = uint32(tcpLayer.DstPort)
		proto = "tcp"
	case layers.LayerTypeUDP:
		udpLayer := pkt.TransportLayer().(*layers.UDP)
		srcPort = uint32(udpLayer.SrcPort)
		dstPort = uint32(udpLayer.DstPort)
		proto = "udp"
	case layers.LayerTypeICMPv4:
		proto = "icmp4"
	case layers.LayerTypeICMPv6:
		proto = "icmp6"
	default:
		// not supported
		return nil, nil
	}

	flow := model.NewFlow(
		model.Peer{
			Addr: netLayer.Src().Raw(),
			Port: srcPort,
		},
		model.Peer{
			Addr: netLayer.Dst().Raw(),
			Port: dstPort,
		},
		proto,
		pkt.Metadata().Timestamp,
		model.PeerStat{
			Bytes:   uint64(pkt.Metadata().Length),
			Packets: 1,
		},
	)

	_ = x.flowMap.Put(flow)

	return &model.Record{}, nil
}

func (x *Engine) Tick(now time.Time) (*model.Record, error) {
	return &model.Record{
		FlowLogs: x.flowMap.Expire(now.Add(-x.timeout)),
	}, nil
}

func (x *Engine) Flush() *model.Record {
	return &model.Record{
		FlowLogs: x.flowMap.Flush(),
	}
}

func (x *Engine) FlowCount() int {
	return x.flowMap.Len()
}
