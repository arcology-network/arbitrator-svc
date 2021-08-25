package workers

import (
	"math"
	"math/big"

	"github.com/HPISTechnologies/evm/common"
	"github.com/HPISTechnologies/evm/consensus"
	"github.com/HPISTechnologies/evm/core/types"
	"github.com/HPISTechnologies/evm/core/vm"
	"github.com/HPISTechnologies/evm/params"
	adaptor "github.com/HPISTechnologies/vm-adaptor/evm"
)

type fakeChain struct {
}

func (chain *fakeChain) GetHeader(common.Hash, uint64) *types.Header {
	return &types.Header{}
}

func (chain *fakeChain) Engine() consensus.Engine {
	return nil
}

var coinbase = common.BytesToAddress([]byte{100, 100, 100})

func createTestConfig() *adaptor.Config {
	vmConfig := vm.Config{}
	cfg := &adaptor.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: new(big.Int).SetUint64(10000000),
		ParentHash:  common.Hash{},
		Time:        new(big.Int).SetUint64(10000000),
		Coinbase:    &coinbase,
		GasLimit:    math.MaxUint64,
		Difficulty:  new(big.Int).SetUint64(10000000),
	}
	cfg.Chain = new(fakeChain)
	return cfg
}

func createTestConfigOnBlock(bn uint64) *adaptor.Config {
	cfg := createTestConfig()
	cfg.BlockNumber = new(big.Int).SetUint64(bn)
	// 60 seconds per block.
	cfg.Time = new(big.Int).SetUint64(bn * 60)
	return cfg
}
