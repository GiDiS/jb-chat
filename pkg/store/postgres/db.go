package postgres

import (
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/logger"
	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	DbDriver   string `env:"DB_DRIVER" envDefault:"postgres"`
	DbHost     string `env:"DB_HOST,required"`
	DbPort     int    `env:"DB_PORT,required" envDefault:"5432"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`
}

func MustGetConfig(log logger.Logger) Config {
	config := Config{
		DbDriver:   "postgres",
		DbHost:     "db",
		DbPort:     5432,
		DbUser:     "",
		DbPassword: "",
		DbName:     "jb_chat",
	}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	return config
}

func (cfg Config) Dsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName,
	)
}

func (cfg Config) Driver() string {
	return cfg.DbDriver
}

func ConnectToDBWithConfig(cfg Config, log logger.Logger) *sqlx.DB {
	return ConnectToDB(cfg.Driver(), cfg.Dsn(), log)
}

func ConnectToDB(driver, dsn string, log logger.Logger) *sqlx.DB {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	// Периодический пинг с переустановкой связи
	//go keepAlive(db, time.Second*60, log)

	if db != nil {
		db.SetMaxIdleConns(30)
		db.SetMaxOpenConns(50)
	}

	return db
}

func keepAlive(db *sqlx.DB, timeout time.Duration, log logger.Logger) {
	for {
		time.Sleep(timeout)
		if err := db.Ping(); err != nil {
			log.Errorf("Db pinger error: %v", err)
		}
	}
}
