package repository

import (
	"context"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
)

type XMLDailyMockResponseRepository interface {
	GetByDate(ctx context.Context, date string) (*entity.XMLDailyMockResponse, error)
	DeleteByDate(ctx context.Context, date string) error
	Insert(ctx context.Context, item *entity.XMLDailyMockResponse) (id string, err error)
	Reset(ctx context.Context) (err error)
}
