package controllers

import (
	"net/http"

	"github.com/TiMattos/go-hokkaido/database"
	"github.com/TiMattos/go-hokkaido/models"
	"github.com/gin-gonic/gin"
)

func Saudacao(c *gin.Context) {
	nome := c.Params.ByName("nome")
	c.JSON(200, gin.H{
		"API diz:": "E ai " + nome + " tudo beleza?",
	})
}

func IncluirCliente(c *gin.Context) {
	var cliente models.Cliente
	// o ShouldBind transforma a requisição na model
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	database.DB.Create(&cliente)
	c.JSON(http.StatusOK, cliente)
}
