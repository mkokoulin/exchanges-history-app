package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/KokoulinM/exchanges-history-app/internal/csv"
	"github.com/KokoulinM/exchanges-history-app/internal/models"
	"github.com/go-chi/chi/v5"
)

type Repository interface {
	Ping(ctx context.Context) error
	UploadFile(ctx context.Context, exchangesHistory []models.ExchangesHistory) error
	GetHistory(ctx context.Context) ([]models.ExchangesHistory, error)
	Calculate(ctx context.Context, from, to, payMethod, cryptoCurrency string) (models.ResponseCalculation, error)
	GetInfo(ctx context.Context) (models.ResponseExchangesHistoryInfo, error)
}

type Handlers struct {
	repo    Repository
	baseURL string
}

// New is the handlers constructor
func New(repo Repository, baseURL string) *Handlers {
	return &Handlers{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (h *Handlers) UploadHistory(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	fileName := chi.URLParam(r, "file")
	if fileName == "" {
		http.Error(w, "the parameter is missing", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(buf.Bytes())

	data, err := csv.Reader(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.repo.UploadFile(r.Context(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handlers) GetHistory(w http.ResponseWriter, r *http.Request) {
	exchangesHistory, err := h.repo.GetHistory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(exchangesHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err == nil {
		return
	}
}

func (h *Handlers) CalculateHistory(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if from == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	payMethod := r.URL.Query().Get("payMethod")
	if payMethod == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cryptoCurrency := r.URL.Query().Get("cryptoCurrency")
	if cryptoCurrency == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	calculation, err := h.repo.Calculate(r.Context(), from, to, payMethod, cryptoCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(calculation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err == nil {
		return
	}
}

func (h *Handlers) GetHistoryInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.repo.GetInfo(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err == nil {
		return
	}
}

func (h *Handlers) PingDB(w http.ResponseWriter, r *http.Request) {
	err := h.repo.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
