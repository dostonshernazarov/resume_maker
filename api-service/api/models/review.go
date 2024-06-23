package models

import "time"

type Review struct {
	ReviewId        string    `json:"review_id"`
	EstablishmentId string    `json:"establishment_id"`
	UserId          string    `json:"user_id"`
	Rating          float64   `json:"rating"`
	Comment         string    `json:"comment"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

type CreateReview struct {
	Rating  float64 `json:"rating" default:"4.7"`
	Comment string  `json:"comment" default:"very good!"`
}

type ReviewModel struct {
	ReviewId        string  `json:"review_id"`
	EstablishmentId string  `json:"establishment_id"`
	UserId          string  `json:"user_id"`
	Rating          float64 `json:"rating"`
	Comment         string  `json:"comment"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type ListReviews struct {
	Reviews []*ReviewModel `json:"reviews"`
	Count   uint64         `json:"count"`
}

