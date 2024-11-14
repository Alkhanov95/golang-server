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

var aviationData = []aviation{
	{ID: "1", Title: "737 MAX", Plane: "Boeing", Price: 90000000.00},
	{ID: "2", Title: "A380", Plane: "Airbus", Price: 445000000.00},
	{ID: "3", Title: "A320", Plane: "Airbus", Price: 120000000.00},
}

func main() {
	router := gin.Default()
	router.GET("/aviation", getAviation)
	router.GET("/aviation/:id", getAviationByID)
	router.POST("/aviation", postAviation) // подключение маршрута POST
	router.Run("localhost:8080")
}

// getAviation возвращает весь список самолётов.
func getAviation(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, aviationData)
}

// getAviationByID возвращает информацию о самолёте по его ID.
func getAviationByID(c *gin.Context) {
	id := c.Param("id")

	// Ищем самолёт по ID в списке
	for _, a := range aviationData {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a) // Если нашли, возвращаем его данные
			return
		}
	}

	// Если самолёт не найден, возвращаем ошибку
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "aviation not found"})
}

// postAviation добавляет новый самолёт на основе полученных JSON-данных.
func postAviation(c *gin.Context) {
	var newAviation aviation

	// Привязываем полученные JSON-данные к newAviation.
	if err := c.BindJSON(&newAviation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	// Добавляем новый самолёт в список
	aviationData = append(aviationData, newAviation)
	c.IndentedJSON(http.StatusCreated, newAviation)
}
