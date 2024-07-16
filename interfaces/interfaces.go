package interfaces

import (
	"context"
	"repos/utils"
	"time"
)

type User struct {
	ID   uint
	Name string
}

type Filters struct {
	Offset  int
	Limit   int
	OrderBy string
	Gte     time.Time
	Lte     time.Time
}

type UsersRepo interface {
	GetById(context.Context, int64, ...utils.Options) (*User, error)
	GetAll(context.Context, Filters, ...utils.Options) ([]*User, error)
	Create(context.Context, *User, ...utils.Options) error
	Update(context.Context, *User, map[string]interface{}, ...utils.Options) error
	Delete(context.Context, int64, ...utils.Options) error
}
