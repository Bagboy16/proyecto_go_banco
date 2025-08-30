package services

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bagboy16/banco/models"
	"github.com/bagboy16/banco/util"
	"gorm.io/gorm"
)

func AltaCliente(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Documento del Cliente: ")
	scanner.Scan()
	documento := strings.TrimSpace(scanner.Text())

	fmt.Print("Nombre del Cliente: ")
	scanner.Scan()
	nombre := strings.TrimSpace(scanner.Text())

	fmt.Print("Apellido del Cliente: ")
	scanner.Scan()
	apellido := strings.TrimSpace(scanner.Text())

	cliente := &models.Cliente{
		Documento: documento,
		Nombre:    nombre,
		Apellido:  apellido,
	}

	if err := db.Create(cliente).Error; err != nil {
		fmt.Println("Error al crear cliente:", err)
		fmt.Scanln()
		return
	}

	fmt.Println("Cliente creado exitosamente... Pulse ENTER para continuar")
	fmt.Scanln()

}

func ConsultarCliente(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Documento del Cliente: ")
	scanner.Scan()
	documento := strings.TrimSpace(scanner.Text())

	var cliente models.Cliente
	if err := db.First(&cliente, "documento = ?", documento).Error; err != nil {
		fmt.Println("Cliente no encontrado. Debe crear el cliente primero.")
		fmt.Scanln()
		return
	}
	util.ConsoleClear()
	fmt.Printf("Cliente: %s\n", cliente.Documento)
	fmt.Printf("Nombre: %s %s\n", cliente.Nombre, cliente.Apellido)

	var cuentas []models.Cuenta
	if err := db.Where("documento_cliente = ?", documento).Find(&cuentas).Error; err != nil {
		fmt.Println("Error al buscar cuentas:", err)
		fmt.Scanln()
		return
	}

	if len(cuentas) == 0 {
		fmt.Println("El cliente no tiene cuentas asociadas")
		fmt.Scanln()
		return
	}

	fmt.Println("Cuentas del cliente: ")
	for _, c := range cuentas {
		fmt.Printf("- Número: %s | Tipo: %s | Saldo: %.4f\n", c.Numero, c.Tipo, c.Saldo)
	}
	fmt.Scanln()
}
