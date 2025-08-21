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
	tokenDuration   = 8 * time.Hour
)

// JWTClaims - кастомные claims для нашего токена
type JWTClaims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

// authMiddleware - middleware для проверки аутентификации
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, требуется ли аутентификация
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			// Аутентификация не требуется
			next(w, r)
			return
		}

		// Получаем токен из куки
		tokenString := getTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен
		valid, err := validateToken(tokenString, password)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// getTokenFromRequest - извлекает токен из куки или заголовка
func getTokenFromRequest(r *http.Request) string {
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
func validateToken(tokenString, password string) (bool, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getPasswordHash(password)), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Проверяем, что хэш пароля в токене совпадает с текущим
		currentHash := getPasswordHash(password)
		return claims.PasswordHash == currentHash, nil
	}

	return false, nil
}

// generateToken - генерирует JWT токен
func generateToken(password string) (string, error) {
	claims := JWTClaims{
		PasswordHash: getPasswordHash(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getPasswordHash(password)))
}

// getPasswordHash - возвращает хэш пароля
func getPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// signinHandler - обработчик аутентификации
func signinHandler(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusUnauthorized, errResp{Error: "Неверный пароль"})
		return
	}

	token, err := generateToken(expectedPassword)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{Error: "failed to generate token"})
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     tokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(tokenDuration),
		HttpOnly: true,
		Path:     "/",
	})

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
