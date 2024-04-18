package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/pascaloseko/ems/internal/employees"
	"github.com/pascaloseko/ems/internal/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func splitBearer(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func Middleware(emp employees.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			//validate jwt token
			tokenStr := splitBearer(header)
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := employees.Employee{Username: username}
			id, err := emp.GetEmployeeIdByUsername(r.Context(), username)
			if err != nil || id == 0 {
				log.Printf("GetEmployeeIdByUsername: %v", err)
				http.Error(w, "Invalid token: user not found", http.StatusForbidden)
				return
			}

			user.ID = id
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			// Call the next handler, which can use the user in its request context.
			next.ServeHTTP(w, r)
		})
	}

}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *employees.Employee {
	raw, _ := ctx.Value(userCtxKey).(*employees.Employee)
	return raw
}
