package service

import (
	v1 "app/api/user/v1"
	"app/internal/biz"
	"app/internal/dto"
	"context"
)

type UserService struct {
	u *biz.FindUser
	v1.UnimplementedUserServer
}

func NewUserService(u *biz.FindUser) v1.UserServer {
	return &UserService{u: u}
}

func (s *UserService) Find(ctx context.Context, r *v1.FindRequest) (*v1.FindReply, error) {
	u := &dto.User{Id: r.Id}
	Nickname := s.u.FindUser(u)
	return &v1.FindReply{Nickname: Nickname}, nil
}
