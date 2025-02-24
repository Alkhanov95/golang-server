package handler

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5" // Библиотека для работы с PostgreSQL.
	"github.com/pkg/errors"
)

type Aviation struct {
	ID    uint64  `json:"id"`    // Уникальный идентификатор самолёта (целое число, внешний ключ на flights.id)
	Title string  `json:"title"` // Название авиакомпании или самолёта
	Plane string  `json:"plane"` // Модель самолёта
	Price float64 `json:"price"` // Цена (стоимость самолёта)
}

type Flights struct {
	ID            uint64 `json:"id"` // Уникальный идентификатор рейса (целое число)
	PlaneID       uint64 `json:"plane_id"`
	FlightsNumber string `json:"flights_number"` // Номер рейса
	Departure     string `json:"departure"`      // Аэропорт вылета
	Arrival       string `json:"arrival"`        // Аэропорт прилёта
}

type Handle struct {
	conn *pgx.Conn
}

type PlaneID struct {
	ID            uint64  `json:"id"`    // Уникальный идентификатор самолёта (целое число, внешний ключ на flights.id)
	Title         string  `json:"title"` // Название авиакомпании или самолёта
	Plane         string  `json:"plane"` // Модель самолёта
	Price         float64 `json:"price"` // Цена (стоимость самолёта)
	PlaneID       uint64  `json:"plane_id"`
	FlightsNumber string  `json:"flights_number"` // Номер рейса
	Departure     string  `json:"departure"`      // Аэропорт вылета
	Arrival       string  `json:"arrival"`        // Аэропорт прилёта
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

func (h *Handle) GetAllFlights(c *gin.Context) {
	// Define the SQL query
	query := "SELECT id, plane_id, flights_number, departure, arrival FROM flights"
	// Execute the query
	rows, err := h.conn.Query(context.Background(), query)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, "function getAllFlights query"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения данных"})
		return
	}
	defer rows.Close()

	// Initialize a slice to hold the flight data
	var flightsData []Flights

	// Iterate through the rows and populate the flightData slice
	for rows.Next() {
		var flights Flights
		if err := rows.Scan(&flights.ID, &flights.PlaneID, &flights.FlightsNumber, &flights.Departure, &flights.Arrival); err != nil {
			slog.Error("Ошибка при обработки данных", "error", errors.Wrap(err, "rows scan, тип данных"))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка обработки данных"})
			return
		}
		flightsData = append(flightsData, flights)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		slog.Error("Ошибка при итерации по строкам", "error", errors.Wrap(err, "rows iteration error"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка итерации по данным"})
		return
	}

	// Return the flight data as a JSON response
	c.IndentedJSON(http.StatusOK, flightsData)

}

func (h *Handle) GetFlightsByID(c *gin.Context) {
	id := c.Param("id")
	query := "SELECT id, flights_number, departure, arrival FROM flights WHERE id = $1"

	var flights Flights
	err := h.conn.QueryRow(context.Background(), query, id).Scan(&flights.ID, &flights.FlightsNumber, &flights.Departure, &flights.Arrival)
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

	c.IndentedJSON(http.StatusOK, flights)
}

func (h *Handle) GetPlaneStatsByID(c *gin.Context) {
	var aviation Aviation
	var flights []Flights

	// Получение planeID из параметров запроса
	planeIDStr := c.Param("plane_id")
	planeID, err := strconv.ParseInt(planeIDStr, 10, 64)
	if err != nil {
		slog.Error("Неверный формат plane_id", "error", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Неверный формат plane_id"})
		return
	}

	// Получение информации о самолете
	queryAviation := `SELECT id, title, plane, price FROM aviation WHERE id = $1`
	err = h.conn.QueryRow(context.Background(), queryAviation, planeID).Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Самолет не найден"})
			return
		}
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, "getplanestatsbyid error with QueryRow aviation"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении записи"})
		return
	}

	// Получение информации о рейсах
	queryFlights := `SELECT id, plane_id, flights_number, departure, arrival FROM flights WHERE plane_id = $1`
	rows, err := h.conn.Query(context.Background(), queryFlights, planeID)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, "getplanestatsbyid error with Query flights"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении записи"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var flight Flights
		err := rows.Scan(&flight.ID, &flight.PlaneID, &flight.FlightsNumber, &flight.Departure, &flight.Arrival)
		if err != nil {
			slog.Error("Ошибка при сканировании данных", "error", errors.Wrap(err, "getplanestatsbyid error with Scan flights"))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении записи"})
			return
		}
		flights = append(flights, flight)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Ошибка при итерации по строкам", "error", errors.Wrap(err, "getplanestatsbyid error with rows iteration"))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при получении записи"})
		return
	}

	// Формирование ответа
	stats := struct {
		Aviation Aviation  `json:"aviation"`
		Flights  []Flights `json:"flights"`
	}{
		Aviation: aviation,
		Flights:  flights,
	}

	c.IndentedJSON(http.StatusOK, stats)
}

// func (h *Handle) GetPlaneStatsByID(c *gin.Context) {
// 	id := c.Param("plane_id")
// 	query := "SELECT id, title, plane, price FROM aviation WHERE id = $1"

// 	var aviation Aviation
// 	err := h.conn.QueryRow(context.Background(), query, id).Scan(&aviation.ID, aviation.Title, aviation.Plane, aviation.Price)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			log.Printf("запись с plane_id %s не найдена", id)
// 			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "запись не найдена"})
// 		} else {
// 			log.Printf("Ошибка при получении записи с plane id %s: %v", id, err)
// 			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Ошибка получения записи"})
// 		}
// 		return
// 	}

// 	var flights Flights
// 	queryf := "SELECT id, flights_number, departure, arrival FROM flights WHERE plane_id = $1"
// 	err = h.conn.QueryRow(context.Background(), queryf, id).Scan(&flights.ID, flights.FlightsNumber, flights.Departure, flights.Arrival)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			log.Printf("запись с plane_id %s не найдена", id)
// 			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "запись не найдена"})
// 		} else {
// 			log.Printf("Ошибка при получении записи с plane id %s: %v", id, err)
// 		}
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, //)
// }

// 2 requests
// for handler stats we need 2 request select * from aviation WHERE  id = planeID
// second req : select * from flights weher planeID = planeID
