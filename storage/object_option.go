package storage

import "context"

// ObjectOption is optional parameters used for object call.
type ObjectOption struct {
	BucketName string
	Path       string
	MaxResult  int64
	Projection string
	Context    context.Context
}

func (o ObjectOption) getOrCreateContext() context.Context {
	if o.Context != nil {
		return o.Context
	}
	return context.Background()
}
