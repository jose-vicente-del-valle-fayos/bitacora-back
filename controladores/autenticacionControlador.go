package controladores

import (
	"github.com/gofiber/fiber/v2"
	"nd-back/bbdd"
	"nd-back/modelos"
	"nd-back/utilidades"
	"strconv"
	"time"
)

// Registrar registers a new user
//
//	{
//		"sobrenombre": "Chevi",
//		"nombre":      "José Vicente",
//		"apellidos":   "del Valle Fayos",
//		"correo":      "hola@chevi.soy",
//		"contrasena":  "almasera",
//	    "contrasenaconf": "almasera"
//	}
func Registrar(c *fiber.Ctx) error {
	var datos map[string]string
	if err := c.BodyParser(&datos); err != nil {
		return err
	}
	if datos["contrasena"] != datos["contrasenaconf"] {
		c.Status(400)
		return c.JSON(fiber.Map{"mensaje": "Las contraseñas no coinciden"})
	}
	var usuario modelos.Usuario
	r := bbdd.DB.Where("correo = ?", datos["correo"]).First(&usuario)
	if r.RowsAffected > 0 {
		usuario.PonContrasena(datos["contrasena"])
		bbdd.DB.Save(&usuario)
		return c.JSON(fiber.Map{"mensaje": "contraseña actualizada correctamente"})
	}
	usuario = modelos.Usuario{
		Sobrenombre: datos["sobrenombre"],
		Nombre:      datos["nombre"],
		Apellidos:   datos["apellidos"],
		Correo:      datos["correo"],
		Contrasena:  []byte(""),
	}
	usuario.PonContrasena(datos["contrasena"])
	bbdd.DB.Create(&usuario)
	return c.JSON(usuario)
}

// Ingresar allows to login a user
//
//	{
//		"correo":      "hola@chevi.soy",
//		"contrasena":  "1234"
//	}
func Ingresar(c *fiber.Ctx) error {
	var datos map[string]string
	if err := c.BodyParser(&datos); err != nil {
		return err
	}
	usuario := modelos.Usuario{}
	bbdd.DB.Where("correo = ?", datos["correo"]).First(&usuario)
	if usuario.Id == 0 {
		time.Sleep(3 * time.Second)
		c.Status(404)
		return c.JSON(fiber.Map{"mensaje": "usuario no encontrado"})
	}
	if err := usuario.ComparaContrasenas(datos["contrasena"]); err != nil {
		time.Sleep(3 * time.Second)
		c.Status(400)
		return c.JSON(fiber.Map{"mensaje": "el nombre de usuario o la contraseña no son correctos"})
	}
	token, err := utilidades.GenerarJWT(strconv.Itoa(int(usuario.Id)))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	cookie := fiber.Cookie{
		Name:     "btcr-jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 12), // Expira en medio día
		HTTPOnly: true,
		Secure:   true, // Solo se enviará si la conexión es segura (HTTPS)
		SameSite: "Lax",
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"mensaje": "autenticado",
	})
}

// Usuario checks is a user is logged
func Usuario(c *fiber.Ctx) error {
	cookie := c.Cookies("btcr-jwt")
	id, err := utilidades.ParsearJWT(cookie)
	if err != nil {
		// c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"mensaje": "sin autenticar",
		})
	}
	usuario := modelos.Usuario{}
	bbdd.DB.Where("id = ?", id).First(&usuario)
	return c.JSON(usuario)
}

// Salir makes the JWT expired
func Salir(c *fiber.Ctx) error {
	// Eliminar completamente la cookie "btcr-jwt"
	cookie := fiber.Cookie{
		Name:     "btcr-jwt",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"mensaje": "salido",
	})
}
