package models

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"

	"errors"
	"strings"
	"time"
)

const (
	EtcdTimeout = 5 * time.Second
)

func ConnectEtcd(address, ca, cert, key string) (*clientv3.Client, error) {
	if ca != "" && cert != "" && key != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      cert,
			KeyFile:       key,
			TrustedCAFile: ca,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, errors.New("Init client config failed:" + err.Error())
		}

		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(address, ","),
			DialTimeout: EtcdTimeout,
			TLS:         tlsConfig,
		})
		if err != nil {
			return nil, errors.New("New client tls failed:" + err.Error())
		}
		return cli, nil
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(address, ","),
		DialTimeout: EtcdTimeout,
	})
	if err != nil {
		return nil, errors.New("New client failed:" + err.Error())
	}
	return cli, nil
}
