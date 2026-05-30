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

func (g *MySQLAutoIncrementGenerator) NextCode(_ context.Context) (string, error) {
	return "", errors.New("mysql auto increment generator is not implemented")
}
