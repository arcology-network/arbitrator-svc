package types

import (
	"sync"

	urltype "github.com/arcology-network/concurrenturl/v2/type"
)

var processedEuResultPool sync.Pool
var univaluePool sync.Pool

func init() {
	processedEuResultPool = sync.Pool{
		New: func() interface{} {
			return &ProcessedEuResult{}
		},
	}
	univaluePool = sync.Pool{
		New: func() interface{} {
			return &urltype.Univalue{}
		},
	}
}
