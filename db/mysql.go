package db

import (
	"github.com/RollNA/harbour/zLog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
)

type MysqlConf struct {
	DSN                string `json:"DSN" mapstructure:"DSN"`
	DSNRead            string `json:"DSNRead" mapstructure:"DSNRead"`
	MaxConnections     int    `json:"MaxConnections" mapstructure:"MaxConnection"`
	MaxIdleConnections int    `json:"MaxIdleConnections" mapstructure:"MaxIdleConnections"`
}

var Client *gorm.DB

func InitMysql(conf MysqlConf) {
	var err error
	Client, err = gorm.Open(mysql.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		zLog.Fatal("mysql open failed", zap.Error(err))
	}

	err = Client.Use(dbresolver.Register(dbresolver.Config{
		// use `db2` as sources, `db3`, `db4` as replicas
		Sources:  []gorm.Dialector{mysql.Open(conf.DSN)},
		Replicas: []gorm.Dialector{mysql.Open(conf.DSNRead)},
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		zLog.Fatal("mysql register failed", zap.Error(err))
	}
	if err = Client.Use(tracing.NewPlugin()); err != nil {
		zLog.Fatal("mysql add trace plugin failed", zap.Error(err))
	}

	sqlDB, _ := Client.DB()
	sqlDB.SetMaxOpenConns(conf.MaxConnections)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConnections)
}
