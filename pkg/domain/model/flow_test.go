package model_test

import (
	"net"
	"testing"

	"github.com/m-mizutani/devourer/pkg/domain/model"
)

func TestFlowKey(t *testing.T) {
	base := model.Flow{
		Protocol: "tcp",
		Src: model.Peer{
			Addr: net.IPv4(192, 168, 0, 1),
			Port: 12345,
		},
		Dst: model.Peer{
			Addr: net.IPv4(10, 0, 0, 1),
			Port: 25678,
		},
	}

	testCases := map[string]struct {
		flow   model.Flow
		isSame bool
	}{
		"same": {
			flow:   base,
			isSame: true,
		},
		"same but reverse": {
			flow: model.Flow{
				Protocol: "tcp",
				Src: model.Peer{
					Addr: net.IPv4(10, 0, 0, 1),
					Port: 25678,
				},
				Dst: model.Peer{
					Addr: net.IPv4(192, 168, 0, 1),
					Port: 12345,
				},
			},
			isSame: true,
		},
		"different protocol": {
			flow: model.Flow{
				Protocol: "udp",
				Src: model.Peer{
					Addr: net.IPv4(192, 168, 0, 1),
					Port: 12345,
				},
				Dst: model.Peer{
					Addr: net.IPv4(10, 0, 0, 1),
					Port: 25678,
				},
			},
			isSame: false,
		},
		"different src addr": {
			flow: model.Flow{
				Protocol: "tcp",
				Src: model.Peer{
					Addr: net.IPv4(192, 168, 0, 2),
					Port: 12345,
				},
				Dst: model.Peer{
					Addr: net.IPv4(10, 0, 0, 1),

					Port: 25678,
				},
			},
			isSame: false,
		},
		"different src port": {
			flow: model.Flow{
				Protocol: "tcp",
				Src: model.Peer{
					Addr: net.IPv4(192, 168, 0, 1),
					Port: 12346,
				},
				Dst: model.Peer{
					Addr: net.IPv4(10, 0, 0, 1),
					Port: 25678,
				},
			},
			isSame: false,
		},
		"different dst addr": {
			flow: model.Flow{
				Protocol: "tcp",
				Src: model.Peer{
					Addr: net.IPv4(192, 168, 0, 1),
					Port: 12345,
				},
				Dst: model.Peer{
					Addr: net.IPv4(10, 0, 0, 2),
					Port: 25678,
				},
			},
			isSame: false,
		},
		"different dst port": {
			flow: model.Flow{
				Protocol: "tcp",
				Src: model.Peer{
					Addr: net.IPv4(192, 168, 0, 1),
					Port: 12345,
				},
				Dst: model.Peer{
					Addr: net.IPv4(10, 0, 0, 1),
					Port: 25679,
				},
			},
			isSame: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			key1 := model.CalcFlowKey(&base.Src, &base.Dst, base.Protocol)
			key2 := model.CalcFlowKey(&tc.flow.Src, &tc.flow.Dst, tc.flow.Protocol)
			if tc.isSame && key1 != key2 {
				t.Errorf("Different key: %d != %d", key1, key2)
			} else if !tc.isSame && key1 == key2 {
				t.Errorf("Same key: %d == %d", key1, key2)
			}
		})
	}
}
