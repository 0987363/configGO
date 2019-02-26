package models

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"

	"golang.org/x/net/context"

	"errors"
	"strings"
	"time"
)

const (
	EtcdTimeout = 5 * time.Second
)

var client *clientv3.Client

func ConnectEtcd(address, ca, cert, key string) error {
	if ca != "" && cert != "" && key != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      cert,
			KeyFile:       key,
			TrustedCAFile: ca,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return errors.New("Init client config failed:" + err.Error())
		}

		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(address, ","),
			DialTimeout: EtcdTimeout,
			TLS:         tlsConfig,
		})
		if err != nil {
			return errors.New("New client tls failed:" + err.Error())
		}
		client = cli
		return nil
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(address, ","),
		DialTimeout: EtcdTimeout,
	})
	if err != nil {
		return errors.New("New client failed:" + err.Error())
	}
	client = cli
	return nil
}

func UnRegisterService(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdTimeout)
	_, err := client.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}

	return nil
}

func RegisterService(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdTimeout)
	_, err := client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}

	return nil
}
