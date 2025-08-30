package services

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/bagboy16/banco/models"
	"github.com/bagboy16/banco/util"
	"gorm.io/gorm"
)

func HacerDeposito(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Numero de cuenta: ")
	scanner.Scan()
	numCuenta := strings.TrimSpace(scanner.Text())

	var cuenta models.Cuenta

	if err := db.First(&cuenta, "numero = ?", numCuenta).Error; err != nil {
		fmt.Println("Cuenta no encontrada")
		fmt.Scanln()
		return
	}

	fmt.Print("Monto de la transacción: ")
	scanner.Scan()
	var monto float64
	_, err := fmt.Sscan(scanner.Text(), &monto)
	if err != nil || monto <= 0 {
		fmt.Println("Monto inválido")
		fmt.Scanln()
		return
	}

	cuenta.Saldo += monto

	if err := db.Save(&cuenta).Error; err != nil {
		fmt.Println("Error al actualizar saldo:", err)
		fmt.Scanln()
		return
	}

	transaccion := &models.Transaccion{
		Fecha:        time.Now(),
		NumeroCuenta: cuenta.Numero,
		TipoID:       1,
		Monto:        monto,
		Descripcion:  "Depósito",
	}

	if err := db.Create(transaccion).Error; err != nil {
		fmt.Println("Depósito realizado. Erro al registrar transacción:", err)
		fmt.Scanln()
		return
	}

	fmt.Println("Depósito exitoso. Saldo acutal: ", cuenta.Saldo)
	fmt.Scanln()

}

func HacerRetiro(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Numero de cuenta: ")
	scanner.Scan()
	numCuenta := strings.TrimSpace(scanner.Text())

	var cuenta models.Cuenta

	if err := db.First(&cuenta, "numero = ?", numCuenta).Error; err != nil {
		fmt.Println("Cuenta no encontrada")
		fmt.Scanln()
		return
	}

	fmt.Printf("Saldo: %f\n", cuenta.Saldo)
	fmt.Print("Monto de la transacción: ")
	scanner.Scan()
	var monto float64
	_, err := fmt.Sscan(scanner.Text(), &monto)
	if err != nil || monto <= 0 || monto > cuenta.Saldo {
		fmt.Println("Monto inválido")
		fmt.Scanln()
		return
	}

	cuenta.Saldo -= monto

	if err := db.Save(&cuenta).Error; err != nil {
		fmt.Println("Error al actualizar saldo:", err)
		fmt.Scanln()
		return
	}

	transaccion := &models.Transaccion{
		Fecha:        time.Now(),
		NumeroCuenta: cuenta.Numero,
		TipoID:       2,
		Monto:        monto,
		Descripcion:  "Retiro por taquilla",
	}

	if err := db.Create(transaccion).Error; err != nil {
		fmt.Println("Retiro realizado. Error al registrar transacción:", err)
		fmt.Scanln()
		return
	}

	fmt.Println("Retiro exitoso. Saldo actual: ", cuenta.Saldo)
	fmt.Scanln()

}

func HacerTransferencia(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Numero de cuenta Origen: ")
	scanner.Scan()
	numeroCuentaOrigen := strings.TrimSpace(scanner.Text())

	var cuentaOrigen models.Cuenta

	if err := db.First(&cuentaOrigen, "numero = ?", numeroCuentaOrigen).Error; err != nil {
		fmt.Println("Cuenta origen no encontrada")
		fmt.Scanln()
		return
	}

	fmt.Print("Numero de cuenta Destino: ")
	scanner.Scan()
	numeroCuentaDestino := strings.TrimSpace(scanner.Text())

	if numeroCuentaDestino == numeroCuentaOrigen {
		fmt.Println("La cuenta destino debe ser diferente a la cuenta origen")
		fmt.Scanln()
		return
	}

	var cuentaDestino models.Cuenta

	if err := db.First(&cuentaDestino, "numero = ?", numeroCuentaDestino).Error; err != nil {
		fmt.Println("Cuenta destino no encontrada")
		fmt.Scanln()
		return
	}

	fmt.Printf("Saldo: %f\n", cuentaOrigen.Saldo)
	fmt.Print("Monto de la transacción: ")
	scanner.Scan()
	var monto float64
	_, err := fmt.Sscan(scanner.Text(), &monto)
	if err != nil || monto <= 0 || monto > cuentaOrigen.Saldo {
		fmt.Println("Monto inválido")
		fmt.Scanln()
		return
	}

	err = db.Transaction(
		func(tx *gorm.DB) error {
			cuentaOrigen.Saldo -= monto
			if err := tx.Save(&cuentaOrigen).Error; err != nil {
				return err
			}
			cuentaDestino.Saldo += monto
			if err := tx.Save(&cuentaDestino).Error; err != nil {
				return err
			}
			transOrigen := &models.Transaccion{
				Fecha:        time.Now(),
				NumeroCuenta: cuentaOrigen.Numero,
				TipoID:       3,
				Monto:        -monto,
				Descripcion:  fmt.Sprintf("Transferencia a %s", cuentaDestino.Numero),
			}
			if err := tx.Create(transOrigen).Error; err != nil {
				return err
			}

			transDestino := &models.Transaccion{
				Fecha:        time.Now(),
				NumeroCuenta: cuentaDestino.Numero,
				TipoID:       3,
				Monto:        monto,
				Descripcion:  fmt.Sprintf("Transferencia desde %s", cuentaOrigen.Numero),
			}
			if err := tx.Create(transDestino).Error; err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		fmt.Println("Error al realizar la transferencia:", err)
		fmt.Scanln()
		return
	}

	fmt.Println("Transferencia exitosa")
	fmt.Scanln()

}

func ConsultarTransacciones(scanner *bufio.Scanner, db *gorm.DB) {
	util.ConsoleClear()
	fmt.Print("Numero de cuenta: ")
	scanner.Scan()
	numCuenta := strings.TrimSpace(scanner.Text())

	var cuenta models.Cuenta

	if err := db.Preload("Cliente").First(&cuenta, "numero = ?", numCuenta).Error; err != nil {
		fmt.Println("Cuenta no encontrada")
		fmt.Scanln()
		return
	}

	util.ConsoleClear()
	fmt.Printf("Cliente: %s\n", cuenta.Cliente.Documento)

	var transacciones []models.Transaccion
	if err := db.Preload("Tipo").Where("numero_cuenta = ?", numCuenta).Order("fecha desc").Find(&transacciones).Error; err != nil {
		fmt.Println("Error al buscar transacciones:", err)
		fmt.Scanln()
		return
	}

	if len(transacciones) == 0 {
		fmt.Println("No hay transacciones hechas en esta cuenta")
		fmt.Scanln()
		return
	}

	fmt.Println("Transacciones: ")
	for _, t := range transacciones {
		fmt.Printf("- Fecha: %s | Tipo: %s | Monto: %.4f | %s\n", t.Fecha, t.Tipo.Nombre, t.Monto, t.Descripcion)
	}
	fmt.Scanln()
}
