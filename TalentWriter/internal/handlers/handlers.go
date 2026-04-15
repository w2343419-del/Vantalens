package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"vantalens/talentwriter/internal/auth"
	"vantalens/talentwriter/internal/comment"
	"vantalens/talentwriter/internal/config"
	"vantalens/talentwriter/internal/models"
)

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:1313",
			"http://localhost:9090",
			"http://localhost:9091",
			"http://127.0.0.1:1313",
			"http://127.0.0.1:9090",
			"http://127.0.0.1:9091",
			"https://w2343419-del.github.io",
		}
		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		h(w, r)
	}
}

type loginRequest struct {
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	username := strings.TrimSpace(req.User)
	if username == "" {
		username = strings.TrimSpace(req.Username)
	}
	password := req.Pass
	if password == "" {
		password = req.Password
	}
	cfg := config.GetConfig()
	if cfg == nil || username != "admin" {
		RespondJSON(w, 401, models.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	if cfg.AdminToken != "" && password != cfg.AdminToken {
		RespondJSON(w, 401, models.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	if cfg.AdminToken == "" && strings.TrimSpace(password) == "" {
		RespondJSON(w, 401, models.APIResponse{Success: false, Message: "Unauthorized"})
		return
	}
	accessToken, _ := auth.CreateJWT("admin", "access")
	refreshToken, _ := auth.CreateJWT("admin", "refresh")
	RespondJSON(w, 200, models.APIResponse{Success: true, Data: map[string]string{"token": accessToken, "access_token": accessToken, "refresh_token": refreshToken}})
}

func HandleGetComments(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	comments, _ := comment.GetComments(path)
	RespondJSON(w, 200, models.APIResponse{Success: true, Data: comments})
}

type addCommentRequest struct {
	Author  string `json:"author"`
	Email   string `json:"email"`
	Content string `json:"content"`
	Parent  string `json:"parent"`
}

func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	var req addCommentRequest
	json.NewDecoder(r.Body).Decode(&req)
	c, err := comment.AddComment(path, req.Author, req.Email, req.Content, r.RemoteAddr, r.UserAgent(), req.Parent)
	if err != nil {
		RespondJSON(w, 500, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	RespondJSON(w, 200, models.APIResponse{Success: true, Data: c})
}

func HandleApproveComment(w http.ResponseWriter, r *http.Request) {
	if !auth.RequireAuth(w, r) {
		return
	}
	path := r.URL.Query().Get("path")
	id := r.URL.Query().Get("id")
	if err := comment.ApproveComment(path, id); err != nil {
		RespondJSON(w, 500, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	RespondJSON(w, 200, models.APIResponse{Success: true})
}

func HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	if !auth.RequireAuth(w, r) {
		return
	}
	path := r.URL.Query().Get("path")
	id := r.URL.Query().Get("id")
	if err := comment.DeleteComment(path, id); err != nil {
		RespondJSON(w, 500, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	RespondJSON(w, 200, models.APIResponse{Success: true})
}

func HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings := comment.LoadSettings()
	RespondJSON(w, 200, models.APIResponse{Success: true, Data: settings})
}

func HandleSaveSettings(w http.ResponseWriter, r *http.Request) {
	if !auth.RequireAuth(w, r) {
		return
	}
	var settings models.CommentSettings
	json.NewDecoder(r.Body).Decode(&settings)
	if err := comment.SaveSettings(settings); err != nil {
		RespondJSON(w, 500, models.APIResponse{Success: false, Message: err.Error()})
		return
	}
	RespondJSON(w, 200, models.APIResponse{Success: true})
}
