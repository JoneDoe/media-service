package models

type Version struct {
	Original struct {
		Filename, Url string
		Size          int
	}
}

type FSFile struct {
	FileName string `json:"fileName"`
	Url      string `json:"url"`
}
