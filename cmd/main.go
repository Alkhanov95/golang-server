package main

import (
	"aviation/my-api/internal/handler"
	"context" // Для управления контекстом выполнения.
	"log"     // Для логирования ошибок и сообщений.

	// Для работы с HTTP-сервером.
	"log/slog" // Используем для логирования

	"github.com/gin-gonic/gin" // Фреймворк для создания веб-приложений.
	"github.com/jackc/pgx/v5"  // Библиотека для работы с PostgreSQL.
	"github.com/pkg/errors"    // Используем для обёртывания ошибок
)

func getConnect(connStr string) (*pgx.Conn, error) {
	// Устанавливаем соединение с базой данных.
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		//логируем ошибку подключения
		slog.Error("Ошибка при подключении к бд", "error", errors.Wrap(err, "ошибка при установление бд PGX connect"))
		return nil, err

	}
	return conn, err
}

func main() {
	// Строка подключения
	connStr := "postgresql://postgres:@postgres:5432/postgres" // Используется имя контейнера "db"

	conn, err := getConnect(connStr)
	if err != nil {
		slog.Error("Ошибка при подключении к бд", "error", errors.Wrap(err, "ошибка при установление бд PGX connect"))
		return
	}

	defer func() {
		// Закрытие соединения при завершении работы.
		if err := conn.Close(context.Background()); err != nil {
			slog.Error("Ошибка при закрытие соеденения", "error", errors.Wrap(err, "closing connection db error (conn.close) "))
		}
	}()
	handle := handler.New(conn)
	router := getRouter(handle)
	// Запуск HTTP-сервера.
	log.Println("Запуск сервера на порту 8080...")
	router.Run("0.0.0.0:8080")
}

func getRouter(handle *handler.Handle) *gin.Engine {
	router := gin.Default()
	router.GET("/aviation", handle.GetAllAviation)
	router.GET("/aviation/:id", handle.GetAviationByID)
	router.POST("/aviation", handle.PostAviation)
	router.PUT("/aviation/:id", handle.PutAviation) // Добавляем PUT маршрут
	router.DELETE("/aviation/:id", handle.DeleteAviationByID)

	return router
}
