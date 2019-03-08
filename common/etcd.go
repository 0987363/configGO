package common

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/0987363/configGO/models"
	"github.com/0987363/viper"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

const (
	baseKey = "/github.com/0987363/configGO"
)

func Registry() chan *Service {
	client := ConnectEtcd()
	if client == nil {
		log.Fatal("Connect to etcd failed.")
		return nil
	}
	c := make(chan *Service, 10)

	go func() {
		for {
			app := <-c
			key := app.Key()
			log.Info("Start update: ", key)

			switch app.Op {
			case watcher.Write, watcher.Create:
				if err := RegisterEtcd(client, key, app.Value); err != nil {
					log.Error("Register service failed:", err)
				}
				continue
			case watcher.Remove:
				if err := UnRegisterEtcd(client, key); err != nil {
					log.Error("Unregister etcd failed:", err)
				}
				continue
			default:
				log.Error("Unsupport op:", app.Op)
			}
		}
	}()

	return c
}

func (app *Service) Key() string {
	url := fmt.Sprintf("%s/%s", app.Project, app.Service)
	//	url := fmt.Sprintf("%s:%s-%s", ip, port, uuid.NewV4().String())
	return filepath.Join(baseKey, url)
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

func RegisterEtcd(client *clientv3.Client, key string, value string, op ...clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), models.EtcdTimeout)
	_, err := client.Put(ctx, key, value, op...)
	cancel()
	if err != nil {
		return err
	}

	return nil
}
