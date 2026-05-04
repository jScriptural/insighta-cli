package models


import (
	"time"
)
type Profile struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability"`
	Age                int       `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryName        string    `json:"country_name,omitempty"`
	CountryProbability float64   `json:"country_probability"`
	CreatedAt          time.Time `json:"created_at"`
}

type Name struct {
	Name string `json"name"`
}

type PageLinks struct {
	Self string `json:"self"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type UserID struct {
	UserID string `json:"user_id"`
}

type Response struct {
	Status     string     `json:"status"`
	Page       int        `json:"page,omitempy"`
	Limit      int        `json:"limit,omitempty"`
	Total      int        `json:"total,omitempty"`
	TotalPages int        `json:"total_pages,omitempty"`
	Links      PageLinks  `json:"links,omitempty"`
	Data       []*Profile `json:"data,omitempty"`
}

type ErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
