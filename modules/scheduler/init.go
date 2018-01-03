package scheduler

import (
	"errors"

	"github.com/edwardhey/flow/models"
	"github.com/edwardhey/flow/modules/nodes"
)

type Scheduler struct {
	Flow *models.Flow
}

func NewScheduler(flow *models.Flow) *Scheduler {
	f := &Scheduler{
		Flow: flow,
	}
	return f
}

func (s *Scheduler) NodeStart(n nodes.INode) {
	//调用k8s api
	err := n.Init()
	if err != nil {
		return
	}
	err = n.Run()
	return
}

func (s *Scheduler) NodeComplete(n *models.Node) error {
	return nil
}

func (s *Scheduler) Run() error {
	if s.Flow.Status != models.KindFlowStatusEndabled {
		return errors.New("flow disabled")
	}
	if s.Flow.ExecStatus >= models.KindFlowExecStatusRunning {
		return nil
	}
	s.Flow.ExecStatus = models.KindFlowExecStatusRunning

	//找到第一个节点，然后开始执行
	n, err := models.GetBeginNodeByFlow(s.Flow)
	if err != nil {
		return err
	}

	go s.NodeStart(n)
	return nil

}

func init() {
	// f := models.Flow
}
