package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/TiMattos/go-hokkaido/database"
	"github.com/TiMattos/go-hokkaido/models"
	"github.com/TiMattos/go-hokkaido/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey []byte

func generateSecretKeyFromPassword(password string, length int) (string, error) {
	// Use SHA-256 para derivar a chave secreta
	hash := sha256.New()
	io.WriteString(hash, password)
	derivedKey := hash.Sum(nil)

	// Codifique os bytes derivados em uma string base64
	secretKey := base64.URLEncoding.EncodeToString(derivedKey[:length])
	return secretKey, nil
}

func init() {
	// Conecte-se ao banco de dados e obtenha as credenciais da API
	database.ConectaComBancoDeDados()
	var apiCredentials models.ApiCredentials
	if err := database.DB.First(&apiCredentials).Error; err != nil {
		fmt.Println("Erro ao obter as credenciais do banco de dados:", err)
		return
	}

	// Gere a chave secreta a partir da senha
	derivedKey, err := generateSecretKeyFromPassword(apiCredentials.SecretKey, 32)
	if err != nil {
		fmt.Println("Erro ao gerar a chave secreta:", err)
		return
	}

	// Atribua a chave secreta derivada à variável global
	secretKey = []byte(derivedKey)

	fmt.Println("Chave secreta gerada com sucesso:", secretKey)
}

func LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Lógica de autenticação - substitua isso com a lógica real
	// Por exemplo, você pode verificar as credenciais no banco de dados
	var apiCredentials models.ApiCredentials
	database.DB.First(&apiCredentials)

	// e gerar um token JWT se as credenciais estiverem corretas.
	if credentials.Username == apiCredentials.User && credentials.Password == apiCredentials.Password {
		token, dateExpiration, err := generateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token, "expiration": dateExpiration})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
}

func generateToken() (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTime.Unix() // Expira em 24 horas

	// Adicione logs para verificar a chave secreta antes de tentar assinar o token
	fmt.Println("Chave secreta utilizada para assinar o token:", secretKey)

	tokenString, err := token.SignedString(secretKey)

	fmt.Println("token gerado:", err, tokenString)

	if err != nil {
		return "", time.Time{}, fmt.Errorf("Erro ao assinar o token: %v", err)
	}

	return tokenString, expirationTime, nil
}

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	fmt.Println("tokenString", tokenString)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação ausente"})
		c.Abort()
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse do token com a mesma chave usada para assiná-lo
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	fmt.Println("token", token, "err", err)

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação inválido"})
		c.Abort()
		return
	}
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

func IncluirVeiculo(c *gin.Context) {
	var veiculo models.Veiculo
	if err := c.ShouldBindJSON(&veiculo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	database.DB.Create(&veiculo)
	c.JSON(http.StatusOK, veiculo)
}

func IncluirServico(c *gin.Context) {
	var servico models.Servico

	if err := c.ShouldBindJSON(&servico); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	database.DB.Create(&servico)
	c.JSON(http.StatusOK, servico)
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

func HealthCheck(c *gin.Context) {
	db, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Erro ao obter o objeto de banco de dados",
		})
		return
	}

	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Falha na conexão com o banco de dados",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func ListarServicosPorIdCliente(c *gin.Context) {
	var servicos []models.Servico
	id := c.Params.ByName("id")
	database.DB.Where("cliente_id = ?", id).Find(&servicos)
	c.JSON(200, servicos)
}

func BuscaServicoPorID(c *gin.Context) {
	var servico models.Servico
	id := c.Params.ByName("id")
	database.DB.First(&servico, id)
	if servico.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Serviço não localizado"})
		logger.GravarLog("Serviço não localizado")
		return
	}
	c.JSON(200, servico)
}
