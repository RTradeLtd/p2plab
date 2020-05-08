package inmemory

import (
	"context"
	"sync"

	"github.com/phayes/freeport"
	"github.com/rs/zerolog"
)

// portHelper enables synchronizing access to free ports
// across multiple different trials and benchmarks, and prevents
// accidental port re-use.
type portHelper struct {
	inUse map[int]bool
	mux   sync.Mutex
}

func (ph *portHelper) getPorts(ctx context.Context, num int) (int, int, error) {
	ph.mux.Lock()
	defer ph.mux.Unlock()
	for {
		zerolog.Ctx(ctx).Info().Msg("getting free ports")
		freePorts, err := freeport.GetFreePorts(2)
		if err != nil {
			return 0, 0, err
		}
		if ph.inUse[freePorts[0]] || ph.inUse[freePorts[1]] {
			zerolog.Ctx(ctx).Warn().Msgf("got in use ports, %v and %v, trying again", freePorts[0], freePorts[1])
			continue
		}
		zerolog.Ctx(ctx).Info().Msg("found available ports")
		ph.inUse[freePorts[0]] = true
		ph.inUse[freePorts[1]] = true
		return freePorts[0], freePorts[1], nil
	}
}

func (ph *portHelper) returnPorts(ctx context.Context, ports []int) {
	ph.mux.Lock()
	defer ph.mux.Unlock()
	returned := 0
	for _, port := range ports {
		if !ph.inUse[port] {
			zerolog.Ctx(ctx).Warn().Msg("trying to return unused port")
			continue
		}
		ph.inUse[port] = false
		returned++
	}
	if returned == 0 {
		zerolog.Ctx(ctx).Warn().Msg("no ports retruend")
	} else {
		zerolog.Ctx(ctx).Info().Msgf("successfully returned %v ports", returned)
	}
}
