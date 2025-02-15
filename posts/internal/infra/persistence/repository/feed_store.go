package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

type FeedStore struct {
	Queries *database.Queries
}

func (f *FeedStore) GetFeed(ctx context.Context, userId uuid.UUID, paginate domain.Paginated) ([]domain.Feed, *errors.CustomError) {
	res, err := f.Queries.GetUserFeed(ctx, database.GetUserFeedParams{
		UserID:  userId,
		Column2: domain.MapToNullString(&paginate.Search),
		Column3: domain.MapToNullString(&paginate.Search),
		Tags:    paginate.Tags,
		Column5: domain.SortPaginateToBool(paginate.Sort),
		Limit:   paginate.Limit,
		Offset:  paginate.Offset,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return domain.DbToUserFeed(res), nil
}
