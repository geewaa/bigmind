package mind

import (
	"log"
	"sync"
	"sync/atomic"
)

type ENVGET func() map[string]any
type ENVSET func(map[string]any)
type ACTION func(Minder) error

type BaseMinder struct {
	// m      *Minder
	action     ACTION
	envget     ENVGET
	envset     ENVSET
	ev         *sync.Map
	notice     chan string
	EnvChanged *atomic.Bool
}

func NewMinder(m *Minder, c chan string, get ENVGET, set ENVSET) *BaseMinder {
	return &BaseMinder{
		// m:      m,
		envget:     get,
		envset:     set,
		notice:     c,
		ev:         new(sync.Map),
		EnvChanged: &atomic.Bool{},
	}
}

// mustEmbedBaseMinder a new struct must embed BaseMinder
func (mr *BaseMinder) mustEmbedBaseMinder() {}

func (mr *BaseMinder) Init(c chan string, get ENVGET, set ENVSET) {
	mr.envget = get
	mr.envset = set
	mr.notice = c
	mr.ev = new(sync.Map)
	mr.EnvChanged = &atomic.Bool{}
}

func (mr *BaseMinder) SetAction(f ACTION) {
	mr.action = f
}

func (mr *BaseMinder) Feedback() {
	// mr.m.Upgrade()
	ev := make(map[string]any)
	mr.ev.Range(
		func(key, value any) bool {
			ev[key.(string)] = value
			return true
		})
	mr.envset(ev)
}
func (mr *BaseMinder) initenv() {
	defer mr.EnvChanged.Store(false)

	ev := mr.envget()
	for k, v := range ev {
		if v != nil {
			mr.ev.Store(k, v)
		}
	}
}

func (mr *BaseMinder) Work(wg *sync.WaitGroup, m Minder) {
	go func() {
		for msg := range mr.notice {
			if msg == "upgrade" {
				// do something
			}
			ev := mr.envget()
			for k, v := range ev {
				if v != nil {
					mr.ev.Store(k, v)
				}
			}
			mr.EnvChanged.Store(true)
		}
	}()
	mr.initenv()
	err := mr.action(m)
	if err != nil {
		log.Printf("%s-got error:%s\n", "BaseMinder", err.Error())
	}
	mr.Feedback()
	wg.Done()
}

// getenv
func (mr *BaseMinder) GetEnv() *sync.Map {
	return mr.ev
}

// getenv
func (mr *BaseMinder) UpdateEnv(nv map[string]any) error {
	for k, v := range nv {
		mr.ev.Store(k, v)
	}
	return nil
}
