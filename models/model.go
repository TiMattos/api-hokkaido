package models

import (
	"gorm.io/gorm"
)

type Cliente struct {
	gorm.Model
	Nome        string    `json:"nome"`
	Email       string    `json:"email"`
	Telefone    string    `json:"telefone"`
	Endereco    string    `json:"endereco"`
	Complemento string    `json:"complemento"`
	Numero      string    `json:"numero"`
	Bairro      string    `json:"bairro"`
	Veiculo     []Veiculo `json:"veiculo"`
	IDCliente   int
}

type Servico struct {
	gorm.Model
	Descricao  string `json:"descricao"`
	Observacao string `json:"observacao"`
	ValorMO    string `json:"valorMO"`
	ValorPecas string `json:"valorPecas"`
	IDCliente  int    `json:"idCliente"`
}

type Veiculo struct {
	gorm.Model
	Marca     string `json:"marca"`
	Modelo    string `json:"modelo"`
	Placa     string `json:"placa"`
	Ano       string `json:"ano"`
	ClienteID int    `json:"idCliente" gorm:"foreignKey:IDCliente"`
}