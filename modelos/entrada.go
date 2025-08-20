package modelos

import (
	"time"
)

// Entrada stores an entry
type Entrada struct {
	Id          uint         `json:"id" gorm:"unique"`
	IdUs        uint         `json:"id_us" gorm:"not null"`
	Usuario     string       `json:"usuario" gorm:"type:VARCHAR(50); not null"`
	Especial    *bool        `json:"especial" gorm:"not null"`
	Titulo      string       `json:"titulo" gorm:"type:VARCHAR(50); not null"`
	Fecha       string       `json:"fecha" gorm:"type:VARCHAR(16); not null"`
	Contenido   string       `json:"contenido" gorm:"not null"`
	Media       string       `json:"media" gorm:"sin-media"`
	Comentarios []Comentario `json:"comentarios" gorm:"foreignKey:IdEnt"`
	TotalCom    uint         `json:"total_com" gorm:"-"`
	Visitas     uint         `json:"visitas" gorm:"default:0"`
}

// CalcularTotalComentarios calculates the total sum of comments
func (entrada *Entrada) CalcularTotalComentarios() {
	entrada.TotalCom = uint(len(entrada.Comentarios))
}

// ValidarFecha validates an entry's date
func (entrada *Entrada) ValidarFecha() bool {
	if len(entrada.Fecha) != 16 || entrada.Fecha[10] != ' ' || entrada.Fecha[13] != ':' {
		return false
	}
	const layoutFechaNumerica string = "2006-01-02 15:04"
	_, err := time.Parse(layoutFechaNumerica, entrada.Fecha)
	return err == nil
}

// ValidarIdUs validates the user's id
func (entrada *Entrada) ValidarIdUs() bool {
	return entrada.IdUs != 0
}

// ValidarUsuario validates a user
func (entrada *Entrada) ValidarUsuario() bool {
	return entrada.Usuario != ""
}

// ValidarTitulo validates a title
func (entrada *Entrada) ValidarTitulo() bool {
	return entrada.Titulo != ""
}

// ValidarContenido validates a title
func (entrada *Entrada) ValidarContenido() bool {
	return entrada.Contenido != ""
}
