package services

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bagboy16/banco/models"
	"gorm.io/gorm"
)

func AltaCuenta(scanner *bufio.Scanner, db *gorm.DB) {
	fmt.Print("Ingrese documento del cliente: \n")
	scanner.Scan()
	documentoCliente := strings.TrimSpace(scanner.Text())

	var cliente models.Cliente
	if err := db.First(&cliente, "documento = ?", documentoCliente).Error; err != nil {
		fmt.Println("Cliente no encontrado. Debe crear el cliente primero.")
		fmt.Scanln()
		return
	}

	fmt.Print("Ingrese tipo de cuenta (Ahorro/Corriente): \n")
	scanner.Scan()
	tipoCuenta := strings.TrimSpace(scanner.Text())

	cuenta := models.Cuenta{
		DocumentoCliente: documentoCliente,
		Tipo:             tipoCuenta,
		Saldo:            0,
	}

	if err := db.Create(&cuenta).Error; err != nil {
		fmt.Println("Error al crear la cuenta:", err)
		fmt.Scanln()
		return
	}

	fmt.Printf("Cuenta %s creada exitosamente\n", cuenta.Numero)
	fmt.Scanln()

}
