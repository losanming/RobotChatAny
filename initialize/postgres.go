package initialize

import (
	"fmt"
	_ "github.com/lib/pq"
	"new-robot-self/config"
	"new-robot-self/model"
	"new-robot-self/pkg/postgres"
	"new-robot-self/public"
)

func InitPgsql() error {
	var c postgres.PostgresConfig
	conf := config.GetConfig().Postgres
	c = postgres.PostgresConfig{
		Host:            conf.Host,
		Port:            conf.Port,
		User:            conf.User,
		Password:        conf.Password,
		DBName:          conf.DBName,
		ConnMaxLifetime: conf.ConnMaxLifetime,
		MaxIdleConns:    conf.MaxIdleConns,
		MaxOpenConns:    conf.MaxOpenConns,
	}
	db, err := postgres.InitPostgres(c)
	public.DB = db
	
	// 同步所有表结构
	tables := []interface{}{
		&model.User{},
		&model.HealthRecord{},
		&model.TodoItem{},
		&model.CalendarEvent{},
		&model.DietRecord{},
		&model.UserPreference{},
		&model.Reminder{},
	}
	
	if err := postgres.SyncTables(public.DB, tables); err != nil {
		return fmt.Errorf("同步表结构失败: %w", err)
	}
	
	fmt.Printf("pgsql 初始化配置成功，已同步 %d 个表\n", len(tables))
	return err
}

func SyncTables(table []interface{}) error {
	return postgres.SyncTables(public.DB, table)
}

func ClosePgsql() {
	public.DB.Close()
}
