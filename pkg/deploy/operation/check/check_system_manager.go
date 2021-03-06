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

package check

import (
	"bytes"
	"fmt"

	"github.com/kpaas-io/kpaas/pkg/deploy/command"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

type CheckSystemManagerOperation struct {
	shellCmd *command.ShellCommand
}

func (ckops *CheckSystemManagerOperation) RunCommands(config *pb.NodeCheckConfig, logChan chan<- *bytes.Buffer) (stdOut, stdErr []byte, err error) {

	itemBuffer := &bytes.Buffer{}

	m, err := machine.NewMachine(config.Node)
	if err != nil {
		return nil, nil, err
	}

	defer m.Close()

	ckops.shellCmd = command.NewShellCommand(m, "ps", "-p 1 | awk /1/'{print $4}'").
		WithDescription("检查机器 system manager 是否满足最低要求").
		WithExecuteLogWriter(itemBuffer)

	// run commands
	stdOut, stdErr, err = ckops.shellCmd.Execute()

	logChan <- itemBuffer

	return
}

// check is system manager is systemd
func CheckSystemManager(systemManager string, desireSysManager string) error {
	if systemManager != desireSysManager {
		return fmt.Errorf("system manager is not systemd")
	}
	return nil
}
