package controladores

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/idna"
	"gorm.io/gorm"
	"io"
	"log"
	"nd-back/bbdd"
	"nd-back/modelos"
	"os"
	"path/filepath"
	"strconv"
)

// TodasEntradas returns an array of entries inside datos
func TodasEntradas(c *fiber.Ctx) error {
	limite, err := strconv.Atoi(c.Query("limite", strconv.Itoa(5)))
	if err != nil {
		return err
	}
	especial, err := strconv.ParseBool(c.Query("especial", "false"))
	if err != nil {
		return err
	}
	ultimaFecha := c.Query("ultima")
	var entradas []modelos.Entrada
	query := bbdd.DB.Preload("Comentarios", func(db *gorm.DB) *gorm.DB {
		return db.Order("fecha ASC")
	}).Where("especial = ?", especial)
	if ultimaFecha != "" {
		query = query.Where("fecha < ?", ultimaFecha)
	}
	query.Order("id desc").Limit(limite).Find(&entradas)
	for i := range entradas {
		entradas[i].CalcularTotalComentarios()
	}
	var total int64
	cuantos := bbdd.DB.Model(&modelos.Entrada{}).Where("especial = ?", especial)
	if ultimaFecha != "" {
		cuantos = cuantos.Where("fecha < ?", ultimaFecha)
	}
	cuantos.Count(&total)
	quedan := total > int64(limite)
	return c.JSON(fiber.Map{
		"datos":  entradas,
		"quedan": quedan,
	})
}

// ExtractoTodas returns an array with entries extract data inside datos
func ExtractoTodas(c *fiber.Ctx) error {
	var entradas []modelos.Entrada
	bbdd.DB.Select("Id", "Titulo", "Fecha", "Contenido").Order("fecha desc").Find(&entradas)
	/*
		if len(entradas) > 0 {
			// Actualizar las visitas de todas las entradas a 0
			for i := range entradas {
				entradas[i].Visitas = 0
			}
			// Actualiza todas las entradas a la vez con las visitas a 0
			bbdd.DB.Model(&entradas).Updates(map[string]interface{}{"Visitas": 0})
		}
	*/
	return c.JSON(fiber.Map{
		"datos": entradas,
	})
}

// CrearEntrada creates an entry
//
//	{
//		"id_us": 1,
//		"usuario": "Chevi",
//		"especial": false,
//		"titulo": "Esta es una entrada fantástica",
//		"fecha": "2024-02-22",
//		"contenido": "Este es un contenido fantástico."
//	}
func CrearEntrada(c *fiber.Ctx) error {
	idUs, _ := strconv.ParseUint(c.FormValue("id_us"), 10, 32)
	usuario := c.FormValue("usuario")
	especial, _ := strconv.ParseBool(c.FormValue("especial"))
	titulo := c.FormValue("titulo")
	fecha := c.FormValue("fecha")
	contenido := c.FormValue("contenido")
	entrada := modelos.Entrada{
		IdUs:      uint(idUs),
		Usuario:   usuario,
		Especial:  &especial,
		Titulo:    titulo,
		Fecha:     fecha,
		Contenido: contenido,
	}
	// if err := c.BodyParser(&entrada); err != nil {
	// 	return err
	// }
	if entrada.ValidarFecha() && entrada.ValidarUsuario() && entrada.ValidarTitulo() && entrada.ValidarContenido() {
		bbdd.DB.Create(&entrada)
		fmt.Println("entrada: " + strconv.Itoa(int(entrada.Id)))
		entrada.Media = SubirMedia(c, entrada.Id)
		bbdd.DB.Model(&entrada).Where("id = ?", entrada.Id).Updates(entrada)
		return c.JSON(entrada)
	}
	return c.JSON(fiber.Map{"mensaje": "error de validación"})
}

// LeerEntrada reads an entry taking entry's id as a URL parameter
func LeerEntrada(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	var entrada modelos.Entrada
	bbdd.DB.Preload("Comentarios", func(db *gorm.DB) *gorm.DB {
		return db.Order("fecha ASC")
	}).Find(&entrada, id)
	entrada.CalcularTotalComentarios()
	return c.JSON(fiber.Map{
		"datos": entrada,
	})
}

