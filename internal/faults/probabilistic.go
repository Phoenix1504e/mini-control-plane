package faults

import "math/rand"

type Injector interface {
    Decide() Decision
}

type ProbabilisticInjector struct {
    DropRate float64
}

func (p *ProbabilisticInjector) Decide() Decision {
    if rand.Float64() < p.DropRate {
        return Decision{
            Action: Drop,
        }
    }

    return Decision{
        Action: Pass,
    }
}
