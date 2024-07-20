package interfaces

import (
	"context"
	"repos/utils"
	"time"
)

type User struct {
	ID        uint
	Name      string
	Age       uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderBy string

const (
	OrderByAge       OrderBy = "age"
	OrderByName      OrderBy = "name"
	OrderByCreatedAt OrderBy = "created_at"
)

type Filters struct {
	Offset       int
	Limit        int
	OrderBy      OrderBy
	CreatedAtGte time.Time
	CreatedAtLte time.Time
	AgeGte       uint8
	AgeLte       uint8
	IDs          []int64
}

type UsersRepo interface {
	GetById(context.Context, int64, ...utils.Options) (*User, error)
	GetAll(context.Context, Filters, ...utils.Options) ([]*User, int64, error)
	Create(context.Context, *User, ...utils.Options) error
	Update(context.Context, *User, map[string]interface{}, ...utils.Options) error
	Delete(context.Context, []int64, ...utils.Options) error
}
