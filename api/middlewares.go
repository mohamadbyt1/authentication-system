package api
import (
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err != nil {
			ResponseJson(w,http.StatusForbidden,map[string]string{"Error": "permission denied"})
			return
		}
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte("secretkey"), nil
		})
		if err != nil || !token.Valid {
			ResponseJson(w,http.StatusForbidden,map[string]string{"Error": "permission denied"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
			if time.Now().After(expirationTime) {
				ResponseJson(w,http.StatusForbidden,map[string]string{"Error": "permission denied"})
				return
			}
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)
		}
		next(w, r)
	}
}
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Incomming request  Method: %s RequestUri: %s ", r.Method, r.RequestURI)
		next.ServeHTTP(w,r)
    }
}