package controller

import (
	"auth/config"
	"auth/internal/infrastrucutre/tool"
	pb "auth/proto/auth"
	pbu "auth/proto/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	pbu.UserServiceClient
	auth *tool.JwtAuther
}

func NewAuthController(cc grpc.ClientConnInterface, cfg *config.Config) *AuthController {
	return &AuthController{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		UserServiceClient:              pbu.NewUserServiceClient(cc),
		auth:                           tool.NewJwtAuther(cfg),
	}
}
func (c *AuthController) Register(ctx context.Context, req *pb.User) (*pb.RegisterResponse, error) {
	userReq := pbu.User{
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	}
	_, err := c.UserServiceClient.SetUser(ctx, &userReq)
	if err != nil {
		return &pb.RegisterResponse{IsRegistred: false}, fmt.Errorf("Can't register user_service")
	}
	return &pb.RegisterResponse{IsRegistred: true}, nil
}
func (c *AuthController) Login(ctx context.Context, req *pb.User) (*pb.Token, error) {
	user, err := c.GetUser(ctx, &pbu.ProfileRequest{Email: req.Email, Phone: req.Phone})
	if err != nil {
		return &pb.Token{Token: ""}, fmt.Errorf("Unknown username or password")
	}
	if user.Password != req.Password {
		return &pb.Token{Token: ""}, fmt.Errorf("Unknown  password")
	}
	token, err := c.auth.GenerateToken(user.Email, user.Phone, 0)
	if err != nil {
		return &pb.Token{Token: ""}, fmt.Errorf("Internale error")
	}
	return &pb.Token{Token: token}, nil

}
func (c *AuthController) Authorised(ctx context.Context, req *pb.Token) (*pb.AuthorisedResponse, error) {
	token := req.GetToken()
	claim, isValid := c.auth.CheckToken(token, 0)
	if !isValid {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unauthorised")
	}
	usr := pbu.ProfileRequest{Email: claim.Email, Phone: claim.Phone}
	user, err := c.UserServiceClient.GetUser(ctx, &usr)
	if err != nil {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unauthorised")
	}
	if claim.Phone != user.GetPhone() && claim.Email != user.GetEmail() {
		return &pb.AuthorisedResponse{IsAuthorised: false}, fmt.Errorf("is unautorised")
	}
	return &pb.AuthorisedResponse{IsAuthorised: true}, nil
}
