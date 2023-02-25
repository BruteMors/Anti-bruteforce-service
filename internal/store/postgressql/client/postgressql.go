package client

import (
	"Anti-bruteforce-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresSql struct {
	Db     *sqlx.DB
	logger *zap.SugaredLogger
	config *config.Config
}

func NewPostgresSql(logger *zap.SugaredLogger, config *config.Config) *PostgresSql {
	return &PostgresSql{logger: logger, config: config}
}

func (p *PostgresSql) Open() error {
	dbSourceName := "host=" + p.config.Database.Host + " " + "dbname=" + p.config.Database.DbName + " " + "port=" + p.config.Database.Port + " " + "user=" + p.config.Database.User + " " + "password=" + p.config.Database.Password + " " + "sslmode=" + p.config.Database.SslMode
	db, err := sqlx.Open("postgres", dbSourceName)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	p.Db = db
	p.logger.Info("Connection to db successfully")
	return nil
}

func (p *PostgresSql) Close() error {
	err := p.Db.Close()
	if err != nil {
		return err
	}
	p.logger.Info("Close db successfully")
	return nil
}
