package mocks

import (
	"context"

	proto "github.com/Flo0807/hsfl-master-ai-cloud-engineering/lib/rpc/auth"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) ValidateToken(ctx context.Context, in *proto.ValidateTokenRequest, opts ...grpc.CallOption) (*proto.ValidateTokenResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*proto.ValidateTokenResponse), args.Error(1)
}
