package model

import "gorm.io/gorm"

type Metadata struct {
	gorm.Model
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}
