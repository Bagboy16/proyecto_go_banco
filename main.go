package main

import (
	"bufio"
	"fmt"

	"os"

	"strings"

	"github.com/bagboy16/banco/db"
	"github.com/bagboy16/banco/models"
	"github.com/bagboy16/banco/services"
	"github.com/bagboy16/banco/util"
)

func main() {
	util.ConsoleClear()
	user := "postgres"
	password := "123456789"
	dbname := "banco"
	host := "localhost"
	port := "5432"

	conexion, err := db.AbrirDB(user, password, dbname, host, port)

	if err != nil {
		panic("⚠️ No se pudo conectar a la DB: " + err.Error())
	}

	if err := models.AutoMigraryLlenar(conexion); err != nil {
		panic("⚠️ No se pudo migrar a la DB: " + err.Error())
	}

	fmt.Println("Conexión exitosa")
	util.ConsoleClear()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		util.ConsoleClear()
		fmt.Println("\n === SISTEMA BANCARIO ===")
		fmt.Println("1 - Alta de Cliente")
		fmt.Println("2 - Alta de Cuenta")
		fmt.Println("3 - Depósito")
		fmt.Println("4 - Retiro")
		fmt.Println("5 - Transferencia")
		fmt.Println("6 - Consultar cuenta")
		fmt.Println("7 - Consultar transacciones")
		fmt.Println("0 - Salir")
		fmt.Print("Escoje una opción: ")
		scanner.Scan()
		opcion := strings.TrimSpace(scanner.Text())

		switch opcion {
		case "1":
			services.AltaCliente(scanner, conexion)
		case "2":
			services.AltaCuenta(scanner, conexion)
		case "3":
			services.HacerDeposito(scanner, conexion)
		case "4":
			services.HacerRetiro(scanner, conexion)
		case "5":
			services.HacerTransferencia(scanner, conexion)
		case "6":
			services.ConsultarCliente(scanner, conexion)
		case "7":
			services.ConsultarTransacciones(scanner, conexion)
		default:
			fmt.Println("Escoja una opción correcta.")
			fmt.Scanln()
		}

	}

}
