package repositories_test

import (
	"context"
	"log"
	"testing"
	"time"

	"repos/interfaces"
	"repos/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	testMysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type testContainerMongo struct {
	container *testMysql.MySQLContainer
}

func NewTestContainerMongo(ctx context.Context) (testContainerMongo, func(ctx context.Context), error) {
	container, err := testMysql.RunContainer(ctx,
		testcontainers.WithImage("mongo:7.0.5"),
		testMysql.WithScripts("main.sql"),
	)

	t := testContainerMongo{container: container}

	return t, t.cleanDB, err
}

func (tcm testContainerMongo) GetConnection(ctx context.Context) string {
	return tcm.container.MustConnectionString(ctx) + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (tcm testContainerMongo) cleanDB(ctx context.Context) {
	if err := tcm.container.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func TestUserMongoRepoGetByID(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0.5")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = mongo.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := mongo.Database("test")

	collection := db.Collection("users")
	if _, err := collection.InsertOne(ctx, interfaces.User{ID: 1, Name: "asdsa", Age: 12}); err != nil {
		panic(err)
	}

	r := repositories.NewUserRepoMongo(db)

	user, err := r.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "asdsa", user.Name, "they should be equal")
}

func TestUserMongoRepoGetAll(t *testing.T) {
	t.Skip()

	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0.5")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = mongo.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := mongo.Database("test")

	r := repositories.NewUserRepoMongo(db)
	collection := db.Collection("users")

	for i := 0; i < 7; i++ {
		if _, err := collection.InsertOne(ctx, interfaces.User{
			ID:        uint(i),
			Name:      "asdsa",
			Age:       10 + uint8(i),
			CreatedAt: time.Date(2024, 4, 10+i, 0, 0, 0, 0, time.Local),
		}); err != nil {
			panic(err)
		}
	}

	docs := []interface{}{
		interfaces.User{
			ID:        1,
			Name:      "first",
			Age:       55,
			CreatedAt: time.Date(2024, 4, 10, 23, 0, 0, 0, time.Local),
		},
		interfaces.User{
			ID:        2,
			Name:      "second",
			Age:       22,
			CreatedAt: time.Date(2024, 4, 11, 23, 0, 0, 0, time.Local),
		},
		interfaces.User{
			ID:        3,
			Name:      "third",
			Age:       40,
			CreatedAt: time.Date(2024, 4, 12, 23, 0, 0, 0, time.Local),
		},
		interfaces.User{
			ID:        4,
			Name:      "forth",
			Age:       30,
			CreatedAt: time.Date(2024, 4, 13, 23, 0, 0, 0, time.Local),
		},
		interfaces.User{
			ID:        5,
			Name:      "five",
			Age:       45,
			CreatedAt: time.Date(2024, 4, 14, 23, 0, 0, 0, time.Local),
		},
		interfaces.User{
			ID:        6,
			Name:      "six",
			Age:       66,
			CreatedAt: time.Date(2024, 4, 15, 23, 0, 0, 0, time.Local),
		},
	}
	if _, err := collection.InsertMany(ctx, docs); err != nil {
		panic(err)
	}

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
		users, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtLte: time.Date(2024, 4, 13, 0, 0, 0, 0, time.Local),
		})
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

	t.Run("gte test", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtGte: time.Date(2024, 4, 14, 0, 0, 0, 0, time.Local),
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
	})

	t.Run("between gte and lte", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{
			CreatedAtGte: time.Date(2024, 4, 13, 0, 0, 0, 0, time.Local),
			CreatedAtLte: time.Date(2024, 4, 15, 0, 0, 0, 0, time.Local),
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 2, len(users), "they should be equal")
	})

	t.Run("order by age", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{
			OrderBy: interfaces.OrderByAge,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, uint8(66), users[0].Age, "they should be equal")
		assert.Equal(t, uint8(55), users[1].Age, "they should be equal")
	})

	t.Run("order by name", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{
			OrderBy: interfaces.OrderByName,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, "third", users[0].Name, "they should be equal")
		assert.Equal(t, "six", users[1].Name, "they should be equal")
	})

	t.Run("age range", func(t *testing.T) {
		users, err := r.GetAll(ctx, interfaces.Filters{
			AgeGte: 22,
			AgeLte: 45,
		})
		if err != nil {
			t.Fatalf("creating user table: %v", err)
		}

		assert.Equal(t, 4, len(users), "they should be equal")
	})
}

func TestUserMongoRepoCreate(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0.5")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = mongo.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := mongo.Database("test")

	r := repositories.NewUserRepoMongo(db)

	if err := r.Create(ctx, &interfaces.User{ID: 1, Name: "asdsa", Age: 12}); err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	user, err := r.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "asdsa", user.Name, "they should be equal")
}

func TestUserMongoRepoDelete(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0.5")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = mongo.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := mongo.Database("test")

	collection := db.Collection("users")
	if _, err := collection.InsertOne(ctx, interfaces.User{ID: 1, Name: "asdsa", Age: 12}); err != nil {
		panic(err)
	}

	r := repositories.NewUserRepoMongo(db)

	if err := r.Delete(ctx, 1); err != nil {
		t.Fatalf("deleting user table: %v", err)
	}

	userDelete, err := r.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Empty(t, userDelete, "they should be equal")
}

func TestUserMongoRepoUpdate(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0.5")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err = mongo.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := mongo.Database("test")
	collection := db.Collection("users")

	u := interfaces.User{ID: 1, Name: "asdsa", Age: 12}
	if _, err := collection.InsertOne(ctx, u); err != nil {
		panic(err)
	}

	r := repositories.NewUserRepoMongo(db)

	val := map[string]interface{}{
		"name": "new name",
	}

	if err := r.Update(ctx, &u, val); err != nil {
		t.Fatalf("deleting user table: %v", err)
	}

	user, err := r.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("creating user table: %v", err)
	}

	assert.Equal(t, "new name", user.Name, "they should be equal")
}
