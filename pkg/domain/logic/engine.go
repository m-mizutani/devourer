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
	var newFlows []*model.Flow

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
	)
	stat := model.PeerStat{
		Bytes:   uint64(pkt.Metadata().Length),
		Packets: 1,
	}

	if x.flowMap.Put(flow, stat) {
		newFlows = append(newFlows, flow)
	}

	return &model.Record{
		NewFlows: newFlows,
	}, nil
}

func (x *Engine) Tick(now time.Time) (*model.Record, error) {
	return &model.Record{}, nil
}

func (x *Engine) Flush() *model.Record {
	return &model.Record{}
}
