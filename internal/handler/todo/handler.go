package todo

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sample-web-http/internal/middleware"
	"sample-web-http/internal/user"

	"sample-web-http/internal/todo"
)

type Handler struct {
	TodoService *todo.Service
	UserService *user.Service
}

type Page struct {
	Title string
	Body  string
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", middleware.RateLimit(h.PageHandler))
	mux.HandleFunc("/all/", h.ListContext)
	mux.HandleFunc("/create", h.CreateContext)
	mux.HandleFunc("/update", h.UpdateContext)
	mux.HandleFunc("/delete", h.DeleteContext)
}

func (h *Handler) PageHandler(w http.ResponseWriter, r *http.Request) {
	pageVar := Page{
		Title: "TO DO APP",
		Body:  "no one can hear me.",
	}
	tmpl, err := template.ParseFiles("templates/page.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, pageVar)
}

func (h *Handler) ListContext(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	result, err := h.TodoService.List(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) CreateContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&todoVar)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := h.TodoService.Create(r.Context(), todoVar.Text, todoVar.CreatedBy)
	if err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) UpdateContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer r.Body.Close()

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
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) DeleteContext(w http.ResponseWriter, r *http.Request) {
	var todoVar todo.Todo
	defer r.Body.Close()

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
	json.NewEncoder(w).Encode(todoVar)
}
