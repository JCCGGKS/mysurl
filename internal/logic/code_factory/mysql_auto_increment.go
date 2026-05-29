package codefactory

import (
	"context"
	"errors"
)

type MySQLAutoIncrementGenerator struct{}

func NewMySQLAutoIncrementGenerator() *MySQLAutoIncrementGenerator {
	return &MySQLAutoIncrementGenerator{}
}

func (g *MySQLAutoIncrementGenerator) Provider() string {
	return ProviderMySQLAutoIncrement
}

func (g *MySQLAutoIncrementGenerator) NextCode(_ context.Context) (string, error) {
	return "", errors.New("mysql auto increment generator is not implemented")
}
