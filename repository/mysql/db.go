package mysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	for {
		err = db.DB.Ping()
		if err == nil {
			break
		}
		log.Info().Msg("connecting mysql server")
		time.Sleep(time.Second * 2)
	}
	log.Info().Msg("conntected mysql server")

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	return db, nil
}
