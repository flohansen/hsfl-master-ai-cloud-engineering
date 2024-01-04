package rpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "hsfl.de/group6/hsfl-master-ai-cloud-engineering/lib/rpc/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"log"
	"strconv"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	userRepository *user.Repository
	hasher         crypto.Hasher
	tokenGenerator auth.TokenGenerator
}

func NewUserServiceServer(userRepository *user.Repository, tokenGenerator auth.TokenGenerator) *UserServiceServer {
	return &UserServiceServer{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *UserServiceServer) ValidateUserToken(_ context.Context, request *proto.ValidateUserTokenRequest) (*proto.ValidateUserTokenResponse, error) {
	claims, err := u.tokenGenerator.VerifyToken(request.Token)
	if err != nil {
		log.Println("Verification failed: ", err.Error())
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	id, err := strconv.ParseUint(claims["id"].(string), 10, 64)
	if err != nil {
		log.Println("Verification failed: Can't find user id in claim.")
		return nil, status.Error(codes.DataLoss, "Can't find user id in claim.")
	}

	user, err := (*u.userRepository).FindById(id)
	if err != nil {
		log.Println("Verification failed: ", err.Error())
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.ValidateUserTokenResponse{User: &proto.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
		Role:  int64(user.Role),
	}}, nil
}
