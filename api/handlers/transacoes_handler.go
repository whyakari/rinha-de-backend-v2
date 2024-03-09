package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/whyakari/rinha-de-backend-v2/database"
)

func HandleTransacoes(c *gin.Context) {
    clienteID := c.Param("id")

    var requestBody struct {
        Valor     int    `json:"valor"`
        Tipo      string `json:"tipo"`
        Descricao string `json:"descricao"`
    }

    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(400, gin.H{"error": "Dados inválidos"})
        return
    }

    var limiteCliente int
    err := db.DB.QueryRow("SELECT limite FROM clientes WHERE id = ?", clienteID).Scan(&limiteCliente)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao obter o limite do cliente"})
        fmt.Println(err)
        return
    }

    var query string
    var args []interface{}

    if requestBody.Tipo == "c" {
        query = "INSERT INTO transacoes (id_cliente, valor, tipo, descricao, realizada_em) VALUES (?, ?, ?, ?, NOW())"
        args = []interface{}{clienteID, requestBody.Valor, requestBody.Tipo, requestBody.Descricao}
    } else if requestBody.Tipo == "d" {
        query = "INSERT INTO transacoes (id_cliente, valor, tipo, descricao, realizada_em) VALUES (?, ?, ?, ?, NOW())"
        args = []interface{}{clienteID, -requestBody.Valor, requestBody.Tipo, requestBody.Descricao}
    }

    _, err = db.DB.Exec(query, args...)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao inserir transação no banco de dados"})
        fmt.Println(err)
        return
    }

    var saldoAtual int
    err = db.DB.QueryRow("SELECT saldo FROM clientes WHERE id = ?", clienteID).Scan(&saldoAtual)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao consultar saldo do cliente"})
        fmt.Println(err)
        return
    }

    var saldoNovo int
    if requestBody.Tipo == "c" {
        saldoNovo = saldoAtual + requestBody.Valor
    } else if requestBody.Tipo == "d" {
        saldoNovo = saldoAtual - requestBody.Valor
    }

    _, err = db.DB.Exec("UPDATE clientes SET saldo = ? WHERE id = ?", saldoNovo, clienteID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao atualizar saldo do cliente"})
        fmt.Println(err)
        return
    }

    c.JSON(200, gin.H{
        "saldo": saldoNovo,
    })
}
