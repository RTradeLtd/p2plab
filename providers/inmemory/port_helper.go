package inmemory

import (
	"sync"

	"github.com/phayes/freeport"
)

// portHelper enables synchronizing access to free ports
// across multiple different trials and benchmarks, and prevents
// accidental port re-use.
type portHelper struct {
	inUse map[int]bool
	mux   sync.Mutex
}

func (ph *portHelper) getPorts(num int) (int, int, error) {
	ph.mux.Lock()
	defer ph.mux.Unlock()
	for {
		freePorts, err := freeport.GetFreePorts(2)
		if err != nil {
			return 0, 0, err
		}
		if ph.inUse[freePorts[0]] || ph.inUse[freePorts[1]] {
			continue
		}
		ph.inUse[freePorts[0]] = true
		ph.inUse[freePorts[1]] = true
		return freePorts[0], freePorts[1], nil
	}
}
