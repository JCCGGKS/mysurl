package codefactory

import "fmt"

const (
	ProviderMySQLAutoIncrement = "mysql_auto_increment"
	ProviderRedisIncr          = "redis_incr"
	ProviderSnowflake          = "snowflake"
)

func New(provider string) (CodeGenerator, error) {
	switch provider {
	case ProviderMySQLAutoIncrement:
		return NewMySQLAutoIncrementGenerator(), nil
	case ProviderRedisIncr:
		return NewRedisIncrGenerator(), nil
	case ProviderSnowflake:
		return NewSnowflakeGenerator(), nil
	default:
		return nil, fmt.Errorf("unsupported short code provider: %s", provider)
	}
}
