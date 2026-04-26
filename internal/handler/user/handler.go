package user

import (
	"encoding/json"
	"html/template"
	"net/http"

	"sample-web-http/internal/authenticate"
	"sample-web-http/internal/user"
)

type SignUp struct {
	Username         string
	Password         string
	RepeatedPassword string
}

type SignIn struct {
	Username string
	Password string
}

type SignInResponse struct {
	Token string `json:"token"`
}

type Handler struct {
	UserService *user.Service
	AuthService *authenticate.TokenService
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/signup-page", h.SignUpPage)
	mux.HandleFunc("/signin-page", h.SignInPage)
	mux.HandleFunc("/signup", h.SignUp)
	mux.HandleFunc("/signin", h.SignIn)
	mux.HandleFunc("/signout", h.SignOut)
}

func (h *Handler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, nil)
}

func (h *Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signIn.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, nil)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpVar SignUp
	var userVar user.User
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&signUpVar)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if signUpVar.Username == "" || signUpVar.Password == "" {
		http.Error(w, "username or password is empty", http.StatusBadRequest)
		return
	}
	if signUpVar.Password != signUpVar.RepeatedPassword {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}
	userVar = user.User{
		UserName: signUpVar.Username,
		Password: signUpVar.Password,
	}

	_, err = h.UserService.Create(r.Context(), &userVar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authToken, err := h.AuthService.GenerateToken(r.Context(), &userVar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := SignInResponse{
		Token: authToken,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var signInData SignIn
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	userData, err := h.UserService.Authenticate(r.Context(), signInData.Username, signInData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	authToken, err := h.AuthService.GenerateToken(r.Context(), userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := SignInResponse{
		Token: authToken,
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	// todo: remove token jwt maybe
}
