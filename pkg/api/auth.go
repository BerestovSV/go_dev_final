package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenCookieName = "token"
)

// JWTClaims - кастомные claims для нашего токена
type JWTClaims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

// authMiddleware - middleware для проверки аутентификации
func (a *API) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.Password == "" {
			next(w, r)
			return
		}

		tokenString := a.getTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		valid, err := a.validateToken(tokenString, a.config.Password)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// getTokenFromRequest - извлекает токен из куки или заголовка
func (a *API) getTokenFromRequest(r *http.Request) string {
	// Пробуем получить из куки
	cookie, err := r.Cookie(tokenCookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Пробуем получить из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

// validateToken - проверяет валидность JWT токена
func (a *API) validateToken(tokenString, password string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JWTSecret), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		currentHash := a.getPasswordHash(password)
		return claims.PasswordHash == currentHash, nil
	}

	return false, nil
}

// generateToken - генерирует JWT токен
func (a *API) generateToken(password string) (string, error) {
	claims := JWTClaims{
		PasswordHash: a.getPasswordHash(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.JWTSecret))
}

func (a *API) getPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// signinHandler - обработчик аутентификации
func (a *API) signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
		return
	}

	type SignInRequest struct {
		Password string `json:"password"`
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid JSON"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSON(w, http.StatusInternalServerError, errResp{Error: "authentication not configured"})
		return
	}

	if req.Password != expectedPassword {
		writeJSON(w, http.StatusUnauthorized, errResp{Error: "wrong password"})
		return
	}

	token, err := a.generateToken(expectedPassword)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{Error: "failed to generate token"})
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(a.config.TokenDuration),
		HttpOnly: true,
		Path:     "/",
	})

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
