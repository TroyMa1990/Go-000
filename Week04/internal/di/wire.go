//+build wireinject 忽略编译

package di

import (
	"app/internal/biz"
	"app/internal/data"

	"github.com/google/wire"
)

func InitFindUser() *biz.FindUser {
	wire.Build(biz.NewFindUser, data.NewUserRepo)
	return &biz.FindUser{}
}
