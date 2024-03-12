package routes

import (
	"github.com/TiMattos/go-hokkaido/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.POST("/login", controllers.LoginHandler)
	authGroup := r.Group("/")
	authGroup.Use(controllers.Authenticate)
	authGroup.GET("/clientes", controllers.ListarClientes)
	authGroup.GET("/veiculos/id/:id", controllers.ListarCarrosPorID)
	authGroup.GET("/veiculos/:placa", controllers.BuscarVeiculoPorPlaca)
	authGroup.POST("/clientes", controllers.IncluirCliente)
	authGroup.GET("/cliente/:id", controllers.BuscarClientePorID)
	authGroup.GET("/cliente/nome/:nome", controllers.BuscarClientePorNome)
	authGroup.GET("/cliente/find/:nome", controllers.BuscarClienteEVeiculosPorNome)
	authGroup.POST("/IncluirVeiculo", controllers.IncluirVeiculo)
	authGroup.POST("/IncluirServico", controllers.IncluirServico)
	authGroup.GET("/servico/lista/:id", controllers.ListarServicosPorIdCliente)
	authGroup.GET("/servico/:id", controllers.BuscaServicoPorID)
	authGroup.GET("/health", controllers.HealthCheck)
	authGroup.POST("/cliente/update", controllers.UpdateCliente)
	authGroup.GET("/cliente/emails", controllers.ListEmailableClients)
	r.Run()

}
