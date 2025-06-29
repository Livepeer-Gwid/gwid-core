// Package types defines types
package types

import "github.com/google/uuid"

type AuthRes struct {
	ID          uuid.UUID `json:"id"`
	Role        string    `json:"role"`
	AccessToken string    `json:"access_token"`
}

type SignupReq struct {
	Name     string `json:"name" binding:"required,min=2,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
