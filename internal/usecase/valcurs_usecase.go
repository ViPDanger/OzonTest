package usecase

import (
	"context"
	"errors"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/domain/repository"
)

type ValCursUseCase interface {
	GetByDateAndName(ctx context.Context, id string, time string, name string) (*entity.ValuteCurs, error)
	DeleteByDateAndName(ctx context.Context, id string, date string, name string) error
	Insert(ctx context.Context, item *entity.ValuteCurs) (id string, err error)
	Reset(ctx context.Context) error
}

func NewValCursUseCase(repository repository.ValCursRepository) ValCursUseCase {
	return &valCursUseCase{repository: repository}
}

type valCursUseCase struct {
	repository repository.ValCursRepository
}

func (uc *valCursUseCase) GetByDateAndName(ctx context.Context, id string, date string, name string) (*entity.ValuteCurs, error) {
	if uc.repository == nil {
		return nil, errors.New("ValCursUseCase.GetByDate(): Nil pointer repository")
	}
	item, err := uc.repository.GetByDateAndName(ctx, id, date, name)

	// СПОРНЫЙ МОМЕНТ. т.к. в доп условиях написано про уникальность данных/ответов, предполгается что все
	// данные будут загружаться по gRPC перед проверкой. поэтому удаляем обьект из бд по нахождению
	if item != nil {
		err = uc.repository.DeleteByDateAndName(ctx, id, date, name)
	}
	return item, err
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

func (uc *valCursUseCase) DeleteByDateAndName(ctx context.Context, id string, date string, name string) error {
	if uc.repository == nil {
		return errors.New("ValCursUseCase.DeleteByDateAndName(): Nil pointer in repository")
	}
	return uc.repository.DeleteByDateAndName(ctx, id, date, name)
}
