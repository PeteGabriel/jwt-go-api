package providers

import (
	"context"
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/domain"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//IUserProvider is the contract other layers can use to rely
//for logic related to user data.
type IUserProvider interface {
	CreateAccount(user *domain.User) error
	UsernameExists(username string) (bool, error)
	FindByUsername(username string) (*domain.User, error)
}

type UserProvider struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserProvider(cfg *config.Settings, client *mongo.Client) IUserProvider {

	col := client.Database(cfg.DbName).Collection("users")

	return &UserProvider{
		userCollection: col,
		ctx:            context.TODO(),
	}
}

func (u UserProvider) CreateAccount(user *domain.User) error {
	
	_, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return errors.Wrap(err, "Error inserting user")
	}

	return nil
}

func (u UserProvider) UsernameExists(username string) (bool, error) {

	usr, err := u.FindByUsername(username)
	if err != nil {
		return false, errors.Wrap(err, "Error checking if username exists")
	}
	
	return (usr != nil && usr.Username == username), nil
}

func (u UserProvider) FindByUsername(username string) (*domain.User, error) {

	var userFound domain.User
	query := bson.D{primitive.E{Key: "username", Value: username}}
	if err := u.userCollection.FindOne(u.ctx, query).Decode(&userFound); err != nil {
		//if theres no documents whatsoever
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error finding by username")
	}

	return &userFound, nil
}
