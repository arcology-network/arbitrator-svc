package arbitrator

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/arcology-network/arbitrator-svc/types"
	"github.com/arcology-network/urlarbitrator-engine/go-wrapper"
)

type ArbitratorImpl struct {
	Txs        []uint32
	Paths      []string
	Reads      []uint32
	Writes     []uint32
	Composite  []bool
	arbitrator unsafe.Pointer
}

func NewArbitratorImpl() *ArbitratorImpl {
	preAllocSize := 5000000
	return &ArbitratorImpl{
		Txs:        make([]uint32, 0, preAllocSize),
		Paths:      make([]string, 0, preAllocSize),
		Reads:      make([]uint32, 0, preAllocSize),
		Writes:     make([]uint32, 0, preAllocSize),
		Composite:  make([]bool, 0, preAllocSize),
		arbitrator: wrapper.Start(),
	}
}

func (arb *ArbitratorImpl) Reset() {
	arb.Txs = arb.Txs[:0]
	arb.Paths = arb.Paths[:0]
	arb.Reads = arb.Reads[:0]
	arb.Writes = arb.Writes[:0]
	arb.Composite = arb.Composite[:0]
}

func (arb *ArbitratorImpl) DetectConflict(groups [][]*types.ProcessedEuResult) ([]uint32, []uint32, []bool, []uint32, []uint32, []time.Duration, time.Time, int) {
	whitelist := make([]uint32, 0, len(groups))
	indexToID := make(map[uint32]uint32)
	tims := make([]time.Duration, 6)
	begintime := time.Now()
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
			arb.Composite = append(arb.Composite, per.Composite...)
		}
		if len(g[0].Txs) == 0 {
			continue
		}
		indexToID[uint32(i)] = g[0].Txs[0]
		whitelist = append(whitelist, uint32(i))
	}
	tims[0] = time.Now().Sub(begintime)

	connectstrTime, buf := wrapper.Insert(arb.arbitrator, arb.Txs, arb.Paths, arb.Reads, arb.Writes, arb.Composite)
	defer wrapper.Clear(arb.arbitrator, buf)
	t0 := time.Now()
	begintime = time.Now()
	txs, g, flags := wrapper.Detect(arb.arbitrator, whitelist)
	fmt.Printf("len(arb.Txs): %d, wrapper.Detect: %v\n", len(arb.Txs), time.Since(t0))
	tims[2] = time.Now().Sub(begintime)
	begintime = time.Now()

	l, r := wrapper.ExportTxs(arb.arbitrator)
	tims[3] = time.Now().Sub(begintime)
	begintime = time.Now()
	left := make([]uint32, len(l))
	right := make([]uint32, len(r))
	for i := range left {
		left[i] = indexToID[l[i]]
		right[i] = indexToID[r[i]]
	}
	tims[4] = time.Now().Sub(begintime)
	tims[5] = connectstrTime
	return txs, g, flags, left, right, tims, time.Now(), len(arb.Txs)
}
