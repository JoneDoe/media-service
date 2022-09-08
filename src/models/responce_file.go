package models

type OutputModel struct {
	FileName string `json:"fileName"`
	Uuid     string `json:"uuid"`
	Url      string `json:"url"`
}

type MetaData struct {
	OutputModel
	Key string `json:"key"`
	Url string `json:"url"`
}
