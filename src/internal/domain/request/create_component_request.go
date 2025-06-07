package request

import "github.com/google/uuid"

type CreateComponentRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Key   string `json:"key" binding:"required" default:""`
	Value string `json:"value" binding:"required"`
}

type CreateMenuComponentRequest struct {
	ID    uuid.UUID
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Key   string `json:"key" binding:"required" default:""`
	Value string `json:"value" binding:"required"`
	Order int    `json:"order" binding:"required"`
}
