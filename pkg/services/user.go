package services

import (
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/data/providers"
	"jwtGoApi/pkg/domain"
	"jwtGoApi/pkg/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//IUserService defines the contract that help controllers
//to execute business logic. It mixes the errors package and
//the data package.
type IUserService interface {
	CreateAccount(user *domain.User) *models.Error
	Login(user *domain.User) (string, *models.Error)
}

type UserService struct {
	provider providers.IUserProvider
	cfg      *config.Settings
}

func NewUserService(cfg *config.Settings, provider providers.IUserProvider) IUserService {
	return &UserService{
		provider: provider,
		cfg:      cfg,
	}
}

func (svc UserService) CreateAccount(user *domain.User) *models.Error {
	//check if the account/user already exists
	if err := svc.validateUsername(user.Username); err != nil {
		return err
	}

	user.ID = primitive.NewObjectID()
	hash, err := hashPassword(user.Password)
	if err != nil {
		return &models.Error{
			Message: "error happened while trying to create an hash of user's password.",
			Code:    500,
			Name:    "SERVER_ERROR",
			Error:   err,
		}
	}

	user.Password = hash
	err = svc.provider.CreateAccount(user)
	if err != nil {
		return &models.Error{
			Message: "error happened while trying to create a new account.",
			Code:    500,
			Name:    "SERVER_ERROR",
			Error:   err,
		}
	}

	return nil

}

func (svc UserService) Login(user *domain.User) (string, *models.Error) {
	//username must exist and account related must have same password
	userFound, err := svc.provider.FindByUsername(user.Username)
	if err != nil {
		return "", &models.Error{
			Message: "error happened while trying to check if account/user already exists",
			Code:    500,
			Name:    "SERVER_ERROR",
			Error:   err,
		}
	}

	if userFound == nil {
		return "", &models.Error{
			Message: "your username and/or password was not found",
			Code:    400,
			Name:    "INVALID_LOGIN",
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		return "", &models.Error{
			Message: "your username and/or password was not found",
			Code:    400,
			Name:    "INVALID_LOGIN",
		}
	}

	token, err := svc.createJwtToken(userFound.ID.Hex())
	if err != nil {
		return "", &models.Error{
			Message: "error happened while trying to create a new token",
			Code:    500,
			Name:    "SERVER_ERROR",
			Error:   err,
		}
	}

	return token, nil
}

func (svc UserService) validateUsername(username string) *models.Error {
	exists, err := svc.provider.UsernameExists(username)
	if err != nil {
		return &models.Error{
			Message: "error happened while trying to check if account/user already exists",
			Code:    500,
			Name:    "SERVER_ERROR",
			Error:   err,
		}
	}

	if exists {
		return &models.Error{
			Message: "account/user already exists with the same username",
			Code:    400,
			Name:    "USERNAME_TAKEN",
		}
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "error hashing password")
	}
	return string(hash), nil
}

func (svc UserService) createJwtToken(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)
	expiresIn, err := strconv.ParseInt(svc.cfg.JwtExpires, 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "error parsing int")
	}

	expiration := time.Duration(int64(time.Minute) * expiresIn)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(expiration).Unix()

	signedTkn, err := token.SignedString([]byte(svc.cfg.JwtSecret))
	if err != nil {
		return "", errors.Wrap(err, "error signing token")
	}
	
	return signedTkn, nil
}
