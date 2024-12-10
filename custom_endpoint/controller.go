package custom_endpoint

import "context"

type ServerAccessValidator func(ctx context.Context, header map[string]string) error
