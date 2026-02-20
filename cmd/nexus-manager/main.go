package main

import (
	"log"
	"net/http"
	"os"

	"nexus/internal/manager/graph"
	"nexus/internal/manager/graph/model"
	"nexus/internal/manager/ui"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize with some dummy data for prototyping
	resolver := &graph.Resolver{
		Agents: []*model.Agent{
			{
				ID:       "agent-001",
				Hostname: "app-server-01",
				LastSeen: "2026-02-19T21:00:00Z",
				Version:  "v1.0.0",
				Status:   "ONLINE",
			},
		},
		Destinations: []*model.DestinationConfig{
			{
				ID:      "dest-0",
				Name:    "Crowdstrike Falcon Streaming",
				Type:    "webhook",
				URL:     "https://api.crowdstrike.com/ingest",
				Enabled: false,
			},
		},
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// GraphQL Endpoints
	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Serve Embedded UI Dashboard
	fs := http.FileServer(http.FS(ui.StaticFS))
	http.Handle("/static/", fs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			content, _ := ui.StaticFS.ReadFile("static/index.html")
			w.Write(content)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Printf("Nexus Control Plane running on http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
