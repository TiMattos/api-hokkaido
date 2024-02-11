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
	Descricao  string `json:"descricao"`
	Observacao string `json:"observacao"`
	ValorMO    string `json:"valorMO"`
	ValorPecas string `json:"valorPecas"`
	KmAtual    string `json:"kmAtual"`
	KmRevisao  string `json:"kmRevisao"`
	ClienteID  int    `gorm:"foreignKey:ID"`
	VeiculoID  int    `gorm:"foreignKey:ID"`
}

type ServicoResponse struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Descricao  string     `json:"descricao"`
	Observacao string     `json:"observacao"`
	ValorMO    string     `json:"valorMO"`
	ValorPecas string     `json:"valorPecas"`
	KmAtual    string     `json:"kmAtual"`
	KmRevisao  string     `json:"kmRevisao"`
	ClienteID  int        `json:"cliente_id"`
	VeiculoID  int        `json:"veiculo_id"`
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
