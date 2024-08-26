package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/svetlana-mel/event-task-planner/internal/lib/jwt"
	sl "github.com/svetlana-mel/event-task-planner/internal/lib/slog"
	"github.com/svetlana-mel/event-task-planner/internal/models"
	base "github.com/svetlana-mel/event-task-planner/internal/repository"
)

const TokenTTL = time.Hour * 24

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user with the specified email already exists")
	ErrUserNotExists      = errors.New("user with the specified email already exists")
	ErrWrongPassword      = errors.New("wrong password")
)

var _ Auth = (*authProvider)(nil)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)

	SignUp(
		ctx context.Context,
		username string,
		email string,
		password string,
	) (userID uint64, err error)
}

type authProvider struct {
	secret      string
	log         *slog.Logger
	userCreator UserCreator
	usrProvider UserProvider
	tokenTTL    time.Duration
}

type UserCreator interface {
	CreateUser(
		ctx context.Context,
		username string,
		email string,
		passHash []byte,
	) (uint64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
}

func NewAuth(
	secret string,
	log *slog.Logger,
	userCreator UserCreator,
	userProvider UserProvider,
	tokenTTL time.Duration,
) Auth {
	return &authProvider{
		userCreator: userCreator,
		usrProvider: userProvider,
		log:         log,
		tokenTTL:    tokenTTL,
		secret:      secret,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
// if user not exists - error
// if wrong password - error
func (ap *authProvider) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "services.auth.Login"

	log := ap.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := ap.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, base.ErrUserNotExists) {
			log.Info("user not exists", sl.AddErrorAtribute(err))
			return "", ErrUserNotExists
		}
		log.Error("error get user", sl.AddErrorAtribute(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Info("invalid password", sl.AddErrorAtribute(err))
		return "", ErrWrongPassword
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, ap.secret, TokenTTL)
	if err != nil {
		log.Error("error generate token", sl.AddErrorAtribute(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// SignUp - Register a new user
// if user already exists - error (== email exists)
func (ap *authProvider) SignUp(
	ctx context.Context,
	username string,
	email string,
	password string,
) (userID uint64, err error) {
	const op = "services.auth.Signup"

	log := ap.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := ap.userCreator.CreateUser(ctx, username, email, passHash)
	if err != nil {
		if errors.Is(err, base.ErrUserAlreadyExists) {
			log.Info("user already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("signup successfully")

	return id, nil
}
