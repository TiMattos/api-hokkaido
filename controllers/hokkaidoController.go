package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/TiMattos/go-hokkaido/database"
	"github.com/TiMattos/go-hokkaido/models"
	"github.com/TiMattos/go-hokkaido/pkg/logger"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey []byte

func generateSecretKeyFromPassword(password string, length int) (string, error) {
	hash := sha256.New()
	io.WriteString(hash, password)
	derivedKey := hash.Sum(nil)
	secretKey := base64.URLEncoding.EncodeToString(derivedKey[:length])
	return secretKey, nil
}

func init() {
	database.ConectaComBancoDeDados()
	var apiCredentials models.ApiCredentials
	if err := database.DB.First(&apiCredentials).Error; err != nil {
		fmt.Println("Erro ao obter as credenciais do banco de dados:", err)
		return
	}

	derivedKey, err := generateSecretKeyFromPassword(apiCredentials.SecretKey, 32)
	if err != nil {
		fmt.Println("Erro ao gerar a chave secreta:", err)
		return
	}

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

	var apiCredentials models.ApiCredentials
	database.DB.First(&apiCredentials)

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
	// Inserir cliente no banco de dados
	if err := database.DB.Create(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao inserir cliente",
		})
		return
	}

	clienteJSON, err := structToJSON(cliente)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao converter cliente para JSON",
		})
		return
	}

	log := models.LogCliente{
		ClienteID: cliente.ID,
		Operacao:  "Insert",
		Usuario:   "admin",
		Json:      clienteJSON,
	}
	if err := InsertLogCliente(c, log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao inserir log de cliente",
		})
		return
	}

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
	logger.GravarLog("SearchCustomerAndVehiclesByName: Iniciando busca de cliente por nome")

	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O parâmetro 'nome' é obrigatório",
		})
		return
	}
	logger.GravarLog("SearchCustomerAndVehiclesByName: Buscando cliente por nome")

	if err := database.DB.Where("nome LIKE ?", "%"+nome+"%").Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar cliente",
		})
		logger.GravarLog("SearchCustomerAndVehiclesByName: Erro ao buscar cliente")
		return
	}

	for i := range clientes {
		clientes[i].Veiculo = ListVeiculos(c, clientes[i].ID)
	}

	if len(clientes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Nenhum cliente localizado com o nome fornecido",
		})
		logger.GravarLog("SearchCustomerAndVehiclesByName: Nenhum cliente localizado com o nome fornecido")
		return
	}
	logger.GravarLog("SearchCustomerAndVehiclesByName: Cliente localizado com sucesso")
	c.JSON(http.StatusOK, clientes)
}

func ListVeiculos(c *gin.Context, id int) []models.Veiculo {
	var veiculos []models.Veiculo
	logger.GravarLog("ListVehicles: Iniciando busca de veículos por ID do cliente")
	if err := database.DB.Where("cliente_id = ?", id).Find(&veiculos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar veículos",
		})
		logger.GravarLog("ListVehicles: Erro ao buscar veículos")
		return nil
	}
	logger.GravarLog("ListVehicles: Veículos localizados com sucesso")
	return veiculos
}

func HealthCheck(c *gin.Context) {
	db, err := database.DB.DB()
	logger.GravarLog("HealthCheck: Iniciando verificação de saúde do banco de dados")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Erro ao obter o objeto de banco de dados",
		})
		return
	}
	logger.GravarLog("HealthCheck: Verificando conexão com o banco de dados")
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Falha na conexão com o banco de dados",
		})
		logger.GravarLog("HealthCheck: Falha na conexão com o banco de dados")
		return
	}
	logger.GravarLog("HealthCheck: Banco de dados está saudável")

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func ListarServicosPorIdCliente(c *gin.Context) {
	var servicos []models.Servico
	logger.GravarLog("ListServicesByClientId: Iniciando busca de serviços por ID do cliente")
	id := c.Params.ByName("id")
	if err := database.DB.Where("cliente_id = ?", id).Find(&servicos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar serviços",
		})
		logger.GravarLog("ListServicesByClientId: Erro ao buscar serviços")
		return
	}

	logger.GravarLog("ListServicesByClientId: Serviços localizados com sucesso")
	c.JSON(200, servicos)
}

func BuscaServicoPorID(c *gin.Context) {
	var servico models.Servico
	logger.GravarLog("SearchServiceById: Iniciando busca de serviço por ID")
	id := c.Params.ByName("id")
	database.DB.First(&servico, id)
	if servico.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Serviço não localizado"})
		logger.GravarLog("SearchServiceById: Serviço não localizado")
		return
	}
	logger.GravarLog("SearchServiceById: Serviço localizado com sucesso")
	c.JSON(200, servico)
}

func UpdateCliente(c *gin.Context) {
	var cliente models.Cliente
	logger.GravarLog("UpdateClient: Iniciando atualização de cliente")
	id := c.Params.ByName("id")
	database.DB.First(&cliente, id)
	if cliente.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Cliente não localizado"})
		logger.GravarLog("Cliente não localizado")
		return
	}
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	if err := database.DB.Save(&cliente).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar cliente",
		})
		return
	}
	logger.GravarLog("UpdateClient: cliente atualizado com sucesso")

	clienteJSON, err := structToJSON(cliente)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao converter cliente para JSON",
		})
		return
	}

	log := models.LogCliente{
		ClienteID: cliente.ID,
		Operacao:  "Update",
		Usuario:   "admin",
		Json:      clienteJSON,
	}
	if err := InsertLogCliente(c, log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao inserir log de cliente",
		})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

func InsertLogCliente(c *gin.Context, logCliente models.LogCliente) error {
	logger.GravarLog("InsertLog: Iniciando inserção de log de cliente")
	if err := database.DB.Create(&logCliente).Error; err != nil {
		return err
	}
	logger.GravarLog("InsertLog: log de cliente inserido com sucesso")
	return nil

}

func InsertLogVeiculo(c *gin.Context, logVeiculo models.LogVeiculo) error {

	if err := database.DB.Create(&logVeiculo).Error; err != nil {
		return err
	}

	return nil

}

func structToJSON(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func ListEmailableClients(c *gin.Context) {
	var clientes []models.Cliente
	var emails []string

	if err := database.DB.Where("email IS NOT NULL").Find(&clientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar clientes",
		})
		logger.GravarLog("listEmailableClients: Erro ao listar emails.")
		return
	}

	for _, cliente := range clientes {
		// Validar se o email é válido antes de adicioná-lo à lista
		if err := checkmail.ValidateFormat(cliente.Email); err == nil {
			emails = append(emails, cliente.Email)
		}
	}

	c.JSON(http.StatusOK, emails)
}

func ListDailyServices(c *gin.Context) {
	var servicos []models.Servico
	data := c.Param("data")

	logger.GravarLog("ListDailyServices: Iniciando busca de serviços diários")
	// Construir a consulta SQL para comparar a data truncada com a data no banco de dados
	if err := database.DB.Where("data_revisao = ?", data).Find(&servicos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar serviços",
		})
		logger.GravarLog("ListDailyServices: Erro ao buscar serviços")
		return
	}

	logger.GravarLog("ListDailyServices: Serviços localizados com sucesso")
	c.JSON(http.StatusOK, servicos)
}
