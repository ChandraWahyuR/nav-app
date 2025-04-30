package main

import (
	"fmt"
	"log"
	"proyek1/app"
	"proyek1/config"
	"proyek1/db/migrations"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	serve := gin.Default()
	cfg := config.EnvFile()
	logger := logrus.New()
	// Koneksi ke database
	db, err := config.InitDatabase(*cfg)
	if err != nil {
		logger.Fatal("Gagal menghubungkan ke database:", err)
		return
	}

	// Inisialisasi tabel
	err = migrations.CreateTables(db)
	if err != nil {
		log.Fatal("Gagal membuat tabel:", err)
	}

	// Inisialisasi JWT (pakai secret dari env)
	jwt := config.NewJWT(logger)
	// Jalankan Bootstrap
	bootstrap := &app.BootstrapConfig{
		App: serve,
		DB:  db,
		Log: logger,
		JWT: jwt,
		Cfg: cfg,
	}
	app.App(bootstrap)

	// Jalankan server
	err = serve.Run("localhost:8081")
	if err != nil {
		fmt.Println("Server tidak bisa dijalankan:", err)
	}
}
