package controllers

import (
	"api/db"
	"api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
)

func DocumentsIndex(c *gin.Context) {
	var documents []models.Document
	db.DB.
		Preload("CreatedBy.Account").
		Preload(clause.Associations).
		Find(&documents)

	c.Header("Content-Type", "application/vnd.api+json; charset=utf-8")
	c.JSON(200, gin.H{
		"documents": documents,
	})
}

func DocumentsCreate(c *gin.Context) {
	d := models.Document{}

	err := c.Bind(&d)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result := db.DB.Create(&d)
	if result.Error != nil {
		_ = c.AbortWithError(http.StatusBadRequest, result.Error)
		return
	}

	c.Header("Content-Type", "application/vnd.api+json; charset=utf-8")
	c.JSON(http.StatusCreated, gin.H{
		"document": d,
	})
}
