// Copyright 2020 Shanghai JingDuo Information Technology co., Ltd.
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
	"github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	RegisterExecutor(ActionTypeDeployWorker, new(deployWorkerExecutor))
}

type deployWorkerExecutor struct {
	deployNodeExecutor
}

func (executor *deployWorkerExecutor) Execute(act Action) *protos.Error {

	action, ok := act.(*DeployWorkerAction)
	if !ok {
		return errOfTypeMismatched(new(DeployWorkerAction), act)
	}

	return executor.deployNodeExecutor.Deploy(act, action.config)
}
