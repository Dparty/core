package services

import (
	"fmt"

	"github.com/Dparty/common/cloud"
	"github.com/Dparty/feieyun"
	"github.com/Dparty/model"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

var DB *gorm.DB

var CosClient cloud.CosClient
var Bucket string
var BillPrinter feieyun.Printer

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	user := viper.GetString("feieyun.user")
	ukey := viper.GetString("feieyun.ukey")
	url := viper.GetString("feieyun.url")
	BillPrinter = feieyun.NewPrinter(user, ukey, url)
}

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	Bucket = viper.GetString("cos.Bucket")
	CosClient.Region = viper.GetString("cos.Region")
	CosClient.SecretID = viper.GetString("cos.SecretID")
	CosClient.SecretKey = viper.GetString("cos.SecretKey")
}

func Init(inject *gorm.DB) {
	DB = inject
	model.Init(DB)
}
