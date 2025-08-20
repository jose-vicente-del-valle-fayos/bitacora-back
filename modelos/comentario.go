package modelos

import (
	"time"
)

// Comentario stores a comment
type Comentario struct {
	Id         uint   `json:"id" gorm:"unique"`
	IdEnt      uint   `json:"id_ent" gorm:"not null"`
	Usuario    string `json:"usuario" gorm:"type:VARCHAR(50); not null"`
	Correo     string `json:"correo" gorm:"type:VARCHAR(50); not null"`
	Fecha      string `json:"fecha" gorm:"type:VARCHAR(16); not null"`
	Comentario string `json:"comentario" gorm:"not null"`
}

// ValidarFecha validates a date
func (comentario *Comentario) ValidarFecha() bool {
	if len(comentario.Fecha) != 16 || comentario.Fecha[10] != ' ' || comentario.Fecha[13] != ':' {
		return false
	}
	const layoutFechaNumerica string = "2006-01-02 15:04"
	_, err := time.Parse(layoutFechaNumerica, comentario.Fecha)
	return err == nil
}

// ValidarIdEnt validates an entry's id
func (comentario *Comentario) ValidarIdEnt() bool {
	return comentario.IdEnt != 0
}

// ValidarUsuario validates a user
func (comentario *Comentario) ValidarUsuario() bool {
	return true
}

// ValidarCorreo validates an email
func (comentario *Comentario) ValidarCorreo() bool {
	return comentario.Correo != ""
}

// ValidarComentario validates a comment
func (comentario *Comentario) ValidarComentario() bool {
	return comentario.Comentario != ""
}
