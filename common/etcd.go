package common

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/0987363/configGO/models"
	"github.com/0987363/configGO/service"
	log "github.com/sirupsen/logrus"
	"github.com/0987363/viper"
	"go.etcd.io/etcd/clientv3"
)

func Registry(lnAddr string) {
	address := viper.GetString("address")
	addrs := strings.Split(address, ":")
	port := addrs[1]
	if port == "0" {
		rAddrs := strings.Split(lnAddr, ":")
		port = rAddrs[len(rAddrs)-1]
	}

	ip := addrs[0]
	if addrs[0] == "0.0.0.0" || addrs[0] == "" {
		ip = models.GetHostIP()
	}

	client := ConnectEtcd()
	if client == nil {
		return
	}

	resp, err := client.Grant(context.TODO(), 60)
	if err != nil {
		log.Fatal("Grant etcd failed:", err)
	}
	key := etcdKey(ip, port)
	value := etcdValue(ip, port)
	if err := RegisterEtcd(client, key, value, clientv3.WithLease(resp.ID)); err != nil {
		log.Fatal("Register service failed:", err)
	}

	m := ReadConfig([]string{"test2", "pikachu"}...)
	if m != nil {
		res, _ := json.Marshal(m)
		if err := RegisterEtcd(client, filepath.Join(key, "cattle/pikachu"), string(res), clientv3.WithLease(resp.ID)); err != nil {
			log.Fatal("Register service failed:", err)
		}
	}

	idleConnsClosed := make(chan struct{})
	service.AddCloseHook(func() {
		defer client.Close()
		close(idleConnsClosed)
		if err := UnRegisterEtcd(client, key); err != nil {
			log.Error("Unregister etcd failed:", err)
		}
	})

	go func() {
		for {
			select {
			case <-idleConnsClosed:
				return
			case <-time.After(20 * time.Second):
				_, err := client.KeepAliveOnce(context.TODO(), resp.ID)
				if err != nil {
					log.Error("Keepalive etcd service failed:", err)
				}
			}
		}
	}()
}

func etcdKey(ip, port string) string {
	url := fmt.Sprintf("%s:%s", ip, port)
	//	url := fmt.Sprintf("%s:%s-%s", ip, port, uuid.NewV4().String())
	return filepath.Join(viper.GetString("etcd.url"), url)
}

func etcdValue(ip, port string) string {
	return fmt.Sprintf("http://%s:%s", ip, port)
}

func ConnectEtcd() *clientv3.Client {
	address := viper.GetString("etcd.address")
	ca := viper.GetString("etcd.ca")
	cert := viper.GetString("etcd.cert")
	key := viper.GetString("etcd.key")

	client, err := models.ConnectEtcd(address, ca, cert, key)
	if err != nil {
		log.Error("Connect to etcd failed:", err)
		return nil
	}

	return client
}

func UnRegisterEtcd(client *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), models.EtcdTimeout)
	_, err := client.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}

	return nil
}

func RegisterEtcd(client *clientv3.Client, key string, value string, op clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), models.EtcdTimeout)
	_, err := client.Put(ctx, key, value, op)
	cancel()
	if err != nil {
		return err
	}

	return nil
}
