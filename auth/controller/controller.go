package controller

import (
	"auth/jwt"
	pb "auth/proto/auth"
	pbu "auth/proto/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	pbu.UserServiceClient
	auth *jwt.JwtAuther
}

func NewAuthController(cc grpc.ClientConnInterface) *AuthController {
	return &AuthController{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		UserServiceClient:              pbu.NewUserServiceClient(cc),
		auth:                           jwt.NewJwtAuther(),
	}
}
func (c *AuthController) Register(ctx context.Context, req *pb.User) (*pb.RegisterResponse, error) {
	_, err := c.SetUser(ctx, &pbu.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return &pb.RegisterResponse{Response: "Cant't register user_service"}, fmt.Errorf("Can't register user_service")
	}
	return &pb.RegisterResponse{Response: "Sucsess"}, fmt.Errorf("Sucseful register")
}
func (c *AuthController) Login(ctx context.Context, req *pb.User) (*pb.Token, error) {
	user, err := c.GetUser(ctx, &pbu.ProfileRequest{Email: req.Email})
	if err != nil {
		return &pb.Token{Token: ""}, fmt.Errorf("Unknow username or password")
	}

	token, err := c.auth.GenerateToken(user.Email, user.Password)
	if err != nil {
		return &pb.Token{Token: ""}, fmt.Errorf("Internale error")
	}
	return &pb.Token{Token: token}, nil

}
func (c *AuthController) Authorised(ctx context.Context, req *pb.Token) (*pb.AuthorisedResponse, error) {
	token := req.GetToken()
	email, _, isValid := c.auth.CheckToken(token)
	if !isValid {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unauthorised")
	}

	user, err := c.UserServiceClient.GetUser(ctx, &pbu.ProfileRequest{Email: email})
	if err != nil {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unauthorised")
	}

	if user.GetEmail() != email {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unauthorised")
	}
	return &pb.AuthorisedResponse{IsAuthorised: true}, nil
}
