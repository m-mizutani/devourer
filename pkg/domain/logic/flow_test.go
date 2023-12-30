package logic_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/devourer/pkg/domain/logic"
	"github.com/m-mizutani/devourer/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestFlows(t *testing.T) {
	flowMap := logic.NewFlowMap()

	now := time.Now()

	flows := []*model.Flow{
		{

			Protocol: "tcp",
			Src: model.Peer{
				Addr: []byte{192, 168, 0, 1},
				Port: 12345,
			},
			Dst: model.Peer{
				Addr: []byte{10, 0, 0, 1},
				Port: 25678,
			},
			FirstSeenAt: now,
			LastSeenAt:  now,
		},
		{

			Protocol: "tcp",
			Src: model.Peer{
				Addr: []byte{192, 168, 0, 1},
				Port: 12345,
			},
			Dst: model.Peer{
				Addr: []byte{10, 0, 0, 2},
				Port: 25678,
			},
			FirstSeenAt: now,
			LastSeenAt:  now,
		},
		{

			Protocol: "tcp",
			Src: model.Peer{
				Addr: []byte{192, 168, 0, 1},
				Port: 12345,
			},
			Dst: model.Peer{
				Addr: []byte{10, 0, 0, 1},
				Port: 80,
			},
			FirstSeenAt: now,
			LastSeenAt:  now,
		},
	}

	t.Run("put flow", func(t *testing.T) {
		gt.True(t, flowMap.Put(flows[0]))
		flows[0].LastSeenAt = now.Add(time.Second)
		gt.False(t, flowMap.Put(flows[0]))
		f := flowMap.Get(flows[0].Key())
		gt.Equal(t, f.LastSeenAt, now.Add(time.Second))
	})
}
