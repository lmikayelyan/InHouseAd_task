package handler

import (
	"InHouseAd/internal/model"
	"InHouseAd/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type Category interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
}

type category struct {
	svc    service.Category
	logger zerolog.Logger
}

func CategoryHandler(svc service.Category, logger zerolog.Logger) Category {
	return &category{svc: svc, logger: logger}
}

// Create
// @Summary Create category
// @Description Add good to the database of goods
// @Param categoryInput body model.CategoryInputResponse true "Category Info"
// @Produce json
// @Success 200 {object} model.CategoryInputResponse
// @Router /category/create [post]
func (h *category) Create(c *gin.Context) {
	var item model.Category
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Category.Create.BindJson: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Create.BindJson: %v", err)
		return
	}

	if err := h.svc.Create(c, &item); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Category.Create: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Create: %v", err)
		return
	}

	response := model.CategoryInputResponse{
		Name: item.Name,
	}

	c.IndentedJSON(http.StatusOK, response)
	h.logger.Info().Msgf("category %s created successfully", item.Name)
}

// Update
// @Summary Update existing good
// @Description Updates existing goods instance
// @Param id path int true "ID"
// @Param item body model.CategoryInputResponse true "Updated Category"
// @Produce json
// @Success 200 {object} model.CategoryInputResponse
// @Router /category/update/{id} [patch]
func (h *category) Update(c *gin.Context) {
	var item model.Category
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Category.Update.BindJson: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Update.BindJson: %v", err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.goods.Delete.Strconv(id): %v", err)})
		return
	}

	item.ID = uint(id)

	if err := h.svc.Update(c, &item); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Category.Update: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Create: %v", err)
		return
	}

	response := model.CategoryInputResponse{
		Name: item.Name,
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message":      "category update successful",
		"updated_item": response,
	})

	h.logger.Info().Msgf("category %s updated successfully", item.Name)
}

// Delete
// @Summary Delete category
// @Description Deletes category and updates all goods categories (removes deleted category)
// @Param id path int true "ID"
// @Produce json
// @Router /category/remove/{id} [delete]
func (h *category) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Category.Update.BindJson: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Get.Strconv: %v", err)
		return
	}

	if err := h.svc.Delete(c, uint(id)); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Category.Delete: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Delete: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "delete successful",
	})

	h.logger.Info().Msgf("category and it's according items deleted successfully")
}

// Get
// @Summary Get all categories
// @Description Get categories list
// @Produce json
// @Success 200
// @Router /category/list [get]
func (h *category) Get(c *gin.Context) {
	categoryList, err := h.svc.Get(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Category.Get: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.Get: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, categoryList)
	h.logger.Info().Msgf("category.get successful")
}
