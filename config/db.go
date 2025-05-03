package config

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
)

var DB *sql.DB

func ConectDataBase() {
	var err error

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro no ping do banco: %v", err)
	}

	log.Println("Banco de dados conectado com sucesso")
}
