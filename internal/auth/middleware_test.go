package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pascaloseko/ems/internal/employees"
	"github.com/pascaloseko/ems/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name string
		args args
		want *employees.Employee
	}{
		{
			name: "valid token",
			args: args{
				header: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBhc2NhbCJ9.sdfsdfsdf",
			},
			want: &employees.Employee{
				Username: "pascal",
				ID:       1,
			},
		},
		{
			name: "invalid token",
			args: args{
				header: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBhc2NhbCJ9.sdfsdfsdf",
			},
			want: nil,
		},
		{
			name: "no token",
			args: args{
				header: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := jwt.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBhc2NhbCJ9.sdfsdfsdf")
			assert.NoError(t, err)
			id, err := employees.GetEmployeeIdByUsername(username)
			assert.NoError(t, err)

			assert.Equal(t, id, 1)
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.args.header)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := r.Context().Value(userCtxKey).(*employees.Employee)
				if user != nil {
					rr.Header().Set("user", user.Username)
				}
			})
			Middleware()(handler).ServeHTTP(rr, req)
			if tt.want != nil {
				if rr.Header().Get("user") != tt.want.Username {
					t.Errorf("Middleware() = %v, want %v", rr.Header().Get("user"), tt.want.Username)
				}
			} else {
				if rr.Header().Get("user") != "" {
					t.Errorf("Middleware() = %v, want %v", rr.Header().Get("user"), "")
				}
			}
		})
	}
}
