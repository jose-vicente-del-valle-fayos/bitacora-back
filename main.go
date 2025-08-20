package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/net/idna"
	"log"
	"nd-back/bbdd"
	"nd-back/rutas"
	"os"
)

func main() {
	bbdd.Conectar()
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		hostnamePermitido, err := idna.ToASCII(os.Getenv("HOSTNAME_PERMITIDO"))
		if err != nil {
			log.Println("error al convertir el hostname permitido a punycode: ", err)
		}
		hostnameSolicitado, err := idna.ToASCII(c.Hostname())
		if err != nil {
			log.Println("error al convertir el hostname solicitado a punycode: ", err)
		}
		if hostnamePermitido != hostnameSolicitado {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.Next()
	})
	corsDominioPermitido, err := idna.ToASCII(os.Getenv("CORS_DOMINIO_PERMITIDO"))
	if err != nil {
		log.Println("error al convertir el dominio permitido de cors a punycode: ", err)
	}
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization",
		AllowOrigins:     corsDominioPermitido,
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE",
		MaxAge:           86400,
	}))
	app.Static("/media", "/var/data/media")
	rutas.Configuracion(app)
	err = app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
