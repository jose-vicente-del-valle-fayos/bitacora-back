package controladores

import (
	"github.com/gofiber/fiber/v2"
	"nd-back/bbdd"
	"nd-back/modelos"
	"strconv"
)

// TodosUsuarios returs all users
func TodosUsuarios(c *fiber.Ctx) error {
	var usuarios []modelos.Usuario
	bbdd.DB.Preload("Entradas").Preload("Entradas.Comentarios").Find(&usuarios)
	return c.JSON(usuarios)
}

// CrearUsuario creates a user
//
//	{
//		"sobrenombre": "Chevi",
//		"nombre":      "José Vicente",
//		"apellidos":   "del Valle Fayos",
//		"correo":      "hola@chevi.soy"
//	}
func CrearUsuario(c *fiber.Ctx) error {
	var usuario modelos.Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return err
	}
	usuario.PonContrasena("1234")
	bbdd.DB.Create(&usuario)
	return c.JSON(usuario)
}

// LeerUsuario reads a user
func LeerUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	usuario := modelos.Usuario{
		Id: uint(id),
	}
	bbdd.DB.Preload("Entradas").Find(&usuario)
	usuario.CalcularTotalEntradas()
	return c.JSON(usuario)
}

// ActualizarUsuario updates a user
//
//	{
//		"sobrenombre": "Ramón",
//		"nombre":      "Supermones rascajeribrother",
//		"apellidos":   "Bronson",
//		"correo":      "pato@gmail.com"
//	}
func ActualizarUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	usuario := modelos.Usuario{
		Id: uint(id),
	}
	if err := c.BodyParser(&usuario); err != nil {
		return err
	}
	bbdd.DB.Model(&usuario).Updates(usuario)
	return c.JSON(usuario)
}

/*
// BorrarUsuario deletes a user
func BorrarUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	usuario := modelos.Usuario{
		Id: uint(id),
	}
	bbdd.DB.Delete(&usuario)
	return nil
}
*/
