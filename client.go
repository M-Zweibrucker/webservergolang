package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Bid string `json:"bid"`
}

func requestCotacao() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Erro: timeout ao receber resposta do servidor (300ms excedido)")
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("servidor retornou erro: status %d", resp.StatusCode)
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Bid, nil
}

func saveToFile(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)
	return os.WriteFile("cotacao.txt", []byte(content), 0644)
}

func main() {
	bid, err := requestCotacao()
	if err != nil {
		fmt.Println("Error fetching cotacao:", err)
		return
	}

	if err := saveToFile(bid); err != nil {
		fmt.Println("Error saving to file:", err)
		return
	}

	fmt.Println("Cotacao saved successfully:", bid)
}
