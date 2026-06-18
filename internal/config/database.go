package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// OpenMySQL 根据配置创建 GORM 数据库连接。
//
// 你可以把 GORM 理解成 Go 里的 MyBatis-Plus/JPA：
// 它帮我们把结构体和数据库表关联起来，少写很多重复 SQL。
func OpenMySQL(cfg MySQLConfig) (*gorm.DB, error) {
	// DSN 是 Data Source Name，也就是数据库连接字符串。
	// 格式大致是: 用户名:密码@tcp(主机:端口)/数据库名?参数
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// gorm.Open 只负责创建数据库连接对象；后续 Repository 会拿这个 db 去查询表。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	return db, nil
}
