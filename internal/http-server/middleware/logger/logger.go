package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("Logger middleware enabled")

		// обработчик

		fn := func(w http.ResponseWriter, r *http.Request) {

			// собираем исходную информацию о запросе

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Создаем обертку вокруг `http.ResponseWriter`
			// для получения сведений об ответе

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Момент получения запроса, что бы вычислить время обработки

			start := time.Now()

			// Запись отправится в лог defer
			// в этот момент запрос уже будет обработан

			defer func() {
				entry.Info("request complete",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(start).String()),
				)
			}()

			// Передаем управление следующему обработчику в цепочке middleware

			next.ServeHTTP(ww, r)
		}

		// Возвращаем созданный выше обработчик, приводя его к типу http.HandlerFunc

		return http.HandlerFunc(fn)
	}
}
