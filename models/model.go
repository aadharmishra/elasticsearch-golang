package models

type Author struct {
	Name        string  `json:"name,omitempty"`
	ReleaseDate string  `json:"release_date,omitempty"`
	Author      string  `json:"author,omitempty"`
	PageCount   string  `json:"page_count,omitempty"`
	Details     Details `json:"details,omitempty"`
}

type Details struct {
	Index        string `json:"index,omitempty"`
	DocumentID   string `json:"documentID,omitempty"`
	DocumentType string `json:"documentType,omitempty"`
}

type Config struct {
	EsCredentials EsCredentials `json:"esCredentials,omitempty"`
}

type EsCredentials struct {
	Id     string `json:"id,omitempty"`
	ApiKey string `json:"apiKey,omitempty"`
}
