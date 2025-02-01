package handler

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5" // Библиотека для работы с PostgreSQL.
	"github.com/pkg/errors"
)

type Aviation struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Plane string  `json:"plane"`
	Price float64 `json:"price"`
}

type Handle struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Handle {
	handle := Handle{
		conn: conn,
	}

	return &handle
}

// postAviation добавляет новую запись.
func (h *Handle) PostAviation(c *gin.Context) {
	var newAviation Aviation
	if err := c.BindJSON(&newAviation); err != nil {
		slog.Error("Ошибка при разборе данных", "error", errors.Wrap(err, "postaviation error with binJSON"))
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверные входные данные"})
		return
	}

	if newAviation.Price <= 0 {
		slog.Error("Некорректная цена при добавлении записи")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Некорректная цена"})
		return
	}

	query := "INSERT INTO aviation (title, plane, price) VALUES ($1, $2, $3)"
	_, err := h.conn.Exec(context.Background(), query, newAviation.Title, newAviation.Plane, newAviation.Price)
	if err != nil {
		slog.Error("Ошибка при добавлении записи", "error", errors.Wrap(err, "query insert"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка добавления записи"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAviation)
}

// getAllAviation возвращает список всех записей из таблицы aviation.
func (h *Handle) GetAllAviation(c *gin.Context) {
	query := "SELECT id, title, plane, price FROM aviation"
	rows, err := h.conn.Query(context.Background(), query)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, " function getAllAviation query?"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения данных"})
		return
	}
	defer rows.Close()

	var aviationData []Aviation
	for rows.Next() {
		var aviation Aviation
		if err := rows.Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price); err != nil {
			slog.Error("Ошибка при обработки данных", "error", errors.Wrap(err, "rows scan, тип данных"))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обработки данных"})
			return
		}
		aviationData = append(aviationData, aviation)
	}

	c.IndentedJSON(http.StatusOK, aviationData)
}

// getAviationByID возвращает запись по ID.
func (h *Handle) GetAviationByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, title, plane, price FROM aviation WHERE id = $1"

	var aviation Aviation
	err := h.conn.QueryRow(context.Background(), query, id).Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price)
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

// deleteAviation удаляет запись по ID.
func (h *Handle) DeleteAviationByID(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM aviation WHERE id = $1"

	commandTag, err := h.conn.Exec(context.Background(), query, id)
	if err != nil {
		slog.Error("error while deleting data", "error", errors.Wrap(err, "func delete aviation query"))
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
func (h *Handle) PutAviation(c *gin.Context) {
	id := c.Param("id")
	var updatedAviation Aviation

	// Считываем данные для обновления из тела запроса.
	if err := c.BindJSON(&updatedAviation); err != nil {
		slog.Error("Ошибка при разборе данных для обновления записи с ID", "error", errors.Wrap(err, "updated aviation binjson"))
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
	cmdTag, err := h.conn.Exec(context.Background(), query, updatedAviation.Title, updatedAviation.Plane, updatedAviation.Price, id)
	if err != nil {
		slog.Error("Ошибка при обновлении записи с ID", "error", errors.Wrap(err, "query for update aviation"))
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
