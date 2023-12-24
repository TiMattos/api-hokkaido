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
	var veiculos []models.Veiculo
	placa := c.Params.ByName("placa")
	database.DB.First(&veiculos, placa)
	if len(veiculos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Nehum veiculo localizado"})
		logger.GravarLog("Nehum veiculo localizado")
		return
	}
	c.JSON(200, veiculos)

}

func BuscarClientePorNomeEVeiculos(c *gin.Context) {
	var clientes []models.Cliente
	nome := c.Params.ByName("nome")

	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Erro": "O parâmetro 'nome' é obrigatório",
		})
		return
	}

	database.DB.Where("nome LIKE ?", "%"+nome+"%").Find(&clientes)

	if len(clientes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Nenhum cliente localizado com o nome fornecido",
		})
		logger.GravarLog("Nenhum cliente localizado com o nome fornecido")
		return
	}
	c.JSON(http.StatusOK, clientes)
}

func BuscarClienteEVeiculosPorNome(c *gin.Context) {
	var clientes models.Cliente
	var veiculos []models.Veiculo
	nome := c.Params.ByName("nome")

	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O parâmetro 'nome' é obrigatório",
		})
		return
	}
	database.DB.Where("nome LIKE ?", "%"+nome+"%").Find(&clientes)

	clienteID, exists := c.Get("ID")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID do cliente não encontrado no contexto",
		})
		return
	}

	database.DB.Preload("Veiculos").Where("cliente_id =  ?", clienteID).First(&veiculos)
	// Consulta que inclui a tabela de veículos relacionada ao cliente
	//database.DB.Preload("Veiculos").Where("id =  ?", "%"+nome+"%").First(&clientes)

	if clientes.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nenhum cliente localizado com o nome fornecido",
		})
		logger.GravarLog("Nenhum cliente localizado com o nome fornecido")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cliente":  clientes,
		"veiculos": clientes.Veiculo, // Acessa os veículos associados ao cliente
	})
}
