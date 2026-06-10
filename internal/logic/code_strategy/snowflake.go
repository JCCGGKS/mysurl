package codestrategy

import (
	"context"
	"errors"

	"mysurl1/internal/utils"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator(workerID int64) (*SnowflakeGenerator, error) {
	node, err := snowflake.NewNode(workerID)
	if err != nil {
		return nil, err
	}

	return &SnowflakeGenerator{
		node: node,
	}, nil
}

func (g *SnowflakeGenerator) Provider() string {
	return ProviderSnowflake
}

func (g *SnowflakeGenerator) NextCode(ctx context.Context) (string, error) {
	if g == nil || g.node == nil {
		return "", errors.New("snowflake generator node is not configured")
	}

	id := g.node.Generate().Int64()
	if id <= 0 {
		return "", errors.New("snowflake generated invalid id")
	}

	shortCode := utils.EncodeBase62(uint64(id))
	return shortCode, nil
}
