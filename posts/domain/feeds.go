package domain

import (
	"database/sql"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
	"github.com/libileh/eegis/posts/pkg/client"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Feed struct {
	Post
	UserDT0       client.UserDTO `json:"user"`
	TotalComments int64          `json:"total_comments"`
}

// swagger:model Paginated
type Paginated struct {
	Limit  int32    `json:"limit" validate:"required,gte=1,lte=20"`
	Offset int32    `json:"offset" validate:"required,gte=0"`
	Sort   string   `json:"sort" validate:"omitempty,oneof=asc desc"`
	Tags   []string `json:"tags" validate:"omitempty,max=5"`
	Search string   `json:"search" validate:"omitempty,min=1,max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func SortPaginateToBool(sort string) bool {
	if sort == "asc" {
		return false
	}
	return true
}

func DbToUserFeed(dbUserFeeds []database.GetUserFeedRow) []Feed {
	userFeeds := make([]Feed, 0, len(dbUserFeeds)) // Initialize with length 0 and capacity len(dbUserFeeds)
	for _, dbUserFeed := range dbUserFeeds {
		post := Post{
			ID:        dbUserFeed.ID,
			Title:     dbUserFeed.Title,
			Content:   dbUserFeed.Content,
			Tags:      dbUserFeed.Tags,
			UserID:    dbUserFeed.UserID,
			CreatedAt: dbUserFeed.CreatedAt,
			Version:   dbUserFeed.Version,
		}
		Feed := Feed{
			Post:          post,
			UserDT0:       client.UserDTO{ID: post.UserID, Username: dbUserFeed.Username},
			TotalComments: dbUserFeed.TotalComments,
		}
		userFeeds = append(userFeeds, Feed)
	}
	return userFeeds
}

func (fp Paginated) GetFeedPaginationFilter(r *http.Request) Paginated {
	qp := r.URL.Query()

	limitParam := qp.Get("limit")
	if limitParam != "" {
		limit, err := strconv.ParseInt(limitParam, 10, 64)
		if err != nil {
			return fp
		}
		fp.Limit = int32(limit)
	}

	offsetParam := qp.Get("offset")
	if offsetParam != "" {
		offset, err := strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			return fp
		}
		fp.Offset = int32(offset)
	}
	sort := qp.Get("sort")
	if sort != "" {
		fp.Sort = strings.ToLower(sort)
	}
	tags := qp.Get("tags")
	if tags != "" {
		fp.Tags = strings.Split(tags, ",")
	}
	since := qp.Get("since")
	if since != "" {
		fp.Since = parseTime(since)
	}
	until := qp.Get("until")
	if until != "" {
		fp.Until = parseTime(until)
	}
	search := qp.Get("search")
	if search != "" {
		fp.Search = search
	}

	return fp
}

func parseTime(stringTime string) string {
	dateTime, err := time.Parse(time.DateTime, stringTime)
	if err != nil {
		return ""
	}
	return dateTime.Format(time.DateTime)
}

func MapToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
