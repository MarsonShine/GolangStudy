package fzlog

import "context"

func NewContext(ctx context.Context) FzLog {
	if ctx == nil {
		return CreateLog()
	} else {
		return withContext(ctx)
	}
}
