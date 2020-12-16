package data

import (
	"app/internal/biz"
	"app/internal/dto"
)

var _ biz.UserRepo = (biz.UserRepo)(nil)

func NewUserRepo() biz.UserRepo {
	return &userRepo{}
}

type userRepo struct {
	Id int64
}

func (r *userRepo) Query(u *dto.User) string {
	return "nickname"
}
