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
	loginHandler := setupLoginHandler()
	registerHandler := setUpRegisterHandler()

	userRepository := user.NewDemoRepository()
	userController := user.NewDefaultController(userRepository)
	userRouter := router.New(loginHandler, registerHandler, userController)

	if err := http.ListenAndServe(":3001", userRouter); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

func setupLoginHandler() *handler.LoginHandler {
	var jwtToken, _ = auth.NewJwtTokenGenerator(
		auth.JwtConfig{SignKey: "../../auth/test-token"})

	return handler.NewLoginHandler(setupMockRepository(),
		crypto.NewBcryptHasher(), jwtToken)
}

func setUpRegisterHandler() *handler.RegisterHandler {
	return handler.NewRegisterHandler(setupMockRepository(),
		crypto.NewBcryptHasher())
}

func setupMockRepository() user.Repository {
	repository := user.NewDemoRepository()
	userSlice := setupDemoUserSlice()
	for _, newUser := range userSlice {
		_, _ = repository.Create(newUser)
	}

	return repository
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
			Email:    "alan.turin@gmail.com",
			Password: hashedPassword,
			Name:     "Alan Turing",
			Role:     model.Customer,
		},
	}
}
