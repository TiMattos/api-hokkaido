package models

import (
	"time"

	"gorm.io/gorm"
)

type Cliente struct {
	gorm.Model
	Nome           string    `json:"nome"`
	Email          string    `json:"email"`
	Telefone       string    `json:"telefone"`
	Endereco       string    `json:"endereco"`
	Complemento    string    `json:"complemento"`
	Cidade         string    `json:"cidade"`
	Bairro         string    `json:"bairro"`
	Veiculo        []Veiculo `json:"veiculo"`
	Cep            string    `json:"numeroCep"`
	DataNascimento string    `json:"dataNascimento"`
	ID             int       `gorm:"primaryKey"`
	Cpf            string    `json:"cpf"`
}

type Servico struct {
	gorm.Model
	Descricao   string    `json:"descricao"`
	Observacao  string    `json:"observacao"`
	ValorMO     string    `json:"valorMO"`
	ValorPecas  string    `json:"valorPecas"`
	KmAtual     string    `json:"kmAtual"`
	KmRevisao   string    `json:"kmRevisao"`
	ClienteID   int       `gorm:"foreignKey:ID"`
	VeiculoID   int       `gorm:"foreignKey:ID"`
	ID          int       `gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	DataRevisao string    `json:"dataRevisao"`
}

func (s *Servico) BeforeSave(tx *gorm.DB) (err error) {
	if s.ID == 0 {
		// Esta é uma inserção, não atualize CreatedAt
		return
	}

	// Esta é uma atualização, atualize CreatedAt
	s.CreatedAt = time.Now()
	return
}

type Veiculo struct {
	gorm.Model
	Marca     string `json:"marca"`
	Modelo    string `json:"modelo"`
	Placa     string `json:"placa"`
	Ano       string `json:"ano"`
	ClienteID int    `gorm:"foreignKey:ID"`
	ID        int    `gorm:"primaryKey"`
}

type ApiCredentials struct {
	gorm.Model
	User      string `json:"username"`
	Password  string `json:"password"`
	SecretKey string `json:"secretKey"`
}

type LogCliente struct {
	gorm.Model
	ClienteID int    `json:"clienteID"`
	Operacao  string `json:"operacao"`
	Usuario   string `json:"usuario"`
	Json      string `json:"json"`
}

type LogVeiculo struct {
	gorm.Model
	VeiculoID int    `json:"veiculoID"`
	Operacao  string `json:"operacao"`
	Usuario   string `json:"usuario"`
	ClienteID int    `json:"clienteID"`
}
