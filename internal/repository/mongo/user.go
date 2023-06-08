package mongo

import (
	"WebProject_part7/internal/core"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (repository *UserRepository) GetAll(ctx context.Context) ([]*core.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	usersChannel := make(chan []*core.User, 0)

	var err error
	go func() {
		err = repository.retrieveAll(ctx, usersChannel)
	}()

	if err != nil {
		return nil, err
	}

	var userList []*core.User
	select {
	case <-ctxTimeout.Done():
		fmt.Println("Processing timeot in Mongo")
		break
	case userList = <-usersChannel:
		fmt.Println("Finished processing in Mongo")
	}
	return userList, nil
}

func (repository *UserRepository) retrieveAll(ctx context.Context, channel chan<- []*core.User) (err error) {
	cursor, err := repository.collection.Find(ctx, bson.D{})
	users := make([]*core.User, 0)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		return err
	}
	channel <- users
	return nil
}

func (repository *UserRepository) GetById(ctx context.Context, id string) (*core.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	userChannel := make(chan *core.User, 0)

	var err error
	go func() {
		err = repository.retrieveUser(ctx, id, userChannel)
	}()

	if err != nil {
		return nil, err
	}

	var user *core.User

	select {
	case <-ctxTimeout.Done():
		fmt.Println("Processing timeot in Mongo")
		break
	case user = <-userChannel:
		fmt.Println("Finished processing in Mongo")
	}

	return user, nil
}

func (repository *UserRepository) retrieveUser(ctx context.Context, id string, channel chan<- *core.User) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	user := &core.User{}
	err = repository.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(user)

	if err != nil {
		return err
	}

	channel <- user
	return nil
}

func (repository *UserRepository) Save(ctx context.Context, user *core.User) (*core.User, error) {
	result, err := repository.collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}
