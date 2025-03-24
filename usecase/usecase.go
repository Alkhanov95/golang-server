package usecase

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

type AviationUsecase struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *AviationUsecase {
	return &AviationUsecase{
		conn: conn,
	}
}

func (a *AviationUsecase) GetAllAviation(ctx context.Context) ([]Aviation, error) {
	query := "SELECT id, title, plane, price FROM aviation"
	rows, err := a.conn.Query(context.Background(), query)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, " function getAllAviation query?"))
		return nil, errors.Wrap(err, "function getAllAviation query")
	}
	defer rows.Close()
    //fmt.Println(rows.)
	var aviationData []Aviation
	for rows.Next() {
		var aviation Aviation
		if err := rows.Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price); err != nil {
			slog.Error("Ошибка при обработки данных", "error", errors.Wrap(err, "rows scan, тип данных"))
			return nil, errors.Wrap(err, "rows scan")
		}
		aviationData = append(aviationData, aviation)
	}

	return aviationData, nil
}

// GetAviationByID получает запись о самолёте по ID.
func (a *AviationUsecase) getAviationByID(ctx context.Context, id string) (*Aviation, error) {

	query := "SELECT id, title, plane, price FROM aviation WHERE id = $1"

	var aviation Aviation
	err := a.conn.QueryRow(ctx, query, id).Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Запись не найдена", "id", id)
			return nil, errors.New("запись с указанным ID не найдена")
		}

		slog.Error("Ошибка при выполнении запроса", "id", id, "error", err)
		return nil, errors.Wrap(err, "ошибка при получении записи")
	}

	return &aviation, nil
}

func (a *AviationUsecase) DeleteAviationByID(ctx context.Context, id string) error {
	query := "DELETE FROM aviation WHERE id = $1"

	commandTag, err := a.conn.Exec(context.Background(), query, id)
	if err != nil {
		slog.Error("error while deleting data", "error", errors.Wrap(err, "func delete aviation query"))

		return err
	}

	if commandTag.RowsAffected() == 0 {
		slog.Error("Запись с ID %s не найдена для удаления\n", id)
		return err
	}
	return nil
}

func (a *AviationUsecase) PutAviation(ctx context.Context, id string) error {
	var updatedAviation Aviation
	query := "UPDATE aviation SET title = $1, plane = $2, price = $3 WHERE id = $4"
	cmdTag, err := a.conn.Exec(context.Background(), query, updatedAviation.Title, updatedAviation.Plane, updatedAviation.Price, id)
	if err != nil {
		slog.Error("Ошибка при обновлении записи с ID", "error", errors.Wrap(err, "query for update aviation"))

		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Printf("Запись с ID %s не найдена для обновления\n", id)

		return err
	}

	// Возвращаем обновленную запись.
	log.Printf("Запись с ID %s успешно обновлена\n", id)

	return nil
}

func (a *AviationUsecase) GetAllFlights(ctx context.Context) ([]Flights, error) {
	
	// Define the SQL query
	query := "SELECT id, plane_id, flights_number, departure, arrival FROM flights"
	// Execute the query
	rows, err := a.conn.Query(context.Background(), query)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, "function getAllFlights query"))

		 return nil, errors.Wrap(err, "function getAllFlights query")
	}
	defer rows.Close()

	// Initialize a slice to hold the flight data
	var flightsData []Flights
	for rows.Next() {
		var flights Flights
		if err := rows.Scan(&flights.ID, &flights.PlaneID, &flights.FlightsNumber, &flights.Departure, &flights.Arrival); err != nil {
			slog.Error("Ошибка при обработки данных", "error", errors.Wrap(err, "rows scan, тип данных"))
		    return nil, errors.Wrap(err, "func get all flights query")
		}
		flightsData = append(flightsData, flights)
	}
	return flightsData, nil
	

}



func (a *AviationUsecase) GetFlightsByID(ctx context.Context, id string) (Flights, error) {
	
	query := "SELECT id, flights_number, departure, arrival FROM flights WHERE id = $1"

	var flights Flights
	err := a.conn.QueryRow(context.Background(), query, id).Scan(&flights.ID, &flights.FlightsNumber, &flights.Departure, &flights.Arrival)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Запись с ID %s не найдена\n", id)
		} else {
			log.Printf("Ошибка при получении записи с ID %s: %v\n", id, err)
		}
		return flights, err
	}

	return flights, nil
}

func (a *AviationUsecase) GetPlaneStatsByID(ctx context.Context, planeIDStr string) (/*data type*/[]Flights, []Aviation, error)  {
	var aviation []Aviation
	var flights []Flights

	// Получение planeID из параметров запроса
	
	planeID, err := strconv.ParseInt(planeIDStr, 10, 64)
	if err != nil {
		slog.Error("Неверный формат plane_id", "error", err)
		err = errors.Wrap(err, "GetPlaneStatsByID ParseInt") //запусти и посмотри как это выглядит
		return flights, aviation, err
		}

	// Получение информации о самолете
	queryAviation := `SELECT id, title, plane, price FROM aviation WHERE id = $1`
	err = a.conn.QueryRow(context.Background(), queryAviation, planeID).Scan(&aviation.ID, &aviation.Title, &aviation.Plane, &aviation.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return flights, aviation, err
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
