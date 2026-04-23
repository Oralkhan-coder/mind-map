package transport

import (
	"net/http"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService NodeService
}

func NewNodeHandler(nodeService NodeService) *NodeHandler {
	return &NodeHandler{nodeService}
}

func (h *NodeHandler) CreateNode(c *gin.Context) {
	mapId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}
	req := new(dto.NodeCreateRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(core.BadRequest("invalid request body: " + err.Error()))
		return
	}

	res, err := h.nodeService.CreateNode(c.Request.Context(), req, userId.(string), mapId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *NodeHandler) GetByMapId(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}

	res, err := h.nodeService.GetByMapId(c.Request.Context(), id, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
