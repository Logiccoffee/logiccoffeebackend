package content

import (
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	ContentItems []ContentItem `json:"data"`
	Meta         Meta          `json:"meta"`
}

type ContentItem struct {
	ContentID          string     `json:"content_id"`
	Type               string     `json:"type"`
	ContentType        string     `json:"content_type"`
	ContentName        string     `json:"content_name"`
	ContentDescription *string    `json:"content_description"`
	Duration           int        `json:"duration"`
	FileSize           string     `json:"file_size"`
	FilePath           *string    `json:"file_path"`
	Status             string     `json:"status"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	ApprovedBy         *string    `json:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at"`
}

type Meta struct {
	CurrentPage string `json:"current_page"`
	FirstItem   int    `json:"first_item"`
	LastItem    int    `json:"last_item"`
	LastPage    int    `json:"last_page"`
	Total       int    `json:"total"`
}
