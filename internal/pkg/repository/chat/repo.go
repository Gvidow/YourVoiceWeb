package chat

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

const collectionName = "chats"

var ErrUnknowTypeID = errors.New("unknown id type of the inserted record")

type chatRepo struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) *chatRepo {
	return &chatRepo{db.Collection(collectionName)}
}

func (r *chatRepo) SelectAllByOrder(ctx context.Context) ([]ChatDoc, error) {
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.M{"num": 1}))
	if err != nil {
		return nil, fmt.Errorf("call method SelectAllByOrder from *chatRepo: find: %w", err)
	}
	var res []ChatDoc
	err = cursor.All(ctx, &res)
	if err != nil {
		return res, fmt.Errorf("call method SelectAllByOrder from *chatRepo: decode all: %w", err)
	}
	err = cursor.Close(ctx)
	return res, err
}

func (r *chatRepo) DeleteMany(ctx context.Context, ids []string) (int, error) {
	objectIds, err := castStringsToObjectIds(ids)
	if err != nil {
		return 0, fmt.Errorf("call method DeleteMany from *chatRepo: %w", err)
	}
	res, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if res == nil {
		return 0, fmt.Errorf("call method DeleteMany from *chatRepo: %w", err)
	}
	return int(res.DeletedCount), err
}

func (r *chatRepo) SaveSettings(ctx context.Context, id string, settings *Setting) error {
	objID := primitive.NewObjectID()
	err := objID.UnmarshalText([]byte(id))
	if err != nil {
		return fmt.Errorf("call method SaveSettings from *chatRepo: unmarshal ObjectID: %w", err)
	}
	_, err = r.collection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"settings": settings}})
	if err != nil {
		return fmt.Errorf("call method SaveSettings from *chatRepo: update: %w", err)
	}
	return nil
}

func (r *chatRepo) SwapPlaces(ctx context.Context, id1, id2 string) error {
	objID1 := primitive.NewObjectID()
	err := objID1.UnmarshalText([]byte(id1))
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: unmarshal ObjectID from %s: %w", id1, err)
	}
	objID2 := primitive.NewObjectID()
	err = objID2.UnmarshalText([]byte(id2))
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: unmarshal ObjectID from %s: %w", id2, err)
	}

	res := r.collection.FindOne(ctx, bson.M{"_id": objID1}, options.FindOne().SetProjection(bson.M{"num": 1}))
	var num1, num2 struct{ Num int }
	err = res.Decode(&num1)
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: decode document in num: %w", err)
	}

	res = r.collection.FindOne(ctx, bson.M{"_id": objID2}, options.FindOne().SetProjection(bson.M{"num": 1}))
	err = res.Decode(&num2)
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: decode document in num: %w", err)
	}

	_, err = r.collection.UpdateByID(ctx, objID1, bson.M{"$set": bson.M{"num": num2.Num}})
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: save num: %w", err)
	}

	_, err = r.collection.UpdateByID(ctx, objID2, bson.M{"$set": bson.M{"num": num1.Num}})
	if err != nil {
		return fmt.Errorf("call SwapPlaces from *chatRepo: save num: %w", err)
	}
	return nil
}

func (r *chatRepo) AddNewChat(ctx context.Context, title string) (string, error) {
	curr, err := r.collection.Aggregate(ctx, bson.A{bson.M{"$group": bson.M{"_id": nil, "max": bson.M{"$max": "$num"}}}})
	if err != nil {
		return "", fmt.Errorf("call AddNewChat from *chatRepo: couldn't determine the num for the new chat: %w", err)
	}
	var m struct{ Max int }
	if curr.Next(ctx) {
		curr.Decode(&m)
	}

	chatNew := &Chat{
		Title:    title,
		Num:      m.Max + 1,
		Settings: &Setting{},
	}

	res, err := r.collection.InsertOne(ctx, chatNew)
	if err != nil {
		return "", fmt.Errorf("call AddNewChat from *chatRepo: insertion error: %w", err)
	}
	if objID, ok := res.InsertedID.(primitive.ObjectID); ok {
		return objID.Hex(), nil
	}
	return "", fmt.Errorf("call AddNewChat from *chatRepo: %w", ErrUnknowTypeID)
}

func (r *chatRepo) EditChat(ctx context.Context, id string, newTitle string) error {
	objID := primitive.NewObjectID()
	err := objID.UnmarshalText([]byte(id))
	if err != nil {
		return fmt.Errorf("call EditChat from *chatRepo: invalid id: %w", err)
	}
	_, err = r.collection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"title": newTitle}})
	if err != nil {
		return fmt.Errorf("call EditChat from *chatRepo: error when changing the title: %w", err)
	}
	return nil
}

func castStringsToObjectIds(ids []string) ([]primitive.ObjectID, error) {
	res := make([]primitive.ObjectID, len(ids))
	var err error
	for i, strID := range ids {
		err = res[i].UnmarshalText([]byte(strID))
		if err != nil {
			return nil, fmt.Errorf("cast string to ObjectID: %w", err)
		}
	}
	return res, err
}
