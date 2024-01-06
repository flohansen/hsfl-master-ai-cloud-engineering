package _mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
)

type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) ValidateUserToken(ctx context.Context, in *user.ValidateUserTokenRequest, opts ...grpc.CallOption) (*user.ValidateUserTokenResponse, error) {
	args := m.Called(ctx, in)

	var resp *user.ValidateUserTokenResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*user.ValidateUserTokenResponse)
	}

	return resp, args.Error(1)
}
