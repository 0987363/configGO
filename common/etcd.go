package common

import (
	"github.com/0987363/configGO/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegisterService(s string) {
	address := viper.GetString("address")
	url := viper.GetString("etcd.url")
	ca := viper.GetString("etcd.ca")
	cert := viper.GetString("etcd.cert")
	key := viper.GetString("etcd.key")

	if err := models.ConnectEtcd(address, ca, cert, key); err != nil {
		log.Error("Connect to etcd failed:", err)
		return
	}

	if err := models.RegisterService(url, s); err != nil {
		log.Error("Register service failed:", err)
		return
	}
}
