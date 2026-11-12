package article

import (
	"net/http"

	"github.com/grigori-zagladkin/web-work/services/api/internal/apperr"
	"github.com/grigori-zagladkin/web-work/services/api/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) ListPublished(writer http.ResponseWriter, request *http.Request) error {
	pagination, err := httpx.ParsePagination(request)
	if err != nil {
		return err
	}
	items, total, err := handler.service.List(request.Context(), ListFilter{
		Status: StatusPublished,
		Tag:    request.URL.Query().Get("tag"),
		Limit:  pagination.Limit,
		Offset: pagination.Offset(),
	})
	if err != nil {
		return err
	}
	httpx.WriteJSON(writer, http.StatusOK, httpx.NewPage(items, pagination, total))
	return nil
}

func (handler *Handler) GetPublished(writer http.ResponseWriter, request *http.Request) error {
	article, err := handler.service.GetPublished(request.Context(), request.PathValue("slug"))
	if err != nil {
		return err
	}
	httpx.WriteJSON(writer, http.StatusOK, article)
	return nil
}

func (handler *Handler) List(writer http.ResponseWriter, request *http.Request) error {
	pagination, err := httpx.ParsePagination(request)
	if err != nil {
		return err
	}
	status := request.URL.Query().Get("status")
	if !httpx.OneOf(status, StatusDraft, StatusPublished) {
		return apperr.BadRequest("status must be draft or published")
	}
	items, total, err := handler.service.List(request.Context(), ListFilter{
		Status: status,
		Tag:    request.URL.Query().Get("tag"),
		Limit:  pagination.Limit,
		Offset: pagination.Offset(),
	})
	if err != nil {
		return err
	}
	httpx.WriteJSON(writer, http.StatusOK, httpx.NewPage(items, pagination, total))
	return nil
}

func (handler *Handler) Get(writer http.ResponseWriter, request *http.Request) error {
	article, err := handler.service.Get(request.Context(), request.PathValue("id"))
	if err != nil {
		return err
	}
	httpx.WriteJSON(writer, http.StatusOK, article)
	return nil
}

func (handler *Handler) Create(writer http.ResponseWriter, request *http.Request) error {
	var input CreateInput
	if err := httpx.DecodeJSON(writer, request, &input); err != nil {
		return err
	}
	article, err := handler.service.Create(request.Context(), input)
	if err != nil {
		return err
	}
	writer.Header().Set("Location", "/api/v1/admin/articles"+article.ID)
	httpx.WriteJSON(writer, http.StatusCreated, article)
	return nil
}

func (handler *Handler) Update(writer http.ResponseWriter, request *http.Request) error {
	var input UpdateInput
	if err := httpx.DecodeJSON(writer, request, &input); err != nil {
		return err
	}
	article, err := handler.service.Update(request.Context(), request.PathValue("id"), input)
	if err != nil {
		return err
	}
	writer.Header().Set("Location", "/api/v1/admin/articles"+article.ID)
	httpx.WriteJSON(writer, http.StatusOK, article)
	return nil
}

func (handler *Handler) Delete(writer http.ResponseWriter, request *http.Request) error {
	if err := handler.service.Delete(request.Context(), request.PathValue("id")); err != nil {
		return err
	}
	writer.WriteHeader(http.StatusNoContent)
	return nil
}
