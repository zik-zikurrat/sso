package helper

import (
	"fmt"
	"sso/internal/config"
)

func GetDBDsn(cfg *config.Config) string {
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.USER, cfg.PG.PASSWORD, cfg.PG.HOST, cfg.PG.PORT, cfg.PG.NAME, cfg.PG.SSL)
	return dbDSN
}

func GetServerAddr(cfg *config.Config) string {
	addr := fmt.Sprintf("%s:%d", cfg.Server.HOST, cfg.Server.PORT)
	return addr
}
