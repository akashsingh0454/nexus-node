package container

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

// ContainerInfo represents metadata about a running container.
type ContainerInfo struct {
	ID      string            `json:"Id"`
	Names   []string          `json:"Names"`
	Image   string            `json:"Image"`
	Command string            `json:"Command"`
	Labels  map[string]string `json:"Labels"`
}

// Watcher monitors the container runtime (Docker Engine API).
type Watcher struct {
	SocketPath string
	Client     *http.Client
	StopChan   chan struct{}

	// State
	seenImages map[string]bool
	OnNewImage func(image string, id string)
}

func NewWatcher(socketPath string) *Watcher {
	// Default to standard Docker socket if empty
	if socketPath == "" {
		// Linux default: /var/run/docker.sock
		// Windows default: //./pipe/docker_engine
		socketPath = "/var/run/docker.sock"
	}

	return &Watcher{
		SocketPath: socketPath,
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
		StopChan:   make(chan struct{}),
		seenImages: make(map[string]bool),
	}
}

func (w *Watcher) Start() {
	go w.loop()
}

func (w *Watcher) Stop() {
	close(w.StopChan)
}

func (w *Watcher) loop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.StopChan:
			return
		case <-ticker.C:
			w.scanContainers()
		}
	}
}

func (w *Watcher) scanContainers() {
	resp, err := w.Client.Get("http://localhost/containers/json")
	if err != nil {
		// Silent fail if docker not present, typical in non-container envs
		return
	}
	defer resp.Body.Close()

	var containers []ContainerInfo
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		log.Printf("[CONTAINER] Failed to decode response: %v", err)
		return
	}

	for _, c := range containers {
		// Detect new images
		if !w.seenImages[c.Image] {
			log.Printf("[CONTAINER] New Image Discovered: %s (%s)", c.Image, c.ID[:12])
			w.seenImages[c.Image] = true

			if w.OnNewImage != nil {
				go w.OnNewImage(c.Image, c.ID)
			}
		}
	}
}
