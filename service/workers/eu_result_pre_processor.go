package workers

import (
	"github.com/arcology-network/arbitrator-svc/types"
	"github.com/arcology-network/common-lib/common"
	ctypes "github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/component-lib/actor"
)

type EuResultPreProcessor struct {
	actor.WorkerThread
}

func NewEuResultPreProcessor(concurrency int, groupid string) *EuResultPreProcessor {
	p := &EuResultPreProcessor{}
	p.Set(concurrency, groupid)
	return p
}

func (p *EuResultPreProcessor) OnStart() {

}

func (p *EuResultPreProcessor) OnMessageArrived(msgs []*actor.Message) error {
	// t0 := time.Now()
	results := *(msgs[0].Data.(*ctypes.TxAccessRecordses))
	processed := make([]*types.ProcessedEuResult, len(results))
	worker := func(start, end, idx int, args ...interface{}) {
		euresults := args[0].([]interface{})[0].(ctypes.TxAccessRecordses)
		processResults := args[0].([]interface{})[1].(*[]*types.ProcessedEuResult)
		for i := start; i < end; i++ {
			(*processResults)[i] = types.Process(euresults[i])
		}
	}
	common.ParallelWorker(len(results), p.Concurrency, worker, results, &processed)
	p.MsgBroker.Send(actor.MsgPreProcessedEuResults, processed)
	// fmt.Println("EuResultPreProcessor: ", time.Since(t0))
	return nil
}
