package domain

import "go.mongodb.org/mongo-driver/bson/primitive"


type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Username string `bson:"username"`
	Password string `json:"-"`
	
}