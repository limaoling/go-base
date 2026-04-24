package database

import (
	"database/sql"
	"fmt"
	"go-base/config"
	"go-base/model"
	"go-base/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDatabase 初始化数据库（创建数据库、连接、迁移表结构）
func InitDatabase() (*gorm.DB, error) {
	// 1. 创建数据库（如果不存在）
	if err := createDatabaseIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// 2. 连接数据库
	db, err := connectDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 3. 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// 4. 初始化表注释
	if err := model.InitTableComments(db); err != nil {
		utils.Warnf("Warning: failed to init table comments: %v", err)
	}

	// 5. 配置连接池
	if err := setupConnectionPool(db); err != nil {
		return nil, fmt.Errorf("failed to setup connection pool: %w", err)
	}

	return db, nil
}

// createDatabaseIfNotExists 创建数据库（如果不存在）
func createDatabaseIfNotExists() error {
	dsnWithoutDB := config.GlobalConfig.Database.GetDSNWithoutDB()
	dbTemp, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer dbTemp.Close()

	// 检查数据库是否存在
	var exists int
	err = dbTemp.QueryRow(
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?",
		config.GlobalConfig.Database.DBName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database: %w", err)
	}

	// 如果数据库不存在，则创建
	if exists == 0 {
		query := fmt.Sprintf(
			"CREATE DATABASE `%s` DEFAULT CHARACTER SET %s COLLATE %s_general_ci",
			config.GlobalConfig.Database.DBName,
			config.GlobalConfig.Database.Charset,
			config.GlobalConfig.Database.Charset,
		)
		if _, err = dbTemp.Exec(query); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		utils.Infof("✅ Database '%s' created successfully", config.GlobalConfig.Database.DBName)
	} else {
		utils.Infof("✅ Database '%s' already exists", config.GlobalConfig.Database.DBName)
	}

	return nil
}

// connectDatabase 连接数据库
func connectDatabase() (*gorm.DB, error) {
	dsn := config.GlobalConfig.Database.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 准备语句设置
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Ping 测试连接
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	utils.Infof("✅ Database connected successfully")
	return db, nil
}

// autoMigrate 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	utils.Infof("✅ Database tables migrated successfully")
	return nil
}

// setupConnectionPool 配置连接池
func setupConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	dbConfig := config.GlobalConfig.Database

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)                                    // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second) // 连接最大生命周期

	utils.Infof("✅ Connection pool configured (MaxIdle: %d, MaxOpen: %d, MaxLifetime: %ds)",
		dbConfig.MaxIdleConns,
		dbConfig.MaxOpenConns,
		dbConfig.ConnMaxLifetime,
	)

	return nil
}
