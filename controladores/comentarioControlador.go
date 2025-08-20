package controladores

import (
	"github.com/gofiber/fiber/v2"
	"nd-back/bbdd"
	"nd-back/modelos"
	"strconv"
)

// TodosComentarios returns all comments
func TodosComentarios(c *fiber.Ctx) error {
	var comentarios []modelos.Comentario
	bbdd.DB.Find(&comentarios)
	var total int64
	bbdd.DB.Model(&comentarios).Count(&total)
	return c.JSON(fiber.Map{
		"datos": comentarios,
		"meta": fiber.Map{
			"total": total,
		},
	})
}

// CrearComentario creates a comment
//
//	{
//		"id_ent": 	1,
//		"usuario":    "Sara",
//		"correo":     "sara@gmail.com",
//		"fecha":      "2018-04-02",
//		"comentario": "Hola Chevi. Eres el mejor."
//	}
func CrearComentario(c *fiber.Ctx) error {
	var comentario modelos.Comentario
	if err := c.BodyParser(&comentario); err != nil {
		return err
	}
	if comentario.ValidarFecha() && comentario.ValidarIdEnt() && comentario.ValidarUsuario() && comentario.ValidarCorreo() && comentario.ValidarComentario() {
		bbdd.DB.Create(&comentario)
		return c.JSON(comentario)
	}
	return c.JSON(fiber.Map{"mensaje": "error de validación"})
}

// LeerComentario reads a comment taking the comment's id as a URL parameter
func LeerComentario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	comentario := modelos.Comentario{
		Id: uint(id),
	}
	bbdd.DB.Find(&comentario)
	return c.JSON(comentario)
}

// ActualizarComentario updates a comment
//
//	{
//		"id_ent": 	1,
//		"usuario":    "Sara actualizada",
//		"correo":     "saractualizada@gmail.com",
//		"fecha":      "2018-04-03",
//		"comentario": "Hola Chevi. Eres el mejor del mundo mundial."
//	}
func ActualizarComentario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	comentario := modelos.Comentario{
		Id: uint(id),
	}
	if err := c.BodyParser(&comentario); err != nil {
		return err
	}
	if comentario.ValidarFecha() && comentario.ValidarIdEnt() && comentario.ValidarUsuario() && comentario.ValidarCorreo() && comentario.ValidarComentario() {
		bbdd.DB.Model(&comentario).Updates(comentario)
		return c.JSON(comentario)
	}
	return c.JSON(fiber.Map{"mensaje": "error de validación"})
}

// BorrarComentario deletes a comment taking the comment's id as a parameter
func BorrarComentario(c *fiber.Ctx) error {
	idUs, err := strconv.Atoi(c.Params("id_us"))
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	comentario := modelos.Comentario{
		Id: uint(id),
	}
	bbdd.DB.Find(&comentario)
	var entrada modelos.Entrada
	bbdd.DB.Find(&entrada, comentario.IdEnt)
	if entrada.IdUs == uint(idUs) {
		comentario := modelos.Comentario{
			Id: uint(id),
		}
		bbdd.DB.Delete(&comentario)
	}
	return nil
}
