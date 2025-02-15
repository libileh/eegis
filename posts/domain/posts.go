package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

// swagger:model PostPayload
type PostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
}

// swagger:model UpdatePayload
type UpdatePayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// swagger:model UpdateResponse
type UpdateResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int32     `json:"version"`
}

// swagger:model Post
type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    uuid.UUID `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version" validate:"required"`
	Comments  []Comment `json:"comments"`
}

// context: post-creation
func NewPost(payload PostPayload, userID uuid.UUID) (*Post, error) {
	if len(payload.Title) == 0 || len(payload.Content) == 0 {
		return nil, errors.New("both title and content are required")
	}

	if len(payload.Title) > 100 {
		return nil, errors.New("title exceeds 100 character limit")
	}

	return &Post{
		ID:        uuid.New(),
		Title:     payload.Title,
		Content:   payload.Content,
		UserID:    userID,
		Tags:      payload.Tags,
		CreatedAt: time.Now(),
		Version:   1,
	}, nil
}

// context: post-updates
func UpdatePost(request *UpdatePayload, post *Post) (*Post, error) {
	if request.Title == "" && request.Content == "" && request.Tags == nil {
		return nil, errors.New("update payload is empty")
	}

	// Apply changes only where new values are provided
	if request.Title != "" {
		post.Title = request.Title
	}
	if request.Content != "" {
		post.Content = request.Content
	}
	if request.Tags != nil {
		post.Tags = request.Tags
	}

	// Increment version to protect against stale writes
	post.Version += 1
	return post, nil
}

func DbToPost(dbPost database.Post) *Post {
	domainPost := &Post{
		ID:        dbPost.ID,
		Title:     dbPost.Title,
		Content:   dbPost.Content,
		Tags:      dbPost.Tags,
		CreatedAt: dbPost.CreatedAt,
		Version:   dbPost.Version,
		UserID:    dbPost.UserID,
	}

	return domainPost
}

func DbToPosts(dbPosts []database.Post) []Post {
	var posts []Post
	for _, dbPost := range dbPosts {
		post := DbToPost(dbPost)
		posts = append(posts, *post)
	}
	return posts
}
