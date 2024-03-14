package model

type Good struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Categories []uint `json:"categories"`
}
