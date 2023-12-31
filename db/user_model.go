package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Nitesh000/hotel-reservation-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userColl = "users"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper

	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	PostUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, bson.M, types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client, dbname string) *MongoUserStore {
	mongoUserStore := &MongoUserStore{
		client: c,
		coll:   c.Database(dbname).Collection(userColl),
	}

	return mongoUserStore
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	// var user *types.User
	// for cur.Next(ctx) {
	// 	if err := cur.Decode(&user); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, user)
	// }
	// defer cur.Close(ctx)
	return users, nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- droppping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) PostUser(ctx context.Context, user *types.User) (*types.User, error) {
	cur, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = cur.InsertedID.(primitive.ObjectID)

	return &types.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, userId string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if deletecount := res.DeletedCount; deletecount == 0 {
		log.Println("No user is deleted!!")
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, update types.UpdateUserParams) error {
	values := bson.M{}
	values["firstName"] = update.FirstName
	values["lastName"] = update.LastName
	log.Println(values)
	updateObject := bson.D{
		{
			"$set", values,
		},
	}

	_, err := s.coll.UpdateOne(ctx, filter, updateObject)
	if err != nil {
		return err
	}
	return nil
}
