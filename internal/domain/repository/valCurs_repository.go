package repository

import (
	"context"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
)

type ValCursRepository interface {
	GetByDateAndName(ctx context.Context, date string, name string) (*entity.ValuteCurs, error)
	DeleteByDateAndName(ctx context.Context, date string, name string) error
	Insert(ctx context.Context, item *entity.ValuteCurs) (id string, err error)
	Reset(ctx context.Context) (err error)
}
