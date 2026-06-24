package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sso/internal/config"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetDBDsn(cfg *config.Config) string {
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.USER, cfg.PG.PASSWORD, cfg.PG.HOST, cfg.PG.PORT, cfg.PG.NAME, cfg.PG.SSL)
	return dbDSN
}

func GetServerAddr(cfg *config.Config) string {
	addr := fmt.Sprintf("%s:%d", cfg.Server.HOST, cfg.Server.PORT)
	return addr
}

func RandomString(n int) (string, error) {
	result := make([]byte, n)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
