package repositories_test

import (
	"context"
	"log"
	"runtime"
	"testing"
	"time"

	"repos/interfaces"
	"repos/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	testMysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type testContainerMysql struct {
	mysqlContainer *testMysql.MySQLContainer
}

func NewTestContainerMysql(ctx context.Context) (testContainerMysql, func(ctx context.Context), error) {
	image := "mysql:8.0"
	if runtime.GOARCH == "arm64" {
		image = "arm64v8/mysql:8.0"
	}

	mysqlContainer, err := testMysql.RunContainer(ctx,
		testcontainers.WithImage(image),
		testMysql.WithScripts("main.sql"),
	)

	t := testContainerMysql{mysqlContainer: mysqlContainer}

	return t, t.cleanDB, err
}

func (tcm testContainerMysql) GetConnection(ctx context.Context) string {
	return tcm.mysqlContainer.MustConnectionString(ctx) + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (tcm testContainerMysql) cleanDB(ctx context.Context) {
	if err := tcm.mysqlContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func TestUserMysqlRepoGetByID(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, close, err := NewTestContainerMysql(ctx)
	if err != nil {
		t.Fatalf("mounting db container: %v", err)
	}
	defer close(ctx)

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	user, err := r.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "first", user.Name, "they should be equal")
}

func TestUserMysqlRepoGetAll(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, close, err := NewTestContainerMysql(ctx)
	if err != nil {
		t.Fatalf("mounting db container: %v", err)
	}
	defer close(ctx)

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	t.Run("check limit to 2", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 2})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("check limit to 6", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 6, len(users), "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("offset beging", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(6), users[0].ID, "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("offset move", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{Offset: 1, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(5), users[0].ID, "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("lte test", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtLte: time.Date(2024, 4, 13, 0, 0, 0, 0, time.Local),
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(3), users[0].ID, "they should be equal")
		assert.Equal(t, int64(3), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("lte default", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(6), users[0].ID, "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("gte test", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtGte: time.Date(2024, 4, 14, 0, 0, 0, 0, time.Local),
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
		assert.Equal(t, int64(2), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("between gte and lte", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtGte: time.Date(2024, 4, 13, 0, 0, 0, 0, time.Local),
			CreatedAtLte: time.Date(2024, 4, 15, 0, 0, 0, 0, time.Local),
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
		assert.Equal(t, int64(2), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("order by age", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			OrderBy: interfaces.OrderByAge,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint8(66), users[0].Age, "they should be equal")
		assert.Equal(t, uint8(55), users[1].Age, "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("order by name", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			OrderBy: interfaces.OrderByName,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, "third", users[0].Name, "they should be equal")
		assert.Equal(t, "six", users[1].Name, "they should be equal")
		assert.Equal(t, int64(6), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("age range", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			AgeGte: 22,
			AgeLte: 45,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 4, len(users), "they should be equal")
		assert.Equal(t, int64(4), total, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})

	t.Run("by ids", func(t *testing.T) {
		users, total, err := r.GetAll(ctx, interfaces.Filters{
			IDs: []int64{3, 4},
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
		assert.Equal(t, int64(2), total, "they should be equal")
		assert.Equal(t, uint(3), users[0].ID, "they should be equal")
		assert.Equal(t, uint(4), users[1].ID, "they should be equal")
		assert.Nil(t, err, "they should be equal")
	})
}

func TestUserMysqlRepoCreate(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, close, err := NewTestContainerMysql(ctx)
	if err != nil {
		t.Fatalf("mounting db container: %v", err)
	}
	defer close(ctx)

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	if err := r.Create(ctx, &interfaces.User{
		Name: "John Doe",
		Age:  5,
	}); err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	users, _, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 10})
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, uint(6), users[0].ID, "they should be equal")
}

func TestUserMysqlRepoDelete(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, close, err := NewTestContainerMysql(ctx)
	if err != nil {
		t.Fatalf("mounting db container: %v", err)
	}
	defer close(ctx)

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	t.Run("single", func(t *testing.T) {
		u := &interfaces.User{
			Name: "John Doe",
		}

		if err := r.Create(ctx, u); err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		user, err := r.GetById(ctx, int64(u.ID))
		if err != nil {
			t.Fatalf("get user: %v", err)
		}

		assert.Equal(t, uint(7), user.ID, "they should be equal")

		if err := r.Delete(ctx, []int64{int64(u.ID)}); err != nil {
			t.Fatalf("deleting user table: %v", err)
		}

		userDelete, err := r.GetById(ctx, int64(u.ID))
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Empty(t, userDelete, "they should be equal")
		assert.Nil(t, err)
	})

	t.Run("multiple", func(t *testing.T) {
		u1 := &interfaces.User{
			Name: "John Doe",
		}

		if err := r.Create(ctx, u1); err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		u2 := &interfaces.User{
			Name: "John Doe",
		}

		if err := r.Create(ctx, u2); err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		if err := r.Delete(ctx, []int64{int64(u1.ID), int64(u2.ID)}); err != nil {
			t.Fatalf("deleting user table: %v", err)
		}

		userDelete1, err := r.GetById(ctx, int64(u1.ID))
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Empty(t, userDelete1, "they should be equal")
		assert.Nil(t, err)

		userDelete2, err := r.GetById(ctx, int64(u1.ID))
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Empty(t, userDelete2, "they should be equal")
		assert.Nil(t, err)
	})

}

func TestUserMysqlRepoUpdate(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, close, err := NewTestContainerMysql(ctx)
	if err != nil {
		t.Fatalf("mounting db container: %v", err)
	}
	defer close(ctx)

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	u := interfaces.User{
		Name: "John Doe",
		Age:  5,
	}

	if err := r.Create(ctx, &u); err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	user, err := r.GetById(ctx, 7)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "John Doe", user.Name, "they should be equal")

	val := map[string]interface{}{
		"name": "new name",
	}

	if err := r.Update(ctx, &u, val); err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	userUpdated, err := r.GetById(ctx, 7)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "new name", userUpdated.Name, "they should be equal")
}
