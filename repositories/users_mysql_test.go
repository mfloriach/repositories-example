package repositories_test

import (
	"context"
	"fmt"
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

	db, err := gorm.Open(mysql.Open(mysqlContainer.GetConnection(ctx)), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		t.Fatalf("mounting db: %v", err)
	}

	r := repositories.NewUserRepoMysql(db)

	t.Run("check limit to 2", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 2})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
	})

	t.Run("check limit to 6", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 6, len(users), "they should be equal")
	})

	t.Run("offset beging", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(6), users[0].ID, "they should be equal")
	})

	t.Run("offset move", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{Offset: 1, Limit: 6})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(5), users[0].ID, "they should be equal")
	})

	t.Run("lte test", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{Lte: time.Date(2024, 4, 13, 0, 0, 0, 0, time.Local)})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(3), users[0].ID, "they should be equal")
	})

	t.Run("lte default", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint(6), users[0].ID, "they should be equal")
	})

	// t.Run("gte test", func(t *testing.T) {
	// 	users, err := r.GetAll(ctx, interfaces.Filters{Gte: time.Date(2024, 4, 14, 0, 0, 0, 0, time.Local)})
	// 	if err != nil {
	// 		t.Fatalf("creating user table: %v", err)
	// 	}

	// 	assert.Equal(t, 2, len(users), "they should be equal")
	// })

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
	}); err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	users, err := r.GetAll(ctx, interfaces.Filters{Offset: 0, Limit: 10})
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	fmt.Println(users[0].ID)
	fmt.Println(len(users))
}
