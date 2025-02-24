package usecase

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Aviation struct {
	ID    uint64  `json:"id"`    // Уникальный идентификатор самолёта (целое число, внешний ключ на flights.id)
	Title string  `json:"title"` // Название авиакомпании или самолёта
	Plane string  `json:"plane"` // Модель самолёта
	Price float64 `json:"price"` // Цена (стоимость самолёта)
}

type AviationUsecase struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *AviationUsecase {
	return &AviationUsecase{
		conn: conn,
	}
}

func (a *AviationUsecase) CreateAviation(ctx context.Context, aviation Aviation) error {
	query := "INSERT INTO aviation (title, plane, price) VALUES ($1, $2, $3)"
	_, err := a.conn.Exec(context.Background(), query, aviation.Title, aviation.Plane, aviation.Price)
	if err != nil {
		slog.Error("Ошибка при добавлении записи", "error", errors.Wrap(err, "query insert"))
		return errors.Wrap(err, "query insert")
	}
	return nil
}

func (a *AviationUsecase) GetAllAviation(ctx context.Context) ([]Aviation, error) {
	query := "SELECT id, title, plane, price FROM aviation"
	rows, err := a.conn.Query(context.Background(), query)
	if err != nil {
		slog.Error("Ошибка при получении данных", "error", errors.Wrap(err, " function getAllAviation query?"))
		return nil, errors.Wrap(err, "function getAllAviation query")
	}
	defer rows.Close()

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
