package database

import (
	"log"

	"github.com/TiMattos/go-hokkaido/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados() {
	stringDeConexao := "host=localhost user=hokkaido password=hokk@ido dbname=root port=5432 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(stringDeConexao))
	if err != nil {
		log.Panic("Erro ao conectar com o banco de dados.")
	}
	DB.AutoMigrate(&models.Cliente{})
	DB.AutoMigrate(&models.Servico{})
	DB.AutoMigrate(&models.Veiculo{})
}
