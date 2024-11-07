package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias, omitempty"`
}

type Response struct {
	// Status string `json:"status"`
	// Error  string `json:"error, omitempty"`
	resp.Response
	Alias string `json:"alias, omitempty"`
}

const aliasLength = 6

type URLDelete interface {
	DeleteURL(urlDelete string, alias string) (int64, error)
}

func New(log *slog.Logger, URLDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request", middleware.GetReqID(r.Context())),
		)

		var req Request

		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := URLDelete.DeleteURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		responseOK(w, r, alias)
	}

}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
