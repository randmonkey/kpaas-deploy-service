// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package task

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/action"
	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// nodeCheckProcessor implements the specific logic for the node check task.
type nodeCheckProcessor struct {
}

func (p *nodeCheckProcessor) StartTask(task Task) error {
	nodeCheckTask, ok := task.(*nodeCheckTask)
	if !ok {
		return fmt.Errorf("%s: %T", consts.MsgTaskTypeMismatched, task)
	}

	go p.run(nodeCheckTask)

	return nil
}

func (p *nodeCheckProcessor) run(nodeCheckTask *nodeCheckTask) {
	logrus.Debugf("Run node check task: %+v", *nodeCheckTask)

	nodeCheckTask.status = TaskDoing

	// first, split task
	actions, err := p.splitTask(nodeCheckTask)
	if err != nil {
		nodeCheckTask.status = TaskFailed
		nodeCheckTask.err = &pb.Error{
			Reason:     consts.MsgTaskSplitFailed,
			Detail:     err.Error(),
			FixMethods: consts.MsgUnknownFixMethod,
		}
		return
	}

	var wg sync.WaitGroup

	// execute the actions concurrently
	for _, act := range actions {
		wg.Add(1)
		go action.ExecuteAction(act, &wg)
	}

	wg.Wait()

	if err = p.genTaskSummary(nodeCheckTask); err != nil {
		nodeCheckTask.status = TaskFailed
		nodeCheckTask.err = &pb.Error{
			Reason:     consts.MsgTaskGenSummaryFailed,
			Detail:     err.Error(),
			FixMethods: consts.MsgUnknownFixMethod,
		}

		logrus.Error(consts.MsgTaskGenSummaryFailed)
	}

	logrus.Debugf("Finish node check task: %+v", *nodeCheckTask)
}

// Spilt the task into one or more node check actions
func (p *nodeCheckProcessor) splitTask(task *nodeCheckTask) ([]action.Action, error) {
	logrus.Debugf("Start to spolit task")

	if task == nil {
		return nil, fmt.Errorf(consts.MsgEmptyTask)
	}

	task.status = TaskSplitting

	// split task into actions: will create a action for every node, the action type
	// is NodeCheckAction
	actions := make([]action.Action, len(task.nodeConfigs))
	for _, subConfig := range task.nodeConfigs {
		actionCfg := &action.NodeCheckActionConfig{
			NodeCheckConfig: subConfig,
			LogFileBasePath: task.logFilePath,
		}
		act, err := action.NewNodeCheckAction(actionCfg)
		if err != nil {
			return nil, err
		}
		actions = append(actions, act)
	}

	logrus.Debugf("Task has been splitted into %d actions", len(actions))

	// update task status to "splitted"
	task.status = TaskSplitted

	return actions, nil
}

// sum task status
func (p *nodeCheckProcessor) genTaskSummary(t *nodeCheckTask) error {
	logrus.Debugf("Start to gen task summary")

	if t == nil {
		return fmt.Errorf(consts.MsgEmptyTask)
	}

	done := 0
	failed := 0
	// combined error message
	errMsgs := make([]string, 0)
	for _, act := range t.actions {
		switch act.GetStatus() {
		case action.ActionFailed:
			failed++
			errMsgs = append(errMsgs, fmt.Sprintf("%v", act.GetErr()))
		case action.ActionDone:
			done++
		}
	}

	// if any action is failed, the task is failed
	if failed > 0 {
		t.status = TaskFailed
		t.err = &pb.Error{
			Reason:     "one or more checks failed",
			Detail:     fmt.Sprintf("%v", errMsgs),
			FixMethods: "check the detail mssage",
		}
	} else if done == len(t.actions) {
		// if all actions are done, the task is done
		t.status = TaskDone
	}

	return nil
}