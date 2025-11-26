package save

import (
	"errors"
	"net/http"

	"log/slog"

	response "GoURLShortener/internal/lib/api/response"
	"GoURLShortener/internal/lib/logger/sl"
	"GoURLShortener/internal/lib/random"
	"GoURLShortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name URLSaver --output ./mocks --outpkg mocks
type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int, error)
}

const aliasLength = 6 //ToDO: move to config

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalide request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength) //ToDo check for uniqueness
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)

		if errors.Is(err, storage.ErrURLAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Info("url already exists", sl.Err(err))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		log.Info("url saved", slog.Int("id", id), slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
