package hive

import (
	"context"
	"log"
	"nexus/internal/config"
)

// Manager orchestrates the Hive P2P subsystem.
type Manager struct {
	Config    config.HiveConfig
	Discovery *DiscoveryService
	Server    *FileServer
	HostName  string
}

func NewManager(cfg config.HiveConfig, hostname string) *Manager {
	return &Manager{
		Config:    cfg,
		HostName:  hostname,
		Discovery: NewDiscoveryService(cfg.Port, hostname),
		Server:    NewFileServer(cfg.Port, cfg.CacheDir),
	}
}

func (m *Manager) Start(ctx context.Context) {
	if !m.Config.Enabled {
		return
	}
	log.Println("[HIVE] Starting Peer-to-Peer subsystem...")
	m.Server.Start()
	m.Discovery.Start()
}

func (m *Manager) Stop() {
	if !m.Config.Enabled {
		return
	}
	m.Discovery.Stop()
	m.Server.Stop()
}

// GetPeers returns a list of active peers.
func (m *Manager) GetPeers() []Peer {
	peers := make([]Peer, 0, len(m.Discovery.Peers))
	for _, p := range m.Discovery.Peers {
		peers = append(peers, *p)
	}
	return peers
}
