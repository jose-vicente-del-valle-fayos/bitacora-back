package modelos

import (
	"golang.org/x/crypto/bcrypt"
)

// Usuario stores a user
type Usuario struct {
	Id          uint      `json:"id" gorm:"unique"`
	Sobrenombre string    `json:"sobrenombre" gorm:"type:VARCHAR(50); unique; not null"`
	Nombre      string    `json:"nombre" gorm:"type:VARCHAR(50); not null"`
	Apellidos   string    `json:"apellidos" gorm:"type:VARCHAR(50); not null"`
	Correo      string    `json:"correo" gorm:"type:VARCHAR(50); unique; not null"`
	Contrasena  []byte    `json:"-" gorm:"not null"`
	Entradas    []Entrada `json:"entradas" gorm:"foreignKey:IdUs"`
	TotalEnt    uint      `json:"total_ent" gorm:"-"`
}

// PonContrasena converts the password into a hash
func (usuario *Usuario) PonContrasena(contrasena string) {
	hashCont, _ := bcrypt.GenerateFromPassword([]byte(contrasena), 14)
	usuario.Contrasena = hashCont
}

// ComparaContrasenas compares two passwords
func (usuario *Usuario) ComparaContrasenas(contrasena string) error {
	return bcrypt.CompareHashAndPassword(usuario.Contrasena, []byte(contrasena))
}

// CalcularTotalEntradas calculates the total sum of entradas
func (usuario *Usuario) CalcularTotalEntradas() {
	usuario.TotalEnt = uint(len(usuario.Entradas))
}
