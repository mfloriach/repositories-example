package repositories

import (
	"context"
	"fmt"
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

func (r userRepoMysql) GetAll(ctx context.Context, filters interfaces.Filters, opts ...utils.Options) ([]*interfaces.User, int64, error) {
	var users []*interfaces.User

	var offset = 0
	if filters.Offset != 0 {
		offset = filters.Offset
	}

	var limit = utils.Limit
	if filters.Limit != 0 {
		limit = filters.Limit
	}

	var orderBy = interfaces.OrderByCreatedAt
	if filters.OrderBy != "" {
		orderBy = filters.OrderBy
	}

	now := time.Date(2024, 4, 17, 0, 0, 0, 0, time.Local)

	var createdAtGte time.Time = now.AddDate(0, 0, -30)
	if !filters.CreatedAtGte.IsZero() {
		createdAtGte = filters.CreatedAtGte
	}

	var createdAtLte time.Time = now
	if !filters.CreatedAtLte.IsZero() {
		createdAtLte = filters.CreatedAtLte
	}

	stmp := utils.
		ConfigureDB(r.db, opts...).
		WithContext(ctx).
		Debug().
		Where("created_at <= ?", createdAtLte).
		Where("created_at >= ?", createdAtGte)

	if filters.AgeGte != 0 {
		stmp = stmp.Where("age >= ?", filters.AgeGte)
	}

	if filters.AgeLte != 0 {
		stmp = stmp.Where("age <= ?", filters.AgeLte)
	}

	var total int64
	if err := stmp.Count(&total).Error; err != nil {
		return users, 0, err
	}

	err := stmp.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%v DESC", orderBy)).
		Find(&users).
		Error

	return users, total, err
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
