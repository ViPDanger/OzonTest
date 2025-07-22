package usecase

import (
	"context"
	"errors"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/domain/repository"
)

type ValCursUseCase interface {
	GetByDateAndName(ctx context.Context, time string, name string) (*entity.ValuteCurs, error)
	DeleteByDateAndName(ctx context.Context, date string, name string) error
	Insert(ctx context.Context, item *entity.ValuteCurs) (id string, err error)
	Reset(ctx context.Context) error
}

func NewValCursUseCase(repository repository.ValCursRepository) ValCursUseCase {
	return &valCursUseCase{repository: repository}
}

type valCursUseCase struct {
	repository repository.ValCursRepository
}

func (uc *valCursUseCase) GetByDateAndName(ctx context.Context, date string, name string) (*entity.ValuteCurs, error) {
	if uc.repository == nil {
		return nil, errors.New("ValCursUseCase.GetByDate(): Nil pointer repository")
	}
	return uc.repository.GetByDateAndName(ctx, date, name)
}

func (uc *valCursUseCase) Insert(ctx context.Context, item *entity.ValuteCurs) (id string, err error) {
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

func (uc *valCursUseCase) DeleteByDateAndName(ctx context.Context, date string, name string) error {
	if uc.repository == nil {
		return errors.New("ValCursUseCase.DeleteByDateAndName(): Nil pointer in repository")
	}
	return uc.repository.DeleteByDateAndName(ctx, date, name)
}
