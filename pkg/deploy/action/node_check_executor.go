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

package action

import (
	"fmt"

	"github.com/sirupsen/logrus"

	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type nodeCheckExecutor struct {
	config *pb.NodeCheckConfig
}

func (a *nodeCheckExecutor) Execute(act Action) error {
	logrus.Debugf("Start to execute action: %+v", act)

	nodeCheckAction, ok := act.(*nodeCheckAction)
	if !ok {
		return fmt.Errorf("the action type is not match: should be node check action, but is %T", act)
	}

	// The following codes are just for example
	// 1. check kernel version
	// TODO: ssh to check the kernel version
	kernelVersionItem := &nodeCheckItem{
		name:        "kernel version",
		description: "kernel version must not less than 4.4.0",
		status:      nodeCheckItemSucessful,
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, kernelVersionItem)

	// 2. check docker version
	// TODO: ssh to check the docker version
	dockerVersionItem := &nodeCheckItem{
		name:        "docker version check",
		description: "docker version must not less than 17.03.02",
		status:      nodeCheckItemFailed,
		err: &pb.Error{
			Reason:     "docker version is too low: 17.0.0",
			Detail:     "",
			FixMethods: "upgrade docker version to 17.03.02+",
		},
	}
	nodeCheckAction.checkItems = append(nodeCheckAction.checkItems, dockerVersionItem)

	// TODO: other checks

	// TODO: update action status
	nodeCheckAction.status = ActionFailed
	nodeCheckAction.err = dockerVersionItem.err

	logrus.Debugf("End to execute action: %+v", act)
	return nil
}