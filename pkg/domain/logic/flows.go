package logic

import (
	"sync"
	"time"

	"github.com/m-mizutani/devourer/pkg/domain/model"
)

type FlowMap struct {
	storage map[model.FlowKey]*model.Flow
	mutex   sync.RWMutex
}

func NewFlowMap() *FlowMap {
	return &FlowMap{
		storage: make(map[model.FlowKey]*model.Flow),
		mutex:   sync.RWMutex{},
	}
}

func (x *FlowMap) Put(flow *model.Flow, stat model.PeerStat) bool {
	key := flow.Key()

	x.mutex.Lock()
	defer x.mutex.Unlock()

	if exist, ok := x.storage[key]; ok {
		exist.Update(&flow.Src, flow.LastSeenAt, stat)
		return false
	}

	x.storage[key] = flow
	return true
}

func (x *FlowMap) Get(key model.FlowKey) *model.Flow {
	x.mutex.RLock()
	defer x.mutex.RUnlock()
	return x.storage[key]
}

func (x *FlowMap) Expire(key model.FlowKey, at time.Time) bool {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	if exist, ok := x.storage[key]; ok {
		if exist.LastSeenAt.Before(at) {
			delete(x.storage, key)
			return true
		}
	}

	return false
}

func (x *FlowMap) Flush() []*model.Flow {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	flows := make([]*model.Flow, 0, len(x.storage))
	for key := range x.storage {
		flows = append(flows, x.storage[key])
	}
	x.storage = make(map[model.FlowKey]*model.Flow)

	return flows
}
