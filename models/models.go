package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cliente struct {
	Documento string `gorm:"primaryKey"`
	Nombre    string
	Apellido  string
	Cuentas   []*Cuenta `gorm:"foreignKey:DocumentoCliente"`
}

type Cuenta struct {
	Numero           string `gorm:"primaryKey;type:uuid"`
	DocumentoCliente string
	Tipo             string
	Saldo            float64
	Cliente          *Cliente       `gorm:"foreignKey:DocumentoCliente;references:Documento"`
	Transacciones    []*Transaccion `gorm:"foreignKey:NumeroCuenta"`
}

func (c *Cuenta) BeforeCreate(tx *gorm.DB) (err error) {
	c.Numero = uuid.New().String()
	return
}

type TipoTransaccion struct {
	ID            uint8          `gorm:"primaryKey"`
	Nombre        string         `gorm:"unique"`
	Transacciones []*Transaccion `gorm:"foreignKey:TipoID"`
}

type Transaccion struct {
	ID           uint `gorm:"primaryKey"`
	Fecha        time.Time
	NumeroCuenta string
	TipoID       uint8
	Monto        float64
	Descripcion  string
	Tipo         *TipoTransaccion `gorm:"foreignKey:TipoID;references:ID"`
	Cuenta       *Cuenta          `gorm:"foreignKey:NumeroCuenta;references:Numero"`
}

func AutoMigraryLlenar(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Cliente{},
		&Cuenta{},
		&TipoTransaccion{},
		&Transaccion{},
	)
	if err != nil {
		return err
	}

	tipos := []TipoTransaccion{
		{ID: 1, Nombre: "deposito"},
		{ID: 2, Nombre: "retiro"},
		{ID: 3, Nombre: "transferencia"},
	}

	for _, t := range tipos {
		db.FirstOrCreate(&t, t)
	}

	return nil

}
