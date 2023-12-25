package mind

import "sync"

type Minder interface {
	mustEmbedBaseMinder()
	Init(c chan string, get ENVGET, set ENVSET)
	Work(wg *sync.WaitGroup, m Minder)
	SetAction(f ACTION)

	GetEnv() *sync.Map
	UpdateEnv(nv map[string]any) error
}
