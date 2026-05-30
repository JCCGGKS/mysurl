package codestrategy

import (
	"context"
	"errors"

	"mysurl1/internal/dao"
	"mysurl1/internal/utils"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeGenerator struct {
	node *snowflake.Node
	dao  *dao.ShortLinkDAO
}

func NewSnowflakeGenerator(workerID int64, shortLinkDAO *dao.ShortLinkDAO) (*SnowflakeGenerator, error) {
	node, err := snowflake.NewNode(workerID)
	if err != nil {
		return nil, err
	}

	return &SnowflakeGenerator{
		node: node,
		dao:  shortLinkDAO,
	}, nil
}

func (g *SnowflakeGenerator) Provider() string {
	return ProviderSnowflake
}

func (g *SnowflakeGenerator) NextCode(ctx context.Context, originalURL, urlHash string) (string, error) {
	if g == nil || g.node == nil {
		return "", errors.New("snowflake generator node is not configured")
	}
	if g.dao == nil {
		return "", errors.New("snowflake generator dao is not configured")
	}

	id := g.node.Generate().Int64()
	if id <= 0 {
		return "", errors.New("snowflake generated invalid id")
	}

	shortCode := utils.EncodeBase62(uint64(id))
	return g.dao.CreateWithShortCode(ctx, shortCode, originalURL, urlHash)
}
