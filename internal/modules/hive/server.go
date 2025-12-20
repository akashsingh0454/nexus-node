package hive

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// FileServer serves cached files to other peers.
type FileServer struct {
	Port     int
	CacheDir string
	server   *http.Server
}

func NewFileServer(port int, cacheDir string) *FileServer {
	return &FileServer{
		Port:     port,
		CacheDir: cacheDir,
	}
}

func (fs *FileServer) Start() {
	if _, err := os.Stat(fs.CacheDir); os.IsNotExist(err) {
		os.MkdirAll(fs.CacheDir, 0755)
	}

	mux := http.NewServeMux()
	// Serve files from the cache directory
	fileHandler := http.StripPrefix("/cache/", http.FileServer(http.Dir(fs.CacheDir)))
	mux.Handle("/cache/", fileHandler)

	fs.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", fs.Port),
		Handler: mux,
	}

	go func() {
		log.Printf("[HIVE] File Server listening on :%d", fs.Port)
		if err := fs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[HIVE] File Server error: %v", err)
		}
	}()
}

func (fs *FileServer) Stop() {
	if fs.server != nil {
		fs.server.Close()
	}
}
