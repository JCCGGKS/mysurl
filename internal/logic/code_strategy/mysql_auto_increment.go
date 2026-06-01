package codestrategy

import (
	"context"
	"errors"

	"mysurl1/internal/dao"
)

type MySQLAutoIncrementGenerator struct {
	dao *dao.ShortLinkDAO
}

func NewMySQLAutoIncrementGenerator(shortLinkDAO *dao.ShortLinkDAO) *MySQLAutoIncrementGenerator {
	return &MySQLAutoIncrementGenerator{dao: shortLinkDAO}
}

func (g *MySQLAutoIncrementGenerator) Provider() string {
	return ProviderMySQLAutoIncrement
}

func (g *MySQLAutoIncrementGenerator) NextCode(ctx context.Context, normalizedURL, urlHash string) (string, error) {
	if g == nil || g.dao == nil {
		return "", errors.New("mysql auto increment generator dao is not configured")
	}

	return g.dao.CreateWithAutoIncrement(ctx, normalizedURL, urlHash)
}
