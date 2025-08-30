package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AbrirDB(user, password, dbname, host, port string) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	fmt.Println("⚙️ Inicializando conexión a DB con parámetros: ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	fmt.Print("⚙️ Conexión exitosa a Base de Datos\n")
	return db, nil
}
