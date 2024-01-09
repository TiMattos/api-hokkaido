package controllers

import (
	"net/http"

	"github.com/TiMattos/go-hokkaido/database"
	"github.com/TiMattos/go-hokkaido/models"
	"github.com/TiMattos/go-hokkaido/pkg/logger"
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
		logger.GravarLog("Nehum veiculo localizado")
		return
	}
	c.JSON(200, veiculos)
}

func BuscarClientePorID(c *gin.Context) {
	var cliente models.Cliente
	id := c.Params.ByName("id")
	database.DB.First(&cliente, id)
	if cliente.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Cliente não localizado"})
		logger.GravarLog("Cliente não localizado")
		return
	}
	c.JSON(200, cliente)
}

func BuscarClientePorNome(c *gin.Context) {
	var cliente models.Cliente
	nome := c.Params.ByName("nome")
	database.DB.Where("nome = ?", nome).Find(&cliente)
	//database.DB.First(&cliente, nome)
	if cliente.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Cliente não localizado"})
		logger.GravarLog("Cliente não localizado")
		return
	}
	c.JSON(200, cliente)
}

func BuscarVeiculoPorPlaca(c *gin.Context) {
	var veiculo models.Veiculo
	placa := c.Params.ByName("placa")

	// Usar um mapa para as condições da consulta
	condicoes := map[string]interface{}{"placa": placa}

	// Alterar a consulta para usar o mapa de condições
	database.DB.Where(condicoes).First(&veiculo)

	if veiculo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Nenhum veiculo localizado",
		})
		logger.GravarLog("Nenhum veiculo localizado")
		return
	}

	c.JSON(http.StatusOK, veiculo)
}

func BuscarClienteEVeiculosPorNome(c *gin.Context) {
	var clientes []models.Cliente
	//var veiculos []models.Veiculo
	nome := c.Params.ByName("nome")

	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O parâmetro 'nome' é obrigatório",
		})
		return
	}
	database.DB.Where("nome LIKE ?", "%"+nome+"%").Find(&clientes)

	for i := range clientes {
		clientes[i].Veiculo = ListVeiculos(c, clientes[i].ID)
	}

	if len(clientes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nenhum cliente localizado com o nome fornecido",
		})
		logger.GravarLog("Nenhum cliente localizado com o nome fornecido")
		return
	}

	c.JSON(http.StatusOK, clientes)
}

func ListVeiculos(c *gin.Context, id int) []models.Veiculo {
	var veiculos []models.Veiculo
	database.DB.Where("cliente_id = ?", id).Find(&veiculos)
	return veiculos
}
