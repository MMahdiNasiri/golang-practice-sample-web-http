package todo

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"sample-web-http/internal/middleware"
	"sample-web-http/internal/user"

	"sample-web-http/internal/todo"
)

type Handler struct {
	TodoService    *todo.Service
	UserService    *user.Service
	AuthMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

type Page struct {
	Title string
	Body  string
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", middleware.RateLimit(h.PageHandler))
	mux.HandleFunc("/all/", h.AuthMiddleware(h.ListContext))
	mux.HandleFunc("/create", h.AuthMiddleware(h.CreateContext))
	mux.HandleFunc("/update", h.AuthMiddleware(h.UpdateContext))
	mux.HandleFunc("/delete", h.AuthMiddleware(h.DeleteContext))
}

func (h *Handler) PageHandler(w http.ResponseWriter, r *http.Request) {
	pageVar := Page{
		Title: "TO DO APP",
		Body:  "why.",
	}
	tmpl, err := template.ParseFiles("templates/page.html")
	if err != nil {
		return
	}
	err = tmpl.Execute(w, pageVar)
	if err != nil {
		return
	}
}

func (h *Handler) ListContext(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	result, err := h.TodoService.List(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&todoVar)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(int)

	result, err := h.TodoService.Create(r.Context(), todoVar.Text, userID)
	if err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&todoVar)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	result, err := h.TodoService.Update(r.Context(), todoVar.ID, todoVar.Text)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&todoVar)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.TodoService.Delete(r.Context(), todoVar.ID)
	if err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(todoVar)
	if err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}
