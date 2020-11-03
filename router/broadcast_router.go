package router

import (
	"sync/atomic"
	"unsafe"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type broadcastGroupRouter struct {
	GroupRouter
}

type broadcastPoolRouter struct {
	PoolRouter
}

type broadcastRouterState struct {
	routees *actor.PIDSet
	sender  actor.SenderContext
}

func (state *broadcastRouterState) SetSender(sender actor.SenderContext) {
	state.sender = sender
}

func (state *broadcastRouterState) SetRoutees(routees *actor.PIDSet) {
	rts := *routees
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&state.routees)), unsafe.Pointer(&rts))
}

func (state *broadcastRouterState) GetRoutees() *actor.PIDSet {
	rts := (*actor.PIDSet)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&state.routees))))
	return rts.Clone()
}

func (state *broadcastRouterState) RouteMessage(message interface{}) {
	rts := (*actor.PIDSet)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&state.routees))))
	rts.ForEach(func(i int, pid *actor.PID) {
		state.sender.Send(pid, message)
	})
}

func NewBroadcastPool(size int) *actor.Props {
	return (&actor.Props{}).WithSpawnFunc(spawner(&broadcastPoolRouter{PoolRouter{PoolSize: size}}))
}

func NewBroadcastGroup(routees ...*actor.PID) *actor.Props {
	return (&actor.Props{}).WithSpawnFunc(spawner(&broadcastGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}}))
}

func (config *broadcastPoolRouter) CreateRouterState() RouterState {
	return &broadcastRouterState{}
}

func (config *broadcastGroupRouter) CreateRouterState() RouterState {
	return &broadcastRouterState{}
}
