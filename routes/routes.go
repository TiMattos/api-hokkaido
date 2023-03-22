package routes

import (
	"github.com/TiMattos/go-hokkaido/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.GET("/:nome", controllers.Saudacao)
	r.GET("/clientes", controllers.ListarClientes)
	r.GET("/veiculos/:id", controllers.ListarCarrosPorID)
	r.POST("/clientes", controllers.IncluirCliente)
	r.GET("/cliente/:id", controllers.BuscarClientePorID)
	r.GET("/cliente/nome/:nome", controllers.BuscarClientePorNome)
	r.Run()

}
