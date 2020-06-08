package model

import "github.com/pinmonl/pinmonl/model/field"

type Share struct {
	ID          string      `json:"id"`
	UserID      string      `json:"userId"`
	Slug        string      `json:"slug"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ImageID     string      `json:"imageId"`
	Status      ShareStatus `json:"status"`
	CreatedAt   field.Time  `json:"createdAt"`
	UpdatedAt   field.Time  `json:"updatedAt"`
}

func (s Share) MorphKey() string  { return s.ID }
func (s Share) MorphName() string { return "share" }

type ShareStatus int

const (
	PublishedShare ShareStatus = iota
	PreparingShare
	DeleteShare
)

type ShareList []*Share

func (sl ShareList) Keys() []string {
	var keys []string
	for _, s := range sl {
		keys = append(keys, s.ID)
	}
	return keys
}
