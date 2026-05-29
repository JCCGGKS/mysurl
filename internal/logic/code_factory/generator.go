package codefactory

import "context"

type CodeGenerator interface {
	Provider() string
	NextCode(ctx context.Context) (string, error)
}
