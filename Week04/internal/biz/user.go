package biz

import (
	"app/internal/dto"
)

type UserRepo interface {
	Query(*dto.User) string
}

func NewFindUser(repo UserRepo) *FindUser {
	return &FindUser{repo: repo}
}

type FindUser struct {
	repo UserRepo
}

func (s *FindUser) FindUser(u *dto.User) string {
	return s.repo.Query(u)
}
