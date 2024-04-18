package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pascaloseko/ems/internal/employees"
	"github.com/pascaloseko/ems/internal/mockdb"
	"github.com/pascaloseko/ems/internal/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	authorizationType string,
	username string,
) {
	tkn, err := jwt.GenerateToken(username)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, tkn)
	request.Header.Set("Authorization", authorizationHeader)
}

func TestMiddleware(t *testing.T) {
	type args struct {
		buildStubs func(store *mockdb.MockStore)
	}
	tests := []struct {
		name          string
		args          args
		want          *employees.Employee
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		setupAuth     func(t *testing.T, request *http.Request)
	}{
		{
			name: "valid token",
			args: args{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetEmployeeIdByUsername(gomock.Any(), "pascal").AnyTimes().Return(int64(1), nil)
				},
			},
			want: &employees.Employee{
				Username: "pascal",
				ID:       int64(1),
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "bearer", "pascal")
			},
		},
		{
			name: "invalid token",
			args: args{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().GetEmployeeIdByUsername(gomock.Any(), "not a valid user").AnyTimes().Return(int64(0), nil)
				},
			},
			want: nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request) {
			},
		},
		{
			name: "no token",
			args: args{
				buildStubs: func(store *mockdb.MockStore) {
					// No expectation for this case as the function should not be called
				},
			},
			want: nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.args.buildStubs(store)

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			tt.setupAuth(t, req)
			rr := httptest.NewRecorder()
			handler := Middleware(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := ForContext(r.Context())
				assert.Equal(t, tt.want, user)
			}))
			handler.ServeHTTP(rr, req)
			tt.checkResponse(t, rr)
		})
	}
}
