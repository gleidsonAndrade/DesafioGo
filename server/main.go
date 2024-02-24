package main

import (
	_ "bytes"
	"context"
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}
type Cambio struct {
	ID         int    `gorm:"primarykey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)

}

func insertCambio(db *gorm.DB, cambio Cambio) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db.WithContext(ctx).Create(&cambio)
	//stmt, error := db.Prepare("INSERT INTO cambio (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	//if error != nil {
	//	fmt.Println("Ocorreu um erro:", error)
	//	return error
	//}
	//defer stmt.Close()
	//_, error = stmt.Exec(cambio.Usdbrl.Code, cambio.Usdbrl.Codein, cambio.Usdbrl.Name, cambio.Usdbrl.High, cambio.Usdbrl.Low, cambio.Usdbrl.VarBid, cambio.Usdbrl.PctChange, cambio.Usdbrl.Bid, cambio.Usdbrl.Ask, cambio.Usdbrl.Timestamp, cambio.Usdbrl.CreateDate)
	//if error != nil {
	//	fmt.Println("Ocorreu um erro:", error)
	//	return error
	//}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

	log.Println("requisição iniciada!")
	defer log.Println("requisição finalizada!")
	//w.Write([]byte("requisição processada com sucesso!"))

	defer cancel()
	//req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Ocorreu um erro de time out:", err)
		return
	}

	//res, error := io.ReadAll(req.Body)
	res, error := http.DefaultClient.Do(req)
	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
		return
	}
	defer res.Body.Close()

	les, error := io.ReadAll(res.Body)
	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
		return
	}

	var cot Cotacao

	error = json.Unmarshal([]byte(les), &cot)
	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
		return
	}
	fmt.Println(cot)

	//db, error := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, error := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if error != nil {
		fmt.Println("Ocorreu um erro:", error)
		return
	}
	db.AutoMigrate(&Cambio{})
	var cambio Cambio
	cambio.Code = cot.Usdbrl.Code
	cambio.Codein = cot.Usdbrl.Codein
	cambio.Name = cot.Usdbrl.Name
	cambio.High = cot.Usdbrl.High
	cambio.Low = cot.Usdbrl.Low
	cambio.VarBid = cot.Usdbrl.VarBid
	cambio.PctChange = cot.Usdbrl.PctChange
	cambio.Bid = cot.Usdbrl.Bid
	cambio.Ask = cot.Usdbrl.Ask
	cambio.Timestamp = cot.Usdbrl.Timestamp
	cambio.CreateDate = cot.Usdbrl.CreateDate
	error = insertCambio(db, cambio)

	json.NewEncoder(w).Encode(cot.Usdbrl.Bid)
	//if error != nil {
	//	fmt.Println("Ocorreu um erro ao inserir cambio:", error)
	//	return
	//}
}
