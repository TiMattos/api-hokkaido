package routes

import (
	"github.com/TiMattos/go-hokkaido/controllers"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.GET("/:nome", controllers.Saudacao)
	r.POST("/clientes", controllers.IncluirCliente)
	r.Run()

}
