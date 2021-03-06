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

package etcd

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pkiutil"

	"github.com/kpaas-io/kpaas/pkg/deploy/machine"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

// CreateAsCA creates a certificate authority, returning the created CA so it can be used to sign child certs.
func CreateAsCA(cfg *certutil.Config) (*x509.Certificate, crypto.Signer, error) {
	caCert, caKey, err := pkiutil.NewCertificateAuthority(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate %v CA certificate, errors: %v", cfg.CommonName, err)
	}

	return caCert, caKey, nil
}

func CreateFromCA(cfg *certutil.Config, caCrt *x509.Certificate, caKey crypto.Signer) (encodedKey, encodedCrt []byte, err error) {
	cert, key, err := pkiutil.NewCertAndKey(caCrt, caKey, cfg)
	if err != nil {
		return
	}

	encodedKey, err = keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		err = fmt.Errorf("failed to marshal private key to PEM, error: %v", err)
		return
	}

	encodedCrt = pkiutil.EncodeCertPEM(cert)

	return
}

func ToByte(crt *x509.Certificate, key crypto.Signer) (encodedKey, encodedCrt []byte, err error) {

	if key != nil {
		encodedKey, err = keyutil.MarshalPrivateKeyToPEM(key)
		if err != nil {
			err = fmt.Errorf("failed to to marshal private key to PEM, error: %v", err)
			return
		}
	}

	if crt != nil {
		encodedCrt = pkiutil.EncodeCertPEM(crt)
	}

	return
}

func FetchEtcdCertAndKey(etcdNode *pb.Node, baseName string) (*x509.Certificate, crypto.Signer, error) {
	certPath := fmt.Sprintf("%v/%v.crt", localEtcdCADir, baseName)
	keyPath := fmt.Sprintf("%v/%v.key", localEtcdCADir, baseName)

	localCert, err := os.Create(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create local %v cert path:%v, error:%v", baseName, certPath, err)
	}

	localKey, err := os.Create(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create local %v key path:%v, error:%v", baseName, keyPath, err)
	}

	m, err := machine.NewMachine(etcdNode)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create exec client for first etcd node:%v, error:%v", m.GetName(), err)
	}

	if err := m.FetchFile(localCert, DefaultEtcdCACertPath); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etcd %v cert, error:%v", baseName, err)
	}

	if err := m.FetchFile(localKey, defaultEtcdCAKeyPath); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch etcd %v key, error:%v", baseName, err)
	}

	etcdCACrt, etcCAKey, err := pkiutil.TryLoadCertAndKeyFromDisk(localEtcdCADir, baseName)
	if err != nil {
		err = fmt.Errorf("failed to load etcd %v cert and key from:%v, error:%v", baseName, localEtcdCADir, err)
		logrus.Error(err)
		return nil, nil, err
	}

	return etcdCACrt, etcCAKey, nil
}
