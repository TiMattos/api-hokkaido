package routes

import (
	"github.com/TiMattos/go-hokkaido/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.GET("/clientes", controllers.ListarClientes)
	r.GET("/veiculos/id/:id", controllers.ListarCarrosPorID)
	r.GET("/veiculos/:placa", controllers.BuscarVeiculoPorPlaca)
	r.POST("/clientes", controllers.IncluirCliente)
	r.GET("/cliente/:id", controllers.BuscarClientePorID)
	r.GET("/cliente/nome/:nome", controllers.BuscarClientePorNome)
	r.GET("/cliente/find/:nome", controllers.BuscarClienteEVeiculosPorNome)
	r.POST("/IncluirVeiculo", controllers.IncluirVeiculo)
	r.GET("/health", controllers.HealthCheck)
	r.Run()

}
