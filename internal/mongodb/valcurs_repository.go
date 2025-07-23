package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ViPDanger/OzonTest/internal/domain/entity"
	"github.com/ViPDanger/OzonTest/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewValCursRepository(db *mongo.Database) repository.ValCursRepository {
	r := valCursRepository{}
	if db != nil {
		r.collection = db.Collection("ValCurs")
	}
	return &r
}

type valCursRepository struct {
	collection *mongo.Collection
}

func (r *valCursRepository) GetByDateAndName(ctx context.Context, id string, date string, name string) (*entity.ValuteCurs, error) {
	if r.collection == nil {
		return nil, errors.New("valCursRepository.GetByDate(): nil pointer collection")
	}
	filter := bson.M{
		"creatorid": id,
		"date":      date,
		"name":      name,
	}
	var result entity.ValuteCurs
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Ничего не найдено
		}
		return nil, fmt.Errorf("valCursRepository.GetByDate()/%w", err)
	}

	return &result, nil
}
func (r *valCursRepository) DeleteByDateAndName(ctx context.Context, id string, date string, name string) error {
	if r.collection == nil {
		return errors.New("valCursRepository.DeleteByDate(): nil pointer collection")
	}
	filter := bson.M{
		"creatorid": id,
		"date":      date,
		"name":      name,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil // Ничего не найдено
		}
		return fmt.Errorf("valCursRepository.DeleteByDate()/%w", err)
	}

	return nil
}
func (r *valCursRepository) Insert(ctx context.Context, item *entity.ValuteCurs) (id string, err error) {
	if r.collection == nil {
		return "", errors.New("valCursRepository.Insert(): nil pointer collection")
	}
	res, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("valCursRepository.Insert(): failed to convert inserted ID to ObjectID")
	}

	return oid.Hex(), nil
}
func (r *valCursRepository) Reset(ctx context.Context) (err error) {
	if r.collection == nil {
		return errors.New("valCursRepository.Reset(): nil pointer collection")
	}
	_, err = r.collection.DeleteMany(ctx, bson.M{})
	return
}
