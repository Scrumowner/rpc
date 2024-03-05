package controller

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user/internal/models"
	"user/internal/modules/service"
	pb "user/proto"
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

func (c *UserController) List(ctx context.Context, req *pb.EmptyRequest) (*pb.Users, error) {
	users, err := c.service.GetList(ctx)
	if err != nil {
		return &pb.Users{}, err
	}
	res := pb.Users{}
	for _, user := range users {
		res.User = append(res.User, &pb.User{
			Email:    user.Email,
			Phone:    user.Phone,
			Password: user.Password,
		})
	}
	return &res, nil
}
func (c *UserController) GetUser(ctx context.Context, req *pb.ProfileRequest) (*pb.User, error) {
	us := models.User{
		Email: req.GetEmail(),
		Phone: req.GetPhone(),
	}
	user, err := c.service.GetUser(&us)
	if err != nil {
		return &pb.User{}, err
	}
	return &pb.User{
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}, nil

}
func (c *UserController) SetUser(ctx context.Context, req *pb.User) (*pb.EmptyRequest, error) {
	user := models.User{
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	}
	err := c.service.SetUser(&user)
	if err != nil {
		return &pb.EmptyRequest{}, err
	}
	return &pb.EmptyRequest{}, nil
}
