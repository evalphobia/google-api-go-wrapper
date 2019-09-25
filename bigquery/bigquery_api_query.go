package bigquery

func (b *BigQuery) Query(opt QueryOption) (*QueryResponse, error) {
	resp, err := b.RunQuery(opt.ToRequest())
	if err != nil {
		return nil, err
	}
	return &QueryResponse{resp}, nil
}
