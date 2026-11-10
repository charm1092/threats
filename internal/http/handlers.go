package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"kaspersky/internal/threat"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	repo *threat.Repository
}

func NewHTTPHandlers(repo *threat.Repository) *HTTPHandlers {
	return &HTTPHandlers{
		repo: repo,
	}
}

func (h *HTTPHandlers) HandleCreateThreat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var threatDTO ThreatDTO
	if err := json.NewDecoder(r.Body).Decode(&threatDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := threatDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	threatVar, err := h.repo.CreateThreat(
		ctx,
		threatDTO.MD5,
		threatDTO.Name,
		threatDTO.Description,
		threatDTO.Severity,
	)

	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, threat.ErrMD5IsNotAllowed) || errors.Is(err, threat.ErrSeverityIsNotAllowed) {
			writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, threat.ErrMD5Exists) {
			writeJSONError(w, errDTO.ToString(), http.StatusConflict)
		} else {
			writeJSONError(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(threatVar, "", "	")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("не удалось записать http-ответ:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetByMD5(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	md5 := mux.Vars(r)["md5"]

	threatVar, err := h.repo.GetByMD5(ctx, md5)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, threat.ErrMD5NotFound) {
			writeJSONError(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, threat.ErrMD5IsNotAllowed) {
			writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			writeJSONError(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(threatVar, "", "	")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("не удалось записать http-ответ:", err)
		return
	}
}

func (h *HTTPHandlers) HandlerUpdateThreat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	md5 := mux.Vars(r)["md5"]

	if !threat.IsMD5Allowed(md5) {
		errDTO := ErrorDTO{
			Message: "недопустимый md5",
			Time:    time.Now(),
		}
		writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	var dto UpdateThreatDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if dto.NewName == "" || dto.NewSeverity == "" {
		errDTO := ErrorDTO{
			Message: "nothing to update",
			Time:    time.Now(),
		}

		writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	err := h.repo.UpdateThreat(
		ctx,
		md5,
		dto.NewName,
		dto.NewDescription,
		dto.NewSeverity,
	)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, threat.ErrMD5NotFound) {
			writeJSONError(w, errDTO.ToString(), http.StatusNotFound)
		} else if errors.Is(err, threat.ErrMD5IsNotAllowed) {
			writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, threat.ErrSeverityIsNotAllowed) {
			writeJSONError(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			writeJSONError(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	updateThreat, err := h.repo.GetByMD5(ctx, md5)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		writeJSONError(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(updateThreat, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := fmt.Fprintln(w, message); err != nil {
		fmt.Println("не удалось записать http-ответ:", err)
	}
}
