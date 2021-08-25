package workers

import (
	"context"
	"errors"

	ethCommon "github.com/arcology/3rd-party/eth/common"
	"github.com/arcology/arbitrator-svc/accumulator"
	arbitrator "github.com/arcology/arbitrator-svc/impl-arbitrator"
	"github.com/arcology/arbitrator-svc/types"
	ctypes "github.com/arcology/common-lib/types"
	"github.com/arcology/component-lib/actor"
	kafkalib "github.com/arcology/component-lib/kafka/lib"
	"github.com/arcology/component-lib/log"
	"go.uber.org/zap"
)

type RpcService struct {
	actor.WorkerThread
	wbs        *kafkalib.Waitobjs
	msgid      int64
	arbitrator *arbitrator.ArbitratorImpl
}

//return a Subscriber struct
func NewRpcService(lanes int, groupid string) *RpcService {
	rs := RpcService{
		arbitrator: arbitrator.NewArbitratorImpl(),
	}
	rs.Set(lanes, groupid)
	rs.msgid = 0
	return &rs
}

func (rs *RpcService) OnStart() {
	rs.wbs = kafkalib.StartWaitObjects()
}

func (rs *RpcService) OnMessageArrived(msgs []*actor.Message) error {
	var euResults *[]*types.ProcessedEuResult
	for _, v := range msgs {
		switch v.Name {
		case actor.MsgStartSub:

		case actor.MsgEuResultSelected:
			euResults = v.Data.(*[]*types.ProcessedEuResult)
			rs.AddLog(log.LogLevel_Debug, "received euresult***********", zap.Int64("msgid", rs.msgid))
			rs.wbs.Update(rs.msgid, euResults)
		}
	}

	return nil
}

func (rs *RpcService) Arbitrate(ctx context.Context, request *actor.Message, response *ctypes.ArbitratorResponse) error {
	lstMessage := request.CopyHeader()
	rs.ChangeEnvironment(lstMessage)
	params := request.Data.(*ctypes.ArbitratorRequest)
	list := []*ethCommon.Hash{}
	for _, rows := range params.TxsListGroup {
		for _, element := range rows {
			list = append(list, element.TxHash)
		}
	}
	reapinglist := ctypes.ReapingList{
		List: list,
	}

	rs.msgid = rs.msgid + 1
	rs.AddLog(log.LogLevel_Debug, "start arbitrate request***********", zap.Int("txs", len(reapinglist.List)))
	rs.wbs.AddWaiter(rs.msgid)
	rs.MsgBroker.Send(actor.MsgReapinglist, &reapinglist)

	rs.wbs.Waitforever(rs.msgid)
	results := rs.wbs.GetData(rs.msgid)

	var resultSelected *[]*types.ProcessedEuResult
	if results == nil {
		rs.AddLog(log.LogLevel_Error, "select euresults error")
		return errors.New("select euresults error")
	}

	if bValue, ok := results.(*[]*types.ProcessedEuResult); ok {
		resultSelected = bValue
	} else {
		rs.AddLog(log.LogLevel_Error, "select euresults type error")
		return errors.New("select euresults type error")
	}

	if resultSelected != nil && len(*resultSelected) > 0 {
		logid := rs.AddLog(log.LogLevel_Info, "Before NewExecutionSchedule")
		interLog := rs.GetLogger(logid)
		conflictedList, left, right := detectConflict(rs.arbitrator, params.TxsListGroup, resultSelected, interLog)
		rs.AddLog(log.LogLevel_Debug, "arbitrate return results***********", zap.Int("txResults", len(conflictedList)))
		response.ConflictedList = conflictedList
		response.CPairLeft = left
		response.CPairRight = right
		return nil
	}

	return nil
}

func detectConflict(arbitrator *arbitrator.ArbitratorImpl, txsListGroup [][]*ctypes.TxElement, euResults *[]*types.ProcessedEuResult, inlog *actor.WorkerThreadLogger) ([]*ethCommon.Hash, []uint32, []uint32) {
	// Make the results indexable.
	euDict := make(map[ethCommon.Hash]*types.ProcessedEuResult, len(*euResults))
	for _, r := range *euResults {
		euDict[ethCommon.BytesToHash([]byte(r.Hash))] = r
	}
	// Prepare arguments for DetectConflict.
	groups := make([][]*types.ProcessedEuResult, 0, len(txsListGroup))
	var maxBatch uint64 = 0
	for _, g := range txsListGroup {
		group := make([]*types.ProcessedEuResult, 0, len(g))
		for _, e := range g {
			group = append(group, euDict[*e.TxHash])
			if e.Batchid > maxBatch {
				maxBatch = e.Batchid
			}
		}
		groups = append(groups, group)
	}
	// Arbitration.
	arbitrator.Reset()
	ids, _, flags, left, right := arbitrator.DetectConflict(groups)
	// Unique conflict IDs.
	uniqueConflicts := make(map[uint32]struct{})
	for i, conflict := range flags {
		if conflict {
			uniqueConflicts[ids[i]] = struct{}{}
		}
	}
	// Add conflicted groups into conflict list.
	var conflictedList []*ethCommon.Hash
	for id := range uniqueConflicts {
		for _, e := range txsListGroup[id] {
			conflictedList = append(conflictedList, e.TxHash)
		}
	}
	inlog.Log(log.LogLevel_Debug, "Arbitration result", zap.Int("conflictNums", len(conflictedList)))
	// Make batch info indexable.
	batches := make([][]*types.ProcessedEuResult, maxBatch+1)
	for i := range batches {
		batches[i] = make([]*types.ProcessedEuResult, 0, len(*euResults)/(i+1))
	}
	for i, g := range txsListGroup {
		if _, ok := uniqueConflicts[uint32(i)]; ok {
			continue
		}
		for _, e := range g {
			batches[e.Batchid] = append(batches[e.Batchid], euDict[*e.TxHash])
		}
	}
	// Prepare arguments for BatchCheck.
	txs := make([]*types.ProcessedEuResult, 0, len(*euResults))
	for i := uint64(0); i <= maxBatch; i++ {
		txs = append(txs, batches[i]...)
	}
	// Accumulation.
	results := accumulator.CheckBalance(batches)
	// Make results indexable.
	conflictDict := make(map[ethCommon.Hash]struct{})
	for i, conflict := range results {
		if !conflict {
			continue
		}
		conflictDict[ethCommon.BytesToHash([]byte(txs[i].Hash))] = struct{}{}
	}
	// Add conflicted groups into conflict list.
	for i, g := range txsListGroup {
		if _, ok := uniqueConflicts[uint32(i)]; ok {
			continue
		}
		isConflict := false
		for _, e := range g {
			// If any of the member in this group is conflicted, the whole group is conflicted.
			if _, ok := conflictDict[*e.TxHash]; ok {
				isConflict = true
				break
			}
		}
		if isConflict {
			for _, e := range g {
				conflictedList = append(conflictedList, e.TxHash)
			}
		}
		//inlog.Log(log.LogLevel_Debug, "Accumulation result", zap.Int("conflictNums", len(conflictedList)))
	}

	return conflictedList, left, right
}
