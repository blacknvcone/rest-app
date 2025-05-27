package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"rest-app/config"
)

type GormDB struct {
	*gorm.DB
}

type DbConfig struct {
	GormDB *GormDB
	Pool   *pgxpool.Pool
	SQLDB  *sql.DB
}

func (db *DbConfig) CloseConnection() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.SQLDB != nil {
		db.SQLDB.Close()
	}
}

// Init initializes GORM and pgxpool separately.
func Init(gormDSN, pgxDSN string) (*DbConfig, error) {
	var (
		dbConfigVar DbConfig
		loggerGorm  logger.Interface
	)
	configData := config.GetConfig()

	// Configure logger
	if configData.App.Env == "local" {
		loggerGorm = logger.Default.LogMode(logger.Info)
	} else {
		loggerGorm = logger.Default.LogMode(logger.Silent)
	}

	// ✅ Initialize GORM with *sql.DB
	sqlDB, err := sql.Open("postgres", gormDSN)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(configData.DB.MaxOpenConn)
	sqlDB.SetMaxIdleConns(configData.DB.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(configData.DB.MaxLifetimeConn))

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger:                 loggerGorm,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	// ✅ Initialize pgxpool separately
	pgxConfig, err := pgxpool.ParseConfig(pgxDSN)
	if err != nil {
		return nil, err
	}
	pgxConfig.MaxConns = int32(configData.DB.MaxOpenConn)
	pgxConfig.MaxConnIdleTime = time.Second * time.Duration(configData.DB.MaxIdletimeConn)
	pgxConfig.HealthCheckPeriod = time.Second * 15

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	dbConfigVar.GormDB = &GormDB{gormDB}
	dbConfigVar.SQLDB = sqlDB
	dbConfigVar.Pool = pool

	log.Println("Database successfully connected")

	return &dbConfigVar, nil
}