// ActualizarEntrada updates an entry
//
//	{
//		"id_us": 1,
//		"usuario": "Susanita",
//		"especial": false,
//		"titulo": "Esta es una entrada actualizada",
//		"fecha": "2024-02-22",
//		"contenido": "Este es un contenido actualizado.",
//		"comentarios": []
//	}
func ActualizarEntrada(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.FormValue("id"), 10, 32)
	BorrarMedia(uint(id))
	idUs, _ := strconv.ParseUint(c.FormValue("id_us"), 10, 32)
	usuario := c.FormValue("usuario")
	especial, _ := strconv.ParseBool(c.FormValue("especial"))
	titulo := c.FormValue("titulo")
	fecha := c.FormValue("fecha")
	contenido := c.FormValue("contenido")
	entrada := modelos.Entrada{
		Id:        uint(id),
		IdUs:      uint(idUs),
		Usuario:   usuario,
		Especial:  &especial,
		Titulo:    titulo,
		Fecha:     fecha,
		Contenido: contenido,
		Media:     SubirMedia(c, uint(id)),
	}
	if entrada.ValidarFecha() && entrada.ValidarUsuario() && entrada.ValidarTitulo() && entrada.ValidarContenido() {
		bbdd.DB.Model(&entrada).Where("id = ?", id).Updates(entrada)
		return c.JSON(entrada)
	}
	return c.JSON(fiber.Map{"mensaje": "error de validación"})
}

// BorrarEntrada deletes an entry taking the entry's id as a URL parameter
func BorrarEntrada(c *fiber.Ctx) error {
	idUs, err := strconv.Atoi(c.Params("id_us"))
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	var entrada modelos.Entrada
	bbdd.DB.Find(&entrada, id)
	if entrada.IdUs == uint(idUs) {
		BorrarMedia(uint(id))
		bbdd.DB.Delete(&entrada)
	}
	return nil
}

// SubirMedia uploads a file (jpg, jpeg or mp3)
func SubirMedia(c *fiber.Ctx, id uint) string {
	idStr := strconv.FormatUint(uint64(id), 10)
	fileHeader, err := c.FormFile("media-entrada")

	if err != nil {
		fmt.Println("no se pudo obtener el archivo: ", err)
		return "sin-media"
	}

	if fileHeader.Size == 0 {
		fmt.Println("el archivo está vacío")
		return "sin-media"
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("error al abrir el archivo: ", err)
		return "sin-media"
	}
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".dat"
	}
	defer file.Close()

	localDir := "/var/data/media"
	if err := os.MkdirAll(localDir, 0755); err != nil {
		fmt.Println("error creando directorio local: ", err)
		return "sin-media"
	}

	localPath := filepath.Join(localDir, idStr+ext)
	localFile, err := os.Create(localPath)
	if err != nil {
		fmt.Println("error creando archivo local: ", err)
		return "sin-media"
	}

	if _, err := io.Copy(localFile, file); err != nil {
		fmt.Println("error guardando archivo local: ", err)
		localFile.Close()
		return "sin-media"
	}

	localFile.Close()

	diskUrl, err := idna.ToASCII(os.Getenv("PERSISTENT_DISK_URL"))
	if err != nil {
		log.Println("error al convertir la url del disco a punycode: ", err)
	}
	publicURL := fmt.Sprintf("%s/%s%s", diskUrl, idStr, ext)

	return publicURL
}

// BorrarMedia destroys a file (jpg, jpeg or mp3)
func BorrarMedia(id uint) bool {
	idStr := strconv.FormatUint(uint64(id), 10)
	localDir := "/var/data/media"

	pattern := filepath.Join(localDir, idStr+".*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println("error buscando archivo: ", err)
		return false
	}

	if len(files) == 0 {
		fmt.Println("no se encontró archivo para id:", idStr)
		return false
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("error al borrar archivo:", err)
			return false
		}
		fmt.Println("archivo eliminado:", file)
	}

	return true
}
