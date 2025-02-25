package application

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/properties"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/client"
	"github.com/libileh/eegis/users/internal/infra/persistence/cache"
	"time"
)

type ServiceManager struct {
	UserRepoService    *UserRepoService
	UserCacheService   *UserCacheService
	PasswordService    *PasswordService
	NotificationClient *client.HttpNotificationClient
	TopicClient        *client.TopicHTTPClient
}

func NewServiceManager(
	userRepoService *UserRepoService,
	userCacheService *UserCacheService,
	passwordService *PasswordService,
	notificationClient *client.HttpNotificationClient,
	topicClient *client.TopicHTTPClient,
) *ServiceManager {
	return &ServiceManager{
		UserRepoService:    userRepoService,
		UserCacheService:   userCacheService,
		PasswordService:    passwordService,
		NotificationClient: notificationClient,
		TopicClient:        topicClient,
	}
}

type UserRepoService struct {
	UserRepo    domain.UserRepository
	RoleRepo    domain.RoleRepository
	FollowerRpo domain.FollowerRepository
}

type UserCacheService struct {
	UserCache cache.CacheUserRepository
}

func NewUserCacheService(userCache cache.CacheUserRepository) *UserCacheService {
	return &UserCacheService{
		UserCache: userCache,
	}
}

func NewUserService(userRepo domain.UserRepository, roleRepo domain.RoleRepository, followerRepo domain.FollowerRepository) *UserRepoService {
	return &UserRepoService{
		UserRepo:    userRepo,
		RoleRepo:    roleRepo,
		FollowerRpo: followerRepo,
	}
}

type NotificationService interface {
	SendConfirmationEmail(email, token string)
}

type TopicClient interface {
	FollowTopic(topicId uuid.UUID, userId uuid.UUID) *errors.CustomError
}

// todo add to common auth package

// MapToJWTClaims maps a user to JWT claims.
// It takes a user and an authentication properties object as input and returns a JWT.MapClaims object.
// The JWT.MapClaims object contains the user's ID, role, expiration time, issued at time, not before time, issuer, and audience.
// The function also maps the user's role to a context role using the MapToCtx
func (s *ServiceManager) MapToJWTClaims(userId uuid.UUID, roleID int16, authProps properties.AuthProperties) (*jwt.MapClaims, error) {
	//generate token: add claims
	ctxRole, err := auth.MapToCtxRole(float64(roleID))
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"sub":  userId,
		"role": ctxRole,
		"exp":  time.Now().Add(authProps.Token.Exp).Unix(),
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"iss":  authProps.Issuer,
		"aud":  authProps.Audience,
	}
	return &claims, err
}

func (s *ServiceManager) FollowTopic(user *auth.CtxUser, topicId uuid.UUID) *errors.CustomError {

	claims, err := s.MapToJWTClaims(user.ID, user.ContextRole.Value, s.TopicClient.AuthProps)
	if err != nil {
		return errors.NewInternalServerError("failed to map claims: %v", err)
	}
	if err := s.TopicClient.FollowTopic(topicId, user.ID, *claims); err != nil {
		return err
	}
	return nil
}
