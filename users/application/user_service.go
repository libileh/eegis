package application

import (
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/client"
	"github.com/libileh/eegis/users/internal/infra/persistence/cache"
)

type ServiceManager struct {
	UserRepoService     *UserRepoService
	UserCacheService    *UserCacheService
	PasswordService     *PasswordService
	NotificationService *client.HttpNotificationService
}

func NewServiceManager(
	userRepoService *UserRepoService,
	userCacheService *UserCacheService,
	passwordService *PasswordService,
	notificationService *client.HttpNotificationService) *ServiceManager {
	return &ServiceManager{
		UserRepoService:     userRepoService,
		UserCacheService:    userCacheService,
		PasswordService:     passwordService,
		NotificationService: notificationService,
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
