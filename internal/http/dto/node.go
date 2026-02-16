package dto

import (
	"github.com/Oralkhan-coder/mind-map/internal/model"
)

type NodeCreateRequest struct {
	ParentID string         `json:"parent_id,omitempty"`
	Type     string         `json:"type"`
	Position model.Point    `json:"position"`
	Data     model.NodeData `json:"data"`
	Width    float64        `json:"width,omitempty"`
	Height   float64        `json:"height,omitempty"`
}
