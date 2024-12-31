package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID             primitive.ObjectID      `bson:"_id"`
	User           primitive.ObjectID      `bson:"user"`  
	Valid          bool                    `bson:"valid"`
	Refresh_token  string                  `bson:"refresh_token"`
	User_agent     string                  `bson:"user_agent"` 
	Created_at     time.Time               `bson:created_at` 
	Last_used      time.Time               `bson:last_used`
}