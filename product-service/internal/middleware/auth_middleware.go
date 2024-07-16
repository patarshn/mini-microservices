package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"product-service/internal/constant"
	"product-service/internal/models"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		secretKey := []byte(os.Getenv("SECRET_KEY"))
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ResponseBody{
				Error:   true,
				Message: "Authorization header missing",
			})
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ResponseBody{
				Error:   true,
				Message: "Invalid token format",
			})
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ResponseBody{
				Error:   true,
				Message: "Err token",
			})
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ResponseBody{
				Error:   true,
				Message: "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ResponseBody{
				Error:   true,
				Message: "Fail claims token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), constant.USERINFO, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
