package parser

import (
	"context"
	parserpb "parser/gen/pb"
)

type server struct{}

func (s *server) GetCurrentBlock(ctx context.Context, in *parserpb.GetCurrentBlockRequest) (resp *parserpb.GetCurrentBlockResponse, err error) {
	resp = &parserpb.GetCurrentBlockResponse{
		Id: int64(getParserInstance().GetCurrentBlock()),
	}
	return resp, nil
}
