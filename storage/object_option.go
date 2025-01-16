package storage

import "context"

// ObjectOption is optional parameters used for object call.
type ObjectOption struct {
	Context    context.Context
	BucketName string
	Path       string

	CacheControl string
}

func (o ObjectOption) getOrCreateContext() context.Context {
	if o.Context != nil {
		return o.Context
	}
	return context.Background()
}
