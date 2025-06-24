package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func logger(ctx context.Context, msg string) {
	if err := ctx.Err(); err != nil {
		log.Printf("%s: %v", msg, err)
	}
}

func getCotacaoFromAPI() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		logger(ctx, "Failed to create request")
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger(ctx, "Erro: timeout ao chamar API de cotação (200ms excedido)")
		return "", err
	}
	defer resp.Body.Close()

	var c Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		logger(ctx, "Failed to decode response")
		return "", err
	}

	return c.USDBRL.Bid, nil
}

func initDB() error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	return nil
}

func saveToDB(bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		logger(ctx, "Failed to open database")
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		logger(ctx, "Erro: timeout ao persistir no banco de dados (10ms excedido)")
		return err
	}

	return nil
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	bid, err := getCotacaoFromAPI()
	if err != nil {
		http.Error(w, "Failed to get cotacao", http.StatusInternalServerError)
		return
	}

	if err := saveToDB(bid); err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": bid})
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal("Erro ao inicializar banco de dados:", err)
	}

	http.HandleFunc("/cotacao", cotacaoHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
