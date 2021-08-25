package workers

import (
	ethCommon "github.com/arcology/3rd-party/eth/common"
	"github.com/arcology/arbitrator-svc/types"
	ctypes "github.com/arcology/common-lib/types"
	"github.com/arcology/component-lib/actor"
	"github.com/arcology/component-lib/aggregator/aggregator"
	"github.com/arcology/component-lib/log"
	"go.uber.org/zap"
)

type EuResultsAggreSelector struct {
	actor.WorkerThread
	ag *aggregator.Aggregator
}

//return a Subscriber struct
func NewEuResultsAggreSelector(concurrency int, groupid string) *EuResultsAggreSelector {
	agg := EuResultsAggreSelector{}
	agg.Set(concurrency, groupid)
	agg.ag = aggregator.NewAggregator()
	return &agg
}

func (a *EuResultsAggreSelector) OnStart() {
}

func (a *EuResultsAggreSelector) OnMessageArrived(msgs []*actor.Message) error {
	switch msgs[0].Name {
	case actor.MsgAppHash:
		remainingQuantity := a.ag.OnClearInfoReceived()
		a.AddLog(log.LogLevel_Info, "clear pool", zap.Int("remainingQuantity", remainingQuantity))
	case actor.MsgReapinglist:
		reapinglist := msgs[0].Data.(*ctypes.ReapingList)
		result, _ := a.ag.OnListReceived(reapinglist)
		a.SendMsg(result)
	case actor.MsgInclusive:
		inclusive := msgs[0].Data.(*ctypes.InclusiveList)
		inclusive.Mode = ctypes.InclusiveMode_Results
		a.ag.OnClearListReceived(inclusive)
	case actor.MsgPreProcessedEuResults:
		data := msgs[0].Data.([]*types.ProcessedEuResult)
		if data != nil && len(data) > 0 {
			for _, v := range data {
				euResult := v
				result := a.ag.OnDataReceived(ethCommon.BytesToHash([]byte(euResult.Hash)), euResult)
				a.SendMsg(result)
			}
		}
	}
	return nil
}
func (a *EuResultsAggreSelector) SendMsg(selectedData *[]*interface{}) {
	if selectedData != nil {
		euResults := make([]*types.ProcessedEuResult, len(*selectedData))
		for i, euResult := range *selectedData {
			euResults[i] = (*euResult).(*types.ProcessedEuResult)
		}
		a.AddLog(log.LogLevel_CheckPoint, "send gather result", zap.Int("counts", len(euResults)))
		a.MsgBroker.Send(actor.MsgEuResultSelected, &euResults)
	}
}
