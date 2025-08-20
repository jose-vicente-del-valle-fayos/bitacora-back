package modelos

// Correo stores a mail
type Correo struct {
	Nombre  string `json:"nombre"`
	Correo  string `json:"correo"`
	Mensaje string `json:"mensaje"`
}

// ValidarNombre validates a name
func (correo *Correo) ValidarNombre() bool { return correo.Nombre != "" }

// ValidarCorreo validates an email address
func (correo *Correo) ValidarCorreo() bool { return correo.Correo != "" }

// ValidarMensaje validates a message
func (correo *Correo) ValidarMensaje() bool { return correo.Mensaje != "" }
