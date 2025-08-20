package controladores

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"nd-back/bbdd"
	"nd-back/modelos"
	"strconv"
)

func RegistrarVisita(c *fiber.Ctx) error {
	idTemp, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	id := uint(idTemp)
	var entrada modelos.Entrada
	if err := bbdd.DB.First(&entrada, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "entrada no encontrada"})
	}
	bbdd.DB.Model(&entrada).Update("visitas", gorm.Expr("visitas + 1"))
	entrada.CalcularTotalComentarios()
	return c.JSON(fiber.Map{"success": true, "id": id})
}
