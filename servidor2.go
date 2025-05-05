package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MODELOS

type Estudiante struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"unique"`
}

type Profesor struct {
	ID     uint   `gorm:"primaryKey"`
	Nombre string `gorm:"unique" json:"nombre"`
}

func (Profesor) TableName() string {
	return "profesores"
}

type Asignatura struct {
	ID         uint `gorm:"primaryKey"`
	Nombre     string
	ProfesorID uint
}

type Matricula struct {
	ID           uint `gorm:"primaryKey"`
	EstudianteID uint
	AsignaturaID uint
	Ciclo        string
	Nota1        float64 `gorm:"default:0"`
	Nota2        float64 `gorm:"default:0"`
	Supletorio   float64 `gorm:"default:0"`
}

// GLOBAL
var db *gorm.DB

func main() {
	dsn := "root:Scarleth3023@tcp(127.0.0.1:3306)/servidor2?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar con base de datos: " + err.Error())
	}

	// Migrar modelos
	db.AutoMigrate(&Estudiante{}, &Profesor{}, &Asignatura{}, &Matricula{})

	r := gin.Default()

	// ✅ Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 🔥 NUEVO ENDPOINT /ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// RUTAS
	r.POST("/profesores", insertarProfesor)
	r.GET("/profesores", listarProfesores)
	r.DELETE("/profesores/:id", eliminarProfesor)

	r.GET("/asignaturas", listarAsignaturas)
	r.POST("/asignaturas", insertarAsignatura)
	r.DELETE("/asignaturas/:id", eliminarAsignatura)

	r.GET("/estudiantes", obtenerEstudiantes)
	r.POST("/estudiantes", insertarEstudiante)
	r.DELETE("/estudiantes/:id", eliminarEstudiante)

	r.GET("/matriculas", listarMatriculas)
	r.PUT("/matricula/:id", actualizarNotas)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run("0.0.0.0:" + port)

}

// FUNCIONES

func listarProfesores(c *gin.Context) {
	var profesores []Profesor
	db.Find(&profesores)
	c.JSON(http.StatusOK, profesores)
}

func insertarProfesor(c *gin.Context) {
	var nuevoProfesor Profesor
	if err := c.ShouldBindJSON(&nuevoProfesor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos no válidos"})
		return
	}
	if err := db.Create(&nuevoProfesor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar profesor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Profesor registrado"})
}

func eliminarProfesor(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Profesor{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el profesor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Profesor eliminado"})
}

func listarAsignaturas(c *gin.Context) {
	var asignaturas []Asignatura
	db.Find(&asignaturas)
	c.JSON(http.StatusOK, asignaturas)
}

func insertarAsignatura(c *gin.Context) {
	var nueva Asignatura
	if err := c.ShouldBindJSON(&nueva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if err := db.Create(&nueva).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear asignatura"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Asignatura creada"})
}

func eliminarAsignatura(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Asignatura{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la asignatura"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Asignatura eliminada"})
}

func obtenerEstudiantes(c *gin.Context) {
	var estudiantes []Estudiante
	db.Find(&estudiantes)
	c.JSON(http.StatusOK, estudiantes)
}

func insertarEstudiante(c *gin.Context) {
	var nuevo Estudiante
	if err := c.ShouldBindJSON(&nuevo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if err := db.Create(&nuevo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo insertar estudiante"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Estudiante registrado"})
}

func eliminarEstudiante(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Estudiante{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar estudiante"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Estudiante eliminado"})
}

func listarMatriculas(c *gin.Context) {
	var matriculas []Matricula
	db.Find(&matriculas)
	c.JSON(http.StatusOK, matriculas)
}

func actualizarNotas(c *gin.Context) {
	var notas struct {
		Nota1      float64 `json:"nota1"`
		Nota2      float64 `json:"nota2"`
		Supletorio float64 `json:"supletorio"`
	}
	id := c.Param("id")

	if err := c.ShouldBindJSON(&notas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var m Matricula
	if err := db.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Matrícula no encontrada"})
		return
	}

	m.Nota1 = notas.Nota1
	m.Nota2 = notas.Nota2
	m.Supletorio = notas.Supletorio
	db.Save(&m)

	c.JSON(http.StatusOK, gin.H{"status": "Notas actualizadas"})
}
