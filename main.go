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

var conn *pgx.Conn

func main() {
	// Строка подключения
	connStr := "postgresql://postgres:@postgres:5432/postgres" // Используется имя контейнера "db"
	var err error

	// Устанавливаем соединение с базой данных.
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	defer func() {
		// Закрытие соединения при завершении работы.
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Ошибка при закрытии соединения с базой данных: %v\n", err)
		}
	}()

	// Создание маршрутов.
	router := gin.Default()
	router.GET("/aviation", getAllAviation)
	router.GET("/aviation/:id", getAviationByID)
	router.POST("/aviation", postAviation)
	router.PUT("/aviation/:id", putAviation) // Добавляем PUT маршрут
	router.DELETE("/aviation/:id", deleteAviationByID)

	// Запуск HTTP-сервера.
	log.Println("Запуск сервера на порту 8080...")
	router.Run("0.0.0.0:8080")
}

// getAllAviation возвращает список всех записей из таблицы aviation.
func getAllAviation(c *gin.Context) {
	query := "SELECT id, title, plane, price FROM aviation"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Ошибка при получении данных: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения данных"})
		return
	}
	defer rows.Close()

	var aviationData []Aviation
	for rows.Next() {
		var aviation Aviation
		if err := rows.Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price); err != nil {
			log.Printf("Ошибка при обработке данных: %v\n", err)
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
			log.Printf("Запись с ID %s не найдена\n", id)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Запись не найдена"})
		} else {
			log.Printf("Ошибка при получении записи с ID %s: %v\n", id, err)
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
		log.Printf("Ошибка при разборе данных: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверные входные данные"})
		return
	}

	if newAviation.Price <= 0 {
		log.Println("Некорректная цена при добавлении записи")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Некорректная цена"})
		return
	}

	query := "INSERT INTO aviation (title, plane, price) VALUES ($1, $2, $3)"
	_, err := conn.Exec(context.Background(), query, newAviation.Title, newAviation.Plane, newAviation.Price)
	if err != nil {
		log.Printf("Ошибка при добавлении записи: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления записи"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAviation)
}

// deleteAviation удаляет запись по ID.
func deleteAviationByID(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM aviation WHERE id = $1"

	commandTag, err := conn.Exec(context.Background(), query, id)
	if err != nil {
		log.Printf("Ошибка при удалении записи с ID %s: %v\n", id, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка удаления записи"})
		return
	}

	if commandTag.RowsAffected() == 0 {
		log.Printf("Запись с ID %s не найдена для удаления\n", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Запись не найдена"})
		return
	}

	log.Printf("Запись с ID %s успешно удалена\n", id)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Запись успешно удалена"})
}

// putAviation обновляет существующую запись.
func putAviation(c *gin.Context) {
	id := c.Param("id")
	var updatedAviation Aviation

	// Считываем данные для обновления из тела запроса.
	if err := c.BindJSON(&updatedAviation); err != nil {
		log.Printf("Ошибка при разборе данных для обновления записи с ID %s: %v\n", id, err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверные входные данные"})
		return
	}

	// Проверяем корректность цены.
	if updatedAviation.Price <= 0 {
		log.Println("Некорректная цена при обновлении записи")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Некорректная цена"})
		return
	}

	// Выполняем обновление записи.
	query := "UPDATE aviation SET title = $1, plane = $2, price = $3 WHERE id = $4"
	cmdTag, err := conn.Exec(context.Background(), query, updatedAviation.Title, updatedAviation.Plane, updatedAviation.Price, id)
	if err != nil {
		log.Printf("Ошибка при обновлении записи с ID %s: %v\n", id, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обновления записи"})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		log.Printf("Запись с ID %s не найдена для обновления\n", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Запись не найдена"})
		return
	}

	// Возвращаем обновленную запись.
	log.Printf("Запись с ID %s успешно обновлена\n", id)
	c.IndentedJSON(http.StatusOK, updatedAviation)
}
