// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"app/internal/biz"
	"app/internal/data"
)

// Injectors from wire.go:

func InitFindUser() *biz.FindUser {
	userRepo := data.NewUserRepo()
	FindUser := biz.NewFindUser(userRepo)
	return FindUser
}