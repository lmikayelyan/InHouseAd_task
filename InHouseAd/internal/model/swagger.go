package model

type GoodInputResponse struct {
	Name       string `json:"name"`
	Categories []uint `json:"categories"`
}

type CategoryInputResponse struct {
	Name string `json:"name"`
}
