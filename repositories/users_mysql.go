package repositories

import (
	"context"
	"time"

	"repos/interfaces"
	"repos/utils"

	"gorm.io/gorm"
)

type userRepoMysql struct {
	db *gorm.DB
}

func NewUserRepoMysql(db *gorm.DB) interfaces.UsersRepo {
	return &userRepoMysql{db: db}
}

func (r userRepoMysql) GetById(ctx context.Context, id int64, opts ...utils.Options) (*interfaces.User, error) {
	var user *interfaces.User
	err := utils.ConfigureDB(r.db, opts...).WithContext(ctx).Find(&user, id).Error
	return user, err
}

func (r userRepoMysql) GetAll(ctx context.Context, filters interfaces.Filters, opts ...utils.Options) ([]*interfaces.User, error) {
	var users []*interfaces.User
	q := utils.ConfigureDB(r.db, opts...).WithContext(ctx)

	var limit = utils.Limit
	if filters.Limit != 0 {
		limit = filters.Limit
	}

	var offset = 0
	if filters.Offset != 0 {
		offset = filters.Offset
	}

	var gte = time.Date(2024, 4, 17, 0, 0, 0, 0, time.Local).AddDate(0, 0, -30)
	if !filters.Gte.IsZero() {
		gte = filters.Gte
	}

	var lte = time.Date(2024, 4, 17, 0, 0, 0, 0, time.Local)
	if !filters.Lte.IsZero() {
		lte = filters.Lte
	}

	err := q.
		Where("created_at <= ?", lte).
		Where("created_at >= ?", gte).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).
		Error

	return users, err
}

func (r userRepoMysql) Create(ctx context.Context, user *interfaces.User, opts ...utils.Options) error {
	return utils.ConfigureDB(r.db, opts...).WithContext(ctx).Create(user).Error
}

func (r userRepoMysql) Update(ctx context.Context, user *interfaces.User, vals map[string]interface{}, opts ...utils.Options) error {
	return utils.ConfigureDB(r.db, opts...).WithContext(ctx).Model(&user).Updates(vals).Error
}

func (r userRepoMysql) Delete(ctx context.Context, id int64, opts ...utils.Options) error {
	return utils.ConfigureDB(r.db, opts...).WithContext(ctx).Delete(&interfaces.User{}, id).Error
}
