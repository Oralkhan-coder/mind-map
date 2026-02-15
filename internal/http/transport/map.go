package transport

import (
	"net/http"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type MapHandler struct {
	mapService MapService
}

func NewMapHandler(mapService MapService) *MapHandler {
	return &MapHandler{mapService}
}

func (h *MapHandler) GetMaps(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}

	res, err := h.mapService.GetMaps(c.Request.Context(), userId.(string))
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, res)
}

func (h *MapHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}

	res, err := h.mapService.GetByID(c.Request.Context(), id, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *MapHandler) CreateMap(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
	}

	var req dto.MapCURequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(core.BadRequest("invalid request"))
		return
	}

	id, err := h.mapService.CreateMap(c.Request.Context(), &req, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, id)
}

func (h *MapHandler) UpdateMap(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}

	var req dto.MapCURequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(core.BadRequest("invalid request"))
		return
	}

	err := h.mapService.UpdateMap(c.Request.Context(), id, userId.(string), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *MapHandler) DeleteMap(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(core.Unauthorized("user ID not found in context"))
		return
	}

	err := h.mapService.DeleteMap(c.Request.Context(), id, userId.(string))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
