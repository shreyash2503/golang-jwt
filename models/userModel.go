package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID			         primitive.ObjectID  `bson:"_id`
	first_name    		 *string  		     `json:"first_name" validate:"required, min=2, max=100"`
	last_name  			 *string             `json:"last_name" validate:"required, min=2, max=100"`
	password   			 *string             `json:"password" validate:"required, min=6"`
	email                *string             `json:"email" validate:"required, email"`
	phone                *string             `json:"phone" validate:"required, min=10, max=10"` 
	token 			  	 *string             `json:"token"`
	user_type 		     *string             `json:"user_type" validate: required, eq=ADMIN|eq=USER`
	refresh_token        *string             `json:"refresh_token"`
	created_at           time.Time           `json:"created_at"`
	updated_at           time.Time           `json:"updated_at"`
	User_id              *string             `json:"User_id"`

}