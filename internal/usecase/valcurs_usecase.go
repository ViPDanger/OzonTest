package usecase

import (
	"context"
	"errors"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/domain/repository"
)

type ValCursUseCase interface {
	GetByDate(ctx context.Context, time string) (*entity.XMLDailyMockResponse, error)
	DeleteByDate(ctx context.Context, date string) error
	Insert(ctx context.Context, item *entity.XMLDailyMockResponse) (id string, err error)
	Reset(ctx context.Context) error
}

func NewValCursUseCase(repository repository.XMLDailyMockResponseRepository) ValCursUseCase {
	return &valCursUseCase{repository: repository, deleteAfterSearch: true}
}

type valCursUseCase struct {
	repository        repository.XMLDailyMockResponseRepository
	deleteAfterSearch bool
}

func (uc *valCursUseCase) GetByDate(ctx context.Context, date string) (*entity.XMLDailyMockResponse, error) {
	if uc.repository == nil {
		return nil, errors.New("ValCursUseCase.GetByDate(): Nil pointer repository")
	}
	item, err := uc.repository.GetByDate(ctx, date)
	if item != nil && uc.deleteAfterSearch {
		err = uc.repository.DeleteByDate(ctx, date)
	}
	return item, err
}

func (uc *valCursUseCase) Insert(ctx context.Context, item *entity.XMLDailyMockResponse) (id string, err error) {
	if uc.repository == nil || item == nil {
		return "", errors.New("ValCursUseCase.Insert(): Nil pointer")
	}
	return uc.repository.Insert(ctx, item)
}

func (uc *valCursUseCase) Reset(ctx context.Context) error {
	if uc.repository == nil {
		return errors.New("ValCursUseCase.Reset(): Nil pointer in repository")
	}
	return uc.repository.Reset(ctx)
}

func (uc *valCursUseCase) DeleteByDate(ctx context.Context, date string) error {
	if uc.repository == nil {
		return errors.New("ValCursUseCase.DeleteByDateAndName(): Nil pointer in repository")
	}
	return uc.repository.DeleteByDate(ctx, date)
}
