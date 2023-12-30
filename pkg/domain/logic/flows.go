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

func (x *FlowMap) Put(flow *model.Flow) bool {
	key := flow.Key()

	x.mutex.Lock()
	defer x.mutex.Unlock()

	if exist, ok := x.storage[key]; ok {
		exist.Update(&flow.Src, flow.LastSeenAt, flow.SrcStat)
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

func (x *FlowMap) Expire(at time.Time) []*model.Flow {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	var expired []*model.Flow
	for key := range x.storage {
		if x.storage[key].LastSeenAt.Before(at) {
			expired = append(expired, x.storage[key])
			delete(x.storage, key)
		}
	}

	return expired
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

func (x *FlowMap) Len() int {
	x.mutex.RLock()
	defer x.mutex.RUnlock()
	return len(x.storage)
}
