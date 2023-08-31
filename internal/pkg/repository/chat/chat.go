package chat

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	Title    string
	Settings *Setting
	Num      int
}

type ChatDoc struct {
	Chat `bson:"inline"`
	ID   *primitive.ObjectID `bson:"_id"`
}

func (cd ChatDoc) StringID() string {
	return cd.ID.Hex()
}
