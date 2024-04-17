package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/pascaloseko/ems/graph"
	"github.com/pascaloseko/ems/internal/auth"
	"github.com/pascaloseko/ems/internal/employees"
	"github.com/pascaloseko/ems/internal/handlers"
	"github.com/pascaloseko/ems/internal/pkg/db/database"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	store := employees.NewEmployeeStore(db)
	resolver := graph.NewResolver(store)
	handlers := handlers.NewHandlers(resolver)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router.HandleFunc("/login", handlers.LoginHandler)

	// Protected Route: /employees
	router.Group(func(r chi.Router) {
		r.Use(auth.Middleware(store))
		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
		r.Handle("/query", srv)
		r.HandleFunc("/employees", handlers.GetAllEmployeesHandler)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
