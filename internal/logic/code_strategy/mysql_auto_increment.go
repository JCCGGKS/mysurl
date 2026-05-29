package codestrategy

import (
	"context"
	"errors"
)

type MySQLAutoIncrementGenerator struct{}


func (g *MySQLAutoIncrementGenerator) Provider() string {
	return ProviderMySQLAutoIncrement
}

func (g *MySQLAutoIncrementGenerator) NextCode(_ context.Context, _ NextCodeInput) (string, error) {
	return "", errors.New("mysql auto increment generator is not implemented")
}
