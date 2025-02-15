package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

type CommentPayload struct {
	Content string    `json:"content" validate:"required"`
	PostId  uuid.UUID `json:"post_id" validate:"required"`
}

// swagger:model Comment
type Comment struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func DbToUserComment(dbComment database.GetCommentsByPostIDRow) *Comment {
	return &Comment{
		ID:        dbComment.ID,
		Content:   dbComment.Content,
		CreatedAt: dbComment.CreatedAt,
		UserID:    dbComment.Commentuserid,
		PostID:    dbComment.PostID,
		Username:  dbComment.Username,
	}
}
func DbToComments(dbComments []database.GetCommentsByPostIDRow) []Comment {

	comments := []Comment{}
	for _, dbC := range dbComments {
		comments = append(comments, *DbToUserComment(dbC))
	}
	return comments
}

// context: comment-management
/*
### **Comment Management Context**
#### Responsibility:
Manages the comments associated with a post (Create, Read, Update, Delete). Even though `comments` appear as a property in the `Post` entity, treating it as a separate context avoids overloading `Post`'s responsibilities.
#### Features:
- CRUD operations for comments, where each comment references a post.
- Ensure comments are not edited or deleted once added (if required).*/

func AddComment(post *Post, comment Comment) error {
	if post == nil {
		return errors.New("post cannot be nil")
	}

	post.Comments = append(post.Comments, comment)
	return nil
}
