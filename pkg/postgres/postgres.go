package postgres

import (
	"fmt"
	"time"
	"xorm.io/xorm"
)

// InitPostgres 初始化PostgreSQL连接
// 使用xorm作为ORM框架
// 返回数据库实例和可能的错误
type PostgresConfig struct {
	// Host 数据库主机地址
	Host string `yaml:"host"`
	// Port 数据库端口
	Port int `yaml:"port"`
	// User 数据库用户名
	User string `yaml:"user"`
	// Password 数据库密码
	Password string `yaml:"password"`
	// DBName 数据库名称
	DBName string `yaml:"dbname"`
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns"`
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns"`
	// ConnMaxLifetime 连接最大生命周期（秒）
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
}

func InitPostgres(c PostgresConfig) (*xorm.Engine, error) {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)

	// 创建数据库实例
	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %v", err)
	}

	// 设置连接池参数
	engine.SetMaxIdleConns(c.MaxIdleConns)
	engine.SetMaxOpenConns(c.MaxOpenConns)
	engine.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime))

	// 测试数据库连接
	if err = engine.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置日志级别
	engine.ShowSQL(true)

	// 设置时区为本地时区
	engine.DatabaseTZ = time.Local
	engine.TZLocation = time.Local
	return engine, nil
}

// SyncTables 同步所有表结构
// 如果表不存在会创建表，已存在则同步结构
func SyncTables(engine *xorm.Engine, table []interface{}) error {
	if engine == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 在这里添加需要同步的表结构
	// 注意：按照依赖关系顺序添加
	var tables []interface{}
	tables = append(tables, table...)

	// 同步所有表结构
	for _, table := range tables {
		if err := engine.Sync2(table); err != nil {
			return fmt.Errorf("同步表结构失败 %T: %v", table, err)
		}
	}

	return nil
}
