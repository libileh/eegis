package seed

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/repository"
	"log"
	"math/rand"
	"time"
)

const (
	numUsers  = 52
	batchSize = 1000
)

func Seed(store repository.Storage) {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	log.Println("Starting database seeding...")

	// Create separate slices for user IDs
	seedUsers(ctx, store.Users(), numUsers)

	log.Println("Database seeding completed.")
}

func seedUsers(ctx context.Context, userRepo repository.UserRepository, count int) []uuid.UUID {
	var password domain.Password
	err := password.Set("password")
	if err != nil {
		log.Fatalf("Error setting password: %v", err)
	}

	log.Println("Seeding users...")
	userIDs := make([]uuid.UUID, count)
	users := make([]*domain.User, 0, batchSize)

	for i := 0; i < count; i++ {
		userID := uuid.New()
		userIDs[i] = userID

		user := &domain.User{
			ID:       userID,
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@qoraal.com", i),
			Password: password,
			//default role is set
			RoleID: 1,
		}
		users = append(users, user)

		if len(users) == batchSize {
			bulkInsertUsers(ctx, userRepo, users)
			users = users[:0]
		}
	}
	if len(users) > 0 {
		bulkInsertUsers(ctx, userRepo, users)
	}

	return userIDs
}

func bulkInsertUsers(ctx context.Context, store repository.UserRepository, users []*domain.User) {
	for _, u := range users {
		if _, err := store.CreateUser(ctx, u); err != nil {
			log.Printf("Error inserting user batch: %v\n", err)
		}
	}
}
