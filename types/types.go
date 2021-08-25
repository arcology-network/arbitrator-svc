package types

import (
	"math/big"

	ctypes "github.com/HPISTechnologies/common-lib/types"
	urltype "github.com/HPISTechnologies/concurrenturl/v2/type"
	"github.com/HPISTechnologies/concurrenturl/v2/type/commutative"
)

type BalanceTransition struct {
	Path   string
	Tx     uint32
	Origin *big.Int
	Delta  *big.Int
}

type ProcessedEuResult struct {
	Hash        string
	Txs         []uint32
	Paths       []string
	Reads       []uint32
	Writes      []uint32
	AddOrDelete []bool
	Composite   []bool
	Transitions []*BalanceTransition
}

func Process(ars *ctypes.TxAccessRecords) *ProcessedEuResult {
	length := len(ars.Accesses)
	per := &ProcessedEuResult{
		Hash:        ars.Hash,
		Txs:         make([]uint32, length),
		Paths:       make([]string, length),
		Reads:       make([]uint32, length),
		Writes:      make([]uint32, length),
		AddOrDelete: make([]bool, length),
		Composite:   make([]bool, length),
		Transitions: make([]*BalanceTransition, 0, length),
	}
	for i, univalue := range ars.Accesses {
		per.Txs[i] = univalue.GetTx()
		per.Paths[i] = univalue.GetPath()
		per.Reads[i] = univalue.GetReads()
		per.Writes[i] = univalue.GetWrites()
		per.AddOrDelete[i] = univalue.(*urltype.Univalue).AddOrDelete
		per.Composite[i] = univalue.(*urltype.Univalue).Composite
		switch v := univalue.GetValue().(type) {
		case *commutative.Balance:
			if v.GetDelta().Sign() >= 0 {
				continue
			}
			per.Transitions = append(per.Transitions, &BalanceTransition{
				Path:   univalue.GetPath(),
				Tx:     univalue.GetTx(),
				Origin: v.Value().(*big.Int),
				Delta:  v.GetDelta(),
			})
		}
	}
	return per
}
