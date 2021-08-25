package workers

import (
	"github.com/HPISTechnologies/component-lib/actor"
	"github.com/HPISTechnologies/component-lib/log"
	"go.uber.org/zap"
)

type Initializer struct {
	actor.WorkerThread
}

//return a Subscriber struct
func NewInitializer(concurrency int, groupid string) *Initializer {
	in := Initializer{}
	in.Set(concurrency, groupid)

	return &in
}

func (i *Initializer) OnStart() {
}

func (i *Initializer) Stop() {

}

func (i *Initializer) OnMessageArrived(msgs []*actor.Message) error {
	i.AddLog(log.LogLevel_Info, "arbitrator initialize ", zap.String("send command", actor.MsgStartSub))
	i.MsgBroker.Send(actor.MsgStartSub, "")
	//i.MsgBroker.Send(actor.MsgBlockCompleted, actor.MsgBlockCompleted_Success)
	return nil
}
