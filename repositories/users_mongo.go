package repositories

import (
	"context"
	"time"

	"repos/interfaces"
	"repos/utils"

	"go.mongodb.org/mongo-driver/bson"
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

func (r userRepoMongo) GetAll(ctx context.Context, filters interfaces.Filters, opts ...utils.Options) ([]*interfaces.User, int64, error) {

	var limit int64 = utils.Limit
	if filters.Limit != 0 {
		limit = int64(filters.Limit)
	}

	var offset int64 = 0
	if filters.Offset != 0 {
		offset = int64(filters.Offset)
	}

	var sort = "createdat"
	if filters.OrderBy != "" {
		sort = string(filters.OrderBy)
	}

	options := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
		Sort:  bson.D{{sort, -1}},
	}

	f := bson.A{}

	now := time.Date(2024, 4, 17, 0, 0, 0, 0, time.Local)

	var createdAtGte time.Time = now.AddDate(0, 0, -30)
	if !filters.CreatedAtGte.IsZero() {
		createdAtGte = filters.CreatedAtGte
	}

	var createdAtLte time.Time = now
	if !filters.CreatedAtLte.IsZero() {
		createdAtLte = filters.CreatedAtLte
	}

	f = append(f, bson.D{{"createdat", bson.D{{"$lte", createdAtLte}}}})
	f = append(f, bson.D{{"createdat", bson.D{{"$gt", createdAtGte}}}})

	if filters.AgeGte != 0 {
		f = append(f, bson.D{{"age", bson.D{{"$gte", filters.AgeGte}}}})
	}

	if filters.AgeLte != 0 {
		f = append(f, bson.D{{"age", bson.D{{"$lte", filters.AgeLte}}}})
	}

	filter := bson.D{{"$and", f}}

	cursor, err := r.collection.Find(ctx, filter, &options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	var users []*interfaces.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, 0, nil
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

func (r userRepoMongo) Delete(ctx context.Context, ids []int64, opts ...utils.Options) error {
	_, err := r.collection.DeleteMany(ctx, bson.D{{"id", ids}})

	return err
}
