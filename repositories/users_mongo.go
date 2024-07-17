package repositories

import (
	"context"

	"repos/interfaces"
	"repos/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepoMongo struct {
	collection *mongo.Collection
}

func NewUserRepoMongo(collection *mongo.Database) interfaces.UsersRepo {
	return &userRepoMongo{collection: collection.Collection("users")}
}

func (r userRepoMongo) GetById(ctx context.Context, id int64, opts ...utils.Options) (*interfaces.User, error) {
	var user *interfaces.User
	if err := r.collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r userRepoMongo) GetAll(ctx context.Context, filters interfaces.Filters, opts ...utils.Options) ([]*interfaces.User, error) {

	var limit int64 = utils.Limit
	if filters.Limit != 0 {
		limit = int64(filters.Limit)
	}

	// var offset int64 = 1
	// if filters.Offset != 0 {
	// 	offset = int64(filters.Offset)
	// }

	options := options.FindOptions{
		// Skip:  &offset,
		Limit: &limit,
		Sort:  bson.D{{"createdat", -1}},
	}

	var ageGte primitive.D = bson.D{}
	if filters.AgeGte != 0 {
		ageGte = bson.D{{"age", bson.D{{"$gte", filters.AgeGte}}}}
	}

	var ageLte primitive.D = bson.D{}
	if filters.AgeLte != 0 {
		ageLte = bson.D{{"age", bson.D{{"$lte", filters.AgeLte}}}}
	}

	filter := bson.D{
		{"$and",
			bson.A{
				ageGte,
				ageLte,
			}},
	}

	cursor, err := r.collection.Find(ctx, filter, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	var users []*interfaces.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r userRepoMongo) Create(ctx context.Context, user *interfaces.User, opts ...utils.Options) error {
	_, err := r.collection.InsertOne(ctx, user)

	return err
}

func (r userRepoMongo) Update(ctx context.Context, user *interfaces.User, vals map[string]interface{}, opts ...utils.Options) error {
	updates := bson.D{}
	for k, v := range vals {
		update := bson.E{Key: k, Value: v}
		updates = append(updates, update)
	}

	updateFilter := bson.D{{"$set", updates}}

	if _, err := r.collection.UpdateOne(ctx, bson.D{{"id", user.ID}}, updateFilter); err != nil {
		return err
	}

	return nil
}

func (r userRepoMongo) Delete(ctx context.Context, id int64, opts ...utils.Options) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{"id", id}})

	return err
}
