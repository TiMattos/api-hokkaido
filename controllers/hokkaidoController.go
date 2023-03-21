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

func ListarClientes(c *gin.Context) {
	var clientes []models.Cliente
	database.DB.Find(&clientes)
	c.JSON(200, clientes)
}

func ListarCarrosPorID(c *gin.Context) {
	var veiculos []models.Veiculo
	id := c.Params.ByName("id")
	database.DB.Where("cliente_id = ?", id).Find(&veiculos)

	if len(veiculos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Nehum veiculo localizado"})
		return
	}
	c.JSON(200, veiculos)
}
