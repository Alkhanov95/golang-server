package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type aviation struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Plane string  `json:"plane"`
	Price float64 `json:"price"`
}

// todo: убрать в БД (postgres sql) поставить postgres
// добавть таблицу aviation c полями структуры type aviation struct
// var aviationData = []aviation{
// 	{ID: "1", Title: "737 MAX", Plane: "Boeing", Price: 90000000.00},
// 	{ID: "2", Title: "A380", Plane: "Airbus", Price: 445000000.00},
// 	{ID: "3", Title: "A320", Plane: "Airbus", Price: 120000000.00},
// }

// todo:
// var (
// 	conn pgx.Conn
// )

func main() {

	// get db conn

	router := gin.Default()
	router.GET("/aviation", getAviation)
	router.GET("/aviation/:id", getAviationByID)
	router.POST("/aviation", postAviation) // подключение маршрута POST
	router.Run("localhost:8080")
}

// func getAviationByID(...) (*aviation, error) {
// 	// поход в базу за данными
// 	// pgx.exec...
// 	// return ...
// }

// func listAviation(...) (*aviation, error) {
// 	// поход в базу за данными
// 	// pgx.exec...
// 	// return ...
// }

// getAviation возвращает весь список самолётов.
func getAviation(c *gin.Context) {

	// todo: доставать из базы данные и отдавать их в ответе
	// aviationData, err := getAviationFromDB(...)

	c.IndentedJSON(http.StatusOK, aviationData)
}

// getAviationByID возвращает информацию о самолёте по его ID.
func getAviationByID(c *gin.Context) {
	id := c.Param("id")

	// todo: доставать из базы данные и отдавать их в ответе
	// aviationData, err := listAviation(...)
	// c.IndentedJSON(http.StatusOK, aviationData) // Если нашли, возвращаем его данные

	// Если самолёт не найден, возвращаем ошибку
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "aviation not found"})
}

// postAviation добавляет новый самолёт на основе полученных JSON-данных.
func postAviation(c *gin.Context) {
	var newAviation aviation

	// добавить валидацию по полю (чтобы если данные неправильные была чёткая ошибка (invalid price))

	// Привязываем полученные JSON-данные к newAviation.
	if err := c.BindJSON(&newAviation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	// todo: записать данные в БД
	// err := addAviationToDB(newAviation)

	// Добавляем новый самолёт в список
	// aviationData = append(aviationData, newAviation)
	c.IndentedJSON(http.StatusCreated, newAviation)
}
