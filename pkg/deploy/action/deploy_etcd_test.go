// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	certutil "k8s.io/client-go/util/cert"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	"github.com/kpaas-io/kpaas/pkg/deploy/operation/etcd"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

func init() {
	machine.IsTesting = true
}

func TestDeployEtcd(t *testing.T) {
	executor := new(deployEtcdExecutor)

	cacert, signer, err := etcd.CreateAsCA(&certutil.Config{
		CommonName: "test",
	})
	assert.NoError(t, err)

	normalAction, err := NewDeployEtcdAction(&DeployEtcdActionConfig{
		CaCrt: cacert,
		CaKey: signer,
		Node: &pb.Node{
			Name: "normal",
			Ip:   "10.10.10.10",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, normalAction)

	pbErr := executor.Execute(normalAction)
	assert.Nil(t, pbErr)

	errorAction, err := NewDeployEtcdAction(&DeployEtcdActionConfig{
		CaCrt: cacert,
		CaKey: signer,
		Node: &pb.Node{
			Name: "error",
			Ip:   "10.10.10.10",
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, errorAction)

	pbErr = executor.Execute(errorAction)
	assert.NoError(t, err)
	assert.NotNil(t, pbErr)
}
