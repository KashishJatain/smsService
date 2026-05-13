package repository

import (
	"context"
	"fmt"
	"time"
	"github.com/sms/sms-store/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const(
	collectionName= "sms_records"
	queryTimeout= 10 *time.Second
)

type SmsRepository interface{
	Save(ctx context.Context, record *model.SmsRecord) error
	FindByUserID(ctx context.Context,userID string) ([]model.SmsRecord,error)
}
type MongoSmsRepository struct{
	collection *mongo.Collection
}

func NewMongoSmsRepository(db *mongo.Database) (*MongoSmsRepository, error){
	col := db.Collection(collectionName)
	repo := &MongoSmsRepository{collection: col}

	if err := repo.ensureIndexes(); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w",err)
	}
	return repo, nil
}
func (r *MongoSmsRepository) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(),queryTimeout)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:  bson.D{{Key: "user_id", Value: 1},{Key:"created_at",Value: -1}},
			Options: options.Index().SetName("idx_user_id_created_at"),
		},
		{
			Keys: bson.D{{Key: "message_id",Value: 1}},
			Options: options.Index().SetName("idx_message_id").SetUnique(true),
		},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
func (r *MongoSmsRepository) Save(ctx context.Context,record *model.SmsRecord) error{
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	record.ID = primitive.NewObjectID()
	record.CreatedAt = time.Now().UTC()
	_, err := r.collection.InsertOne(ctx,record)
	if err != nil {
		return fmt.Errorf("failed to insert SMS record: %w", err)
	}
	return nil
}
func (r *MongoSmsRepository) FindByUserID(ctx context.Context, userID string) ([]model.SmsRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()
	filter := bson.D{{Key:"user_id", Value: userID}}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query SMS records for userID=%s: %w",userID,err)
	}
	defer cursor.Close(ctx)
	var records []model.SmsRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode SMS records: %w",err)
	}
	return records, nil
}