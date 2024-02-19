package controller

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user/models"
	pb "user/proto/user"
	"user/service"
)

type UserController struct {
	service *service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserController(db *sqlx.DB) *UserController {
	return &UserController{
		service:                        service.NewUserService(db),
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
	}
}

func (c *UserController) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.User, error) {
	user, err := c.service.GetUser(req.GetEmail())
	if err != nil {
		return &pb.User{}, err
	}
	return &pb.User{Email: user.GetEmail(), Password: user.GetPassword()}, nil
}
func (c *UserController) List(ctx context.Context, req *pb.EmptyRequest) (*pb.ListResponse, error) {
	_ = *req
	users, err := c.service.GetList()
	if err != nil {
		return &pb.ListResponse{}, err
	}
	var resp pb.ListResponse
	for _, user := range users {
		resp.User = append(resp.User, &pb.User{Email: user.GetEmail(), Password: user.GetPassword()})
	}
	return &resp, nil
}
func (c *UserController) GetUser(ctx context.Context, req *pb.ProfileRequest) (*pb.User, error) {
	user, err := c.service.GetUser(req.GetEmail())
	if err != nil {
		return &pb.User{}, err
	}
	return &pb.User{Email: user.GetEmail(), Password: user.GetPassword()}, nil
}
func (c *UserController) SetUser(ctx context.Context, req *pb.User) (*pb.EmptyRequest, error) {
	err := c.service.SetUser(models.User{Email: req.GetEmail(), Password: req.GetPassword()})
	if err != nil {
		return &pb.EmptyRequest{}, err
	}
	return &pb.EmptyRequest{}, nil
}
