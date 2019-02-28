package common

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/0987363/configGO/models"
	"github.com/0987363/configGO/service"
	"go.etcd.io/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterService(lnAddr string) {
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

	if client := ConnectEtcd(); client != nil {
		RegisterEtcd(client, etcdKey(ip, port), etcdValue(ip, port))
	}
}

func etcdKey(ip, port string) string {
	url := fmt.Sprintf("%s:%s", ip, port)
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

func RegisterEtcd(client *clientv3.Client, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), models.EtcdTimeout)
	_, err := client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}

	service.AddCloseHook(func() {
		UnRegisterEtcd(client, key)
	})

	return nil
}
