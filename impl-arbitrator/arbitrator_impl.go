package arbitrator

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	types "github.com/HPISTechnologies/arbitrator-svc/types"
	wrapper "github.com/HPISTechnologies/urlarbitrator-engine/go-wrapper"
)

type ArbitratorImpl struct {
	Txs         []uint32
	Paths       []string
	Reads       []uint32
	Writes      []uint32
	AddOrDelete []bool
	Composite   []bool
}

var arbitrator unsafe.Pointer
var initArbitrator sync.Once

func NewArbitratorImpl() *ArbitratorImpl {
	preAllocSize := 5000000
	return &ArbitratorImpl{
		Txs:         make([]uint32, 0, preAllocSize),
		Paths:       make([]string, 0, preAllocSize),
		Reads:       make([]uint32, 0, preAllocSize),
		Writes:      make([]uint32, 0, preAllocSize),
		AddOrDelete: make([]bool, 0, preAllocSize),
		Composite:   make([]bool, 0, preAllocSize),
	}
}

func (arb *ArbitratorImpl) Reset() {
	arb.Txs = arb.Txs[:0]
	arb.Paths = arb.Paths[:0]
	arb.Reads = arb.Reads[:0]
	arb.Writes = arb.Writes[:0]
	arb.AddOrDelete = arb.AddOrDelete[:0]
	arb.Composite = arb.Composite[:0]
}

func (arb *ArbitratorImpl) DetectConflict(groups [][]*types.ProcessedEuResult) ([]uint32, []uint32, []bool, []uint32, []uint32, []time.Duration, time.Time, int) {
	indexToID := make(map[uint32]uint32)
	tims := make([]time.Duration, 6)
	t0 := time.Now()
	for i, g := range groups {
		if len(g) == 0 {
			continue
		}
		for _, per := range g {
			for range per.Paths {
				arb.Txs = append(arb.Txs, uint32(i))
			}
			arb.Paths = append(arb.Paths, per.Paths...)
			arb.Reads = append(arb.Reads, per.Reads...)
			arb.Writes = append(arb.Writes, per.Writes...)
			arb.AddOrDelete = append(arb.AddOrDelete, per.AddOrDelete...)
			arb.Composite = append(arb.Composite, per.Composite...)
		}
		if len(g[0].Txs) == 0 {
			continue
		}
		indexToID[uint32(i)] = g[0].Txs[0]
	}
	tims[0] = time.Since(t0)

	t0 = time.Now()
	buf, ConcateStrTime := wrapper.Insert(getArbitrator(), arb.Txs, arb.Paths, arb.Reads, arb.Writes, arb.AddOrDelete, arb.Composite)
	tims[1] = time.Since(t0)
	defer wrapper.Clear(getArbitrator(), buf)

	t0 = time.Now()
	txs, g, flags := wrapper.Detect(getArbitrator(), uint32(len(arb.Txs)))
	fmt.Printf("len(arb.Txs): %d, wrapper.Detect: %v\n", len(arb.Txs), time.Since(t0))
	tims[2] = time.Since(t0)
	t0 = time.Now()

	l, r := wrapper.ExportTxs(getArbitrator())
	tims[3] = time.Since(t0)
	t0 = time.Now()
	left := make([]uint32, len(l))
	right := make([]uint32, len(r))
	for i := range left {
		left[i] = indexToID[l[i]]
		right[i] = indexToID[r[i]]
	}
	tims[4] = time.Since(t0)
	tims[5] = ConcateStrTime
	return txs, g, flags, left, right, tims, time.Now(), len(arb.Txs)
}

func getArbitrator() unsafe.Pointer {
	initArbitrator.Do(func() {
		arbitrator = wrapper.Start()
	})
	return arbitrator
}
