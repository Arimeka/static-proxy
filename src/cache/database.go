package cache

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"

	"fmt"
)

func NewDBConn() (*gorm.DB, error) {
	v := viper.New()
	v.SetEnvPrefix("cachedb")
	v.BindEnv("adapter")
	v.SetDefault("adapter", "sqlite3")
	v.BindEnv("conn")
	v.SetDefault("conn", "tmp/cache.db")
	v.BindEnv("pool")
	v.SetDefault("pool", 5)
	v.BindEnv("max_conn")
	v.SetDefault("max_conn", 15)
	v.BindEnv("debug")
	v.SetDefault("debug", true)

	config := &Database{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("Invalid cache database config: %v", err)
	}

	db, err := gorm.Open(config.Adapter, config.ConnString)
	if err != nil {
		return db, err
	}

	if err = db.DB().Ping(); err != nil {
		return db, err
	}
	db.DB().SetMaxIdleConns(config.Pool)
	db.DB().SetMaxOpenConns(config.MaxConn)

	db.LogMode(config.Debug)

	db.AutoMigrate(&File{})
	db.Model(&File{}).AddUniqueIndex("idx_file_name", "filename")
	db.Model(&File{}).AddIndex("idx_file_deleted_updated", "deleted","updated_at")

	//if err := testData(db); err != nil {
	//	return db, err
	//}

	return db, nil
}

type Database struct {
	Adapter    string `mapstructure:"adapter"`
	ConnString string `mapstructure:"conn"`
	Pool       int    `mapstructure:"pool"`
	MaxConn    int    `mapstructure:"max_conn"`
	Debug      bool   `mapstructure:"debug"`
}


func testData(db *gorm.DB) error {
	file1 := &File{
		Filename: "cache/saimoe-2012-12.jpg",
		ContentType: "image/jpeg",
		Size: 263961,
	}

	if err := db.Create(file1).Error; err != nil {
		return err
	}
	return nil
}
