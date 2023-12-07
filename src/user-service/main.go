package main

import (
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/api/router"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/auth"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user"
	"hsfl.de/group6/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"log"
	"net/http"
)

func main() {
	userRepository := user.NewDemoRepository()

	loginHandler := setupLoginHandler(userRepository)
	registerHandler := setUpRegisterHandler(userRepository)

	userController := user.NewDefaultController(userRepository)
	userRouter := router.New(loginHandler, registerHandler, userController)

	if err := http.ListenAndServe(":3001", userRouter); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

func setupLoginHandler(userRepository user.Repository) *handler.LoginHandler {
	var jwtToken, _ = auth.NewJwtTokenGenerator(
		auth.JwtConfig{SignKey: "../../auth/test-token"})

	return handler.NewLoginHandler(setupMockRepository(userRepository),
		crypto.NewBcryptHasher(), jwtToken)
}

func setUpRegisterHandler(userRepository user.Repository) *handler.RegisterHandler {
	return handler.NewRegisterHandler(setupMockRepository(userRepository),
		crypto.NewBcryptHasher())
}

func setupMockRepository(userRepository user.Repository) user.Repository {
	userSlice := setupDemoUserSlice()
	for _, newUser := range userSlice {
		_, _ = userRepository.Create(newUser)
	}

	return userRepository
}

func setupDemoUserSlice() []*model.User {
	bcryptHasher := crypto.NewBcryptHasher()
	hashedPassword, _ := bcryptHasher.Hash([]byte("123456"))

	return []*model.User{
		{
			Id:       1,
			Email:    "ada.lovelace@gmail.com",
			Password: hashedPassword,
			Name:     "Ada Lovelace",
			Role:     model.Customer,
		},
		{
			Id:       2,
			Email:    "info-aldi@gmail.com",
			Password: hashedPassword,
			Name:     "Aldi",
			Role:     model.Merchant,
		},
		{
			Id:       3,
			Email:    "info-edeka@gmail.com",
			Password: hashedPassword,
			Name:     "Edeka",
			Role:     model.Merchant,
		},
	}
}
