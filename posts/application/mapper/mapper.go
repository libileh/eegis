package mapper

import (
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

func DbToPost(dbPost database.Post) *domain.Post {
	domainPost := &domain.Post{
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

func DbToPosts(dbPosts []database.Post) []domain.Post {
	var posts []domain.Post
	for _, dbPost := range dbPosts {
		post := DbToPost(dbPost)
		posts = append(posts, *post)
	}
	return posts
}
