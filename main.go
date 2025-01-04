package main

import (
	"context"  // Для управления контекстом выполнения.
	"log"      // Для логирования ошибок и сообщений.
	"net/http" // Для работы с HTTP-сервером.

	"github.com/gin-gonic/gin" // Фреймворк для создания веб-приложений.
	"github.com/jackc/pgx/v5"  // Библиотека для работы с PostgreSQL.
)

type Aviation struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Plane string  `json:"plane"`
	Price float64 `json:"price"`
}

var conn *pgx.Conn // Соединение с базой данных.

func main() {
	// Настройка подключения к PostgreSQL.
	connStr := "postgres://postgres:postgres@localhost:5432/postgres"
	var err error

	// Устанавливаем соединение с базой данных.
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	defer conn.Close(context.Background()) // Закрываем соединение при завершении работы.

	// Создание маршрутов.
	router := gin.Default()
	router.GET("/aviation", getAllAviation)
	router.GET("/aviation/:id", getAviationByID)
	router.POST("/aviation", postAviation)
	router.PUT("/aviation/:id", putAviation) // Добавляем PUT маршрут

	// Запуск HTTP-сервера.
	router.Run("localhost:8080")
}

// getAllAviation возвращает список всех записей из таблицы aviation.
func getAllAviation(c *gin.Context) {
	query := "SELECT id, title, plane, price FROM aviation"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения данных"})
		return
	}
	defer rows.Close()

	var aviationData []Aviation
	for rows.Next() {
		var aviation Aviation
		if err := rows.Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обработки данных"})
			return
		}
		aviationData = append(aviationData, aviation)
	}

	c.IndentedJSON(http.StatusOK, aviationData)
}

// getAviationByID возвращает запись по ID.
func getAviationByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, title, plane, price FROM aviation WHERE id = $1"

	var aviation Aviation
	err := conn.QueryRow(context.Background(), query, id).Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Запись не найдена"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения записи"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, aviation)
}

// postAviation добавляет новую запись.
func postAviation(c *gin.Context) {
	var newAviation Aviation

	if err := c.BindJSON(&newAviation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверные входные данные"})
		return
	}

	if newAviation.Price <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Некорректная цена"})
		return
	}

	query := "INSERT INTO aviation (id, title, plane, price) VALUES ($1, $2, $3, $4)"
	_, err := conn.Exec(context.Background(), query, newAviation.ID, newAviation.Title, newAviation.Plane, newAviation.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления записи"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAviation)
}

// putAviation обновляет существующую запись.
func putAviation(c *gin.Context) {
	id := c.Param("id")
	var updatedAviation Aviation

	// Считываем данные для обновления из тела запроса.
	if err := c.BindJSON(&updatedAviation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверные входные данные"})
		return
	}

	// Проверяем корректность цены.
	if updatedAviation.Price <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Некорректная цена"})
		return
	}

	// Выполняем обновление записи.
	query := "UPDATE aviation SET title = $1, plane = $2, price = $3 WHERE id = $4"
	cmdTag, err := conn.Exec(context.Background(), query, updatedAviation.Title, updatedAviation.Plane, updatedAviation.Price, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления записи"})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Запись не найдена"})
		return
	}

	// Возвращаем обновленную запись.
	c.IndentedJSON(http.StatusOK, updatedAviation)
}
