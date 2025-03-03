package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/posts/internal/infra/requests"
)

// Post represents the core aggregate for posts.
type Post struct {
	ID        uuid.UUID  `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    uuid.UUID  `json:"user_id"`
	Tags      []string   `json:"tags"` //todo replace Tags by list of Topics
	CreatedAt time.Time  `json:"created_at"`
	Version   int32      `json:"version" validate:"required"`
	Comments  []Comment  `json:"comments"`
	Status    PostStatus `json:"status" validate:"required" default:"pending"`
}

// PostStatus represents the valid statuses for a Post.
type PostStatus string

const (
	Pending  PostStatus = "pending"
	Approved PostStatus = "approved"
	Rejected PostStatus = "rejected"
)

// NewPostStatus validates and converts a string to a PostStatus.
func NewPostStatus(status string) (PostStatus, error) {
	validStatuses := map[string]PostStatus{
		string(Pending):  Pending,
		string(Approved): Approved,
		string(Rejected): Rejected,
	}
	if validStatus, exists := validStatuses[status]; exists {
		return validStatus, nil
	}
	return "", errors.New("invalid post status")
}

// Review allows a Post aggregate to transition its status by applying review rules.
// It returns a PostReview instance if the operation is successful.
func (p *Post) Review(ctxUser *auth.CtxUser, status string) (*PostReview, error) {
	// Ensure that the post hasn't already been approved or rejected.
	if p.Status == Approved || p.Status == Rejected {
		return nil, errors.New("post is already approved or rejected")
	}
	// Reject attempting to set the status back to pending.
	if status == string(Pending) {
		return nil, errors.New("post is already pending")
	}
	// Convert the provided status string into a valid PostStatus.
	newStatus, err := NewPostStatus(status)
	if err != nil {
		return nil, err
	}
	// Update the aggregate's status and version.
	p.Status = newStatus
	p.Version++

	// Create a new review record for this state transition.
	return NewReviewPost(p, ctxUser)
}

// PostReview captures details of a review event on a post.
type PostReview struct {
	Id        uuid.UUID  `json:"id"`
	PostId    uuid.UUID  `json:"post_id"`
	Author    uuid.UUID  `json:"author"`
	Reviewer  uuid.UUID  `json:"reviewer"`
	CreatedAt time.Time  `json:"created_at"`
	Status    PostStatus `json:"status"`
}

// NewReviewPost constructs a new PostReview based on the current state of the post.
func NewReviewPost(post *Post, reviewer *auth.CtxUser) (*PostReview, error) {
	return &PostReview{
		Id:        uuid.New(),
		PostId:    post.ID,
		Author:    post.UserID,
		Reviewer:  reviewer.ID,
		CreatedAt: time.Now(),
		Status:    post.Status,
	}, nil
}

type PostStatusChangedEvent struct {
	PostID    uuid.UUID  `json:"post_id"`
	NewStatus PostStatus `json:"new_status"`
	Author    uuid.UUID  `json:"author"`
}

// NewPostStatusChangedEvent creates an event that signals a change in the post's status.
// This event can be used for event publishing in the application layer.
func NewPostStatusChangedEvent(postID uuid.UUID, status PostStatus, author uuid.UUID) *PostStatusChangedEvent {
	return &PostStatusChangedEvent{
		PostID:    postID,
		NewStatus: status,
		Author:    author,
	}
}

// UpdatePayload represents the input for updating a post.
// swagger:model UpdatePayload
type UpdatePayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// UpdateResponse represents the output after updating a post.
// swagger:model UpdateResponse
type UpdateResponse struct {
	ID      uuid.UUID `json:"id"`
	Version int32     `json:"version"`
}

// NewPost creates a new Post aggregate from the given payload and user.
func NewPost(payload requests.PostPayload, user *auth.CtxUser) (*Post, error) {
	if len(payload.Title) == 0 || len(payload.Content) == 0 {
		return nil, errors.New("both title and content are required")
	}
	if len(payload.Title) > 100 {
		return nil, errors.New("title exceeds 100 character limit")
	}
	if user.ContextRole == auth.ADMIN || user.ContextRole == auth.Moderator &&
		(payload.Status == string(Pending) || payload.Status == string(Rejected)) {
		return nil, errors.New("admin or moderator posts don't need moderation")
	}

	return &Post{
		ID:        uuid.New(),
		Title:     payload.Title,
		Content:   payload.Content,
		UserID:    user.ID,
		Tags:      payload.Tags,
		CreatedAt: time.Now(),
		Version:   1,
		Status:    Pending,
	}, nil
}
