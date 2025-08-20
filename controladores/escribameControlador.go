package controladores

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
	"nd-back/modelos"
	"os"
	"strconv"
	"time"
)

var bloquearEnvio = false
var tiempoLlamada = time.Now().Add(-1 * GetEnvDuracion("CORREO_TIMEOUT_TRAMO_5"))
var numLlamadas = map[int]int{
	1: 0,
	2: 0,
	3: 0,
	4: 0,
	5: 0,
}

func GetEnvMaxLlam(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func GetEnvDuracion(key string) time.Duration {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		fmt.Println(err)
	}
	return time.Duration(val) * time.Second
}

func ComprobarBloqueo(tramo int, maxLlamadas int, intervalo time.Duration) bool {
	if time.Now().Sub(tiempoLlamada) > intervalo {
		numLlamadas[tramo] = 0
	}
	return numLlamadas[tramo] > maxLlamadas
}

// Escribeme sends an email to the ND owner
//
//	{
//		"nombre": "Chevi",
//		"correo": "chevielmejor@gmail.com",
//		"mensaje": "Hola chevi. Qué tal estás?"
//	}
func Escribame(c *fiber.Ctx) error {
	tramos := []struct {
		maxLlamadas int
		timeout     time.Duration
	}{
		{GetEnvMaxLlam("CORREO_MAX_LLAMADAS_TRAMO_1"), GetEnvDuracion("CORREO_TIMEOUT_TRAMO_1")},
		{GetEnvMaxLlam("CORREO_MAX_LLAMADAS_TRAMO_2"), GetEnvDuracion("CORREO_TIMEOUT_TRAMO_2")},
		{GetEnvMaxLlam("CORREO_MAX_LLAMADAS_TRAMO_3"), GetEnvDuracion("CORREO_TIMEOUT_TRAMO_3")},
		{GetEnvMaxLlam("CORREO_MAX_LLAMADAS_TRAMO_4"), GetEnvDuracion("CORREO_TIMEOUT_TRAMO_4")},
		{GetEnvMaxLlam("CORREO_MAX_LLAMADAS_TRAMO_5"), GetEnvDuracion("CORREO_TIMEOUT_TRAMO_5")},
	}
	bloquearEnvio = false
	for i, tramo := range tramos {
		if ComprobarBloqueo(i+1, tramo.maxLlamadas, tramo.timeout) {
			bloquearEnvio = true
		}
	}

	var correo modelos.Correo
	if err := c.BodyParser(&correo); err != nil {
		return err
	}
	if correo.ValidarNombre() && correo.ValidarCorreo() && correo.ValidarMensaje() && !bloquearEnvio {
		m := gomail.NewMessage()
		m.SetHeader("From", m.FormatAddress(os.Getenv("CORREO_FROM"), "chevi.soy") /* email */)
		m.SetHeader("To", os.Getenv("CORREO_TO"))
		m.SetAddressHeader("reply-to", correo.Correo, correo.Nombre)
		m.SetHeader("Subject", "Correo desde chevi.soy")
		m.SetBody("text/html", correo.Mensaje)
		port, _ := strconv.Atoi(os.Getenv("CORREO_PORT"))
		d := gomail.NewDialer(os.Getenv("CORREO_SERVER"), port, os.Getenv("CORREO_FROM"), os.Getenv("CORREO_PASS"))
		err := d.DialAndSend(m)
		if err != nil {
			return err
		} else {
			tiempoLlamada = time.Now()
			for t, v := range numLlamadas {
				numLlamadas[t] = v + 1
			}
			return c.JSON(correo)
		}
	} else {
		return errors.New("error de validación")
	}
}
