package auth

import (
	"context"
	"net/http"

	"github.com/pascaloseko/ems/internal/employees"
	"github.com/pascaloseko/ems/internal/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware(emp employees.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := employees.Employee{Username: username}
			id, err := emp.GetEmployeeIdByUsername(r.Context(), username)
			if err != nil {
				next.ServeHTTP(w, r)
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
