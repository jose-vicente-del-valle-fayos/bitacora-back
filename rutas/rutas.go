package rutas

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"nd-back/controladores"
	"nd-back/middlewares"
	"os"
	"strconv"
)

// Configuracion selects a route and calls to a controller
func Configuracion(app *fiber.App) {
	reg, err := strconv.ParseBool(os.Getenv("REGISTRAR_ENABLED"))
	if err != nil {
		fmt.Println(err)
	}
	if reg {
		app.Post("/registrar", controladores.Registrar)
	}
	app.Post("/ingresar", controladores.Ingresar)
	app.Get("/entrada/:id", controladores.LeerEntrada)
	app.Get("/entradas", controladores.TodasEntradas)
	app.Get("/todas", controladores.ExtractoTodas)
	app.Get("/especial", controladores.TodasEntradas)
	app.Get("/comentarios", controladores.TodosComentarios)
	app.Post("/escribame", controladores.Escribame)
	app.Post("/visita/:id", controladores.RegistrarVisita)
	app.Use(middlewares.Autenticado)
	app.Get("/usuario", controladores.Usuario)
	app.Post("/salir", controladores.Salir)
	app.Post("/usuario", controladores.CrearUsuario)
	app.Post("/entrada", controladores.CrearEntrada)
	app.Post("/comentario", controladores.CrearComentario)
	app.Get("/:id", controladores.LeerUsuario)
	app.Get("/comentario/:id", controladores.LeerComentario)
	app.Put("/usuario/:id", controladores.ActualizarUsuario)
	app.Put("/entrada/:id", controladores.ActualizarEntrada)
	app.Put("/comentario/:id", controladores.ActualizarComentario)
	// app.Delete("/usuario/:id", controladores.BorrarUsuario)
	app.Delete("/entrada/:id_us/:id", controladores.BorrarEntrada)
	app.Delete("/comentario/:id_us/:id", controladores.BorrarComentario)
	// app.Get("/usuarios", controladores.TodosUsuarios)
}
