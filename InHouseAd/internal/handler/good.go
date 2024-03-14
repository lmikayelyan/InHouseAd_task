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

type Good interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByCategory(c *gin.Context)
}

type good struct {
	svc    service.Good
	logger zerolog.Logger
}

func GoodHandler(svc service.Good, logger zerolog.Logger) Good {
	return &good{svc: svc, logger: logger}
}

// Create
// @Summary Create good
// @Description Add good to the database of goods.
// @Param goodInput body model.GoodInputResponse true "Good Info"
// @Produce json
// @Success 200 {object} model.GoodInputResponse
// @Router /good/create [post]
func (h *good) Create(c *gin.Context) {
	var item model.Good
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Good.Create.BindJson: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Create.BindJson: %v", err)
		return
	}

	if err := h.svc.Create(c, &item); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Good.Create: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Create: %v", err)
		return
	}

	reponse := model.GoodInputResponse{
		Name:       item.Name,
		Categories: item.Categories,
	}

	c.IndentedJSON(http.StatusOK, reponse)
	h.logger.Info().Msgf("Good %s created successfully", item.Name)
}

// Update
// @Summary Update existing good
// @Description Updates existing goods instance
// @Param id path int true "ID"
// @Param item body model.GoodInputResponse true "Updated good"
// @Produce json
// @Success 200 {object} model.GoodInputResponse
// @Router /good/update/{id} [patch]
func (h *good) Update(c *gin.Context) {
	var item model.Good
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Good.Update.BindJson: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Update.BindJson: %v", err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.goods.Update.Strconv(id): %v", err)})
		return
	}

	item.ID = uint(id)

	if err := h.svc.Update(c, &item); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Good.Update: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Create: %v", err)
		return
	}

	response := model.GoodInputResponse{
		Name:       item.Name,
		Categories: item.Categories,
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":      "good update successful",
		"updated_item": response,
	})

	h.logger.Info().Msgf("Good %s updated successfully", item.Name)
}

// Delete
// @Summary Delete good
// @Description Deletes good instance
// @Param id path int true "ID"
// @Produce json
// @Router /good/remove/{id} [delete]
func (h *good) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Good.Update.Strconv: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Get.Strconv: %v", err)
		return
	}

	if err := h.svc.Delete(c, uint(id)); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Good.Delete: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.Delete: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "delete successful",
	})

	h.logger.Info().Msgf("Good %d deleted successfully", id)
}

// GetByCategory
// @Summary Get goods by category ID
// @Description Get goods list by the inputted category
// @Param category_id path int true "Category ID"
// @Produce json
// @Success 200
// @Router /good/list/{category_id} [get]
func (h *good) GetByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Errorf("handler.Good.GetByCategory.Strconv: %v", err),
		})
		h.logger.Error().Msgf("handler.Good.GetByCategory.Strconv: %v", err)
		return
	}

	categoryList, err := h.svc.GetByCategory(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    http.StatusBadGateway,
			"message": fmt.Errorf("handler.Category.GetByCategory: %v", err),
		})
		h.logger.Error().Msgf("handler.Category.GetByCategory: %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, categoryList)
	h.logger.Info().Msgf("category.get successful")
}
