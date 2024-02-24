package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, error := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)

	if error != nil {
		panic(error)
	}
	res, error := http.DefaultClient.Do(req)
	if error != nil {
		panic(error)
	}
	defer res.Body.Close()
	file, error := os.Create("arquivo.txt")
	if error != nil {
		panic(error)
	}
	defer file.Close()
	_, error = io.Copy(file, res.Body)
	if error != nil {
		panic(error)
	}
	fmt.Println("Arquivo criado com sucesso!")

	//io.Copy(os.Stdout, res.Body)
}
