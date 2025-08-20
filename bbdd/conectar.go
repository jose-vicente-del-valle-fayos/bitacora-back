package bbdd

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nd-back/modelos"
	"os"
)

var DB *gorm.DB

// Conectar connects the API to the database and initializes Usuario, Page and Comentario if they aren't initialized
func Conectar() {
	// db, err := gorm.Open(postgres.Open("postgresql://postgres:almasera@localhost:5432/nd"), &gorm.Config{})
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("no se puede conectar a la base de datos.")
	}
	err = db.AutoMigrate(&modelos.Usuario{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&modelos.Entrada{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&modelos.Comentario{})
	if err != nil {
		fmt.Println(err)
	}
	DB = db
}
