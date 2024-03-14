package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/shop-framework/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func WithMysqlInfo(address string, hand func(cli *gorm.DB) error) error {
	err := config.ViperInit(address)
	if err != nil {
		return err
	}

	var app struct {
		MysqlConf struct {
			Username string
			Password string
			Host     string
			Port     string
			Database string
		} `json:"Mysql"`
	}

	id := viper.GetString("Database.DataId")
	Group := viper.GetString("Database.Group")

	res, err := config.GetConfig(id, Group)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(res), &app); err != nil {
		return err
	}

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		app.MysqlConf.Username,
		app.MysqlConf.Password,
		app.MysqlConf.Host,
		app.MysqlConf.Port,
		app.MysqlConf.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	d, _ := db.DB()
	defer d.Close()

	if err = hand(db); err != nil {
		return err
	}
	return nil
}
