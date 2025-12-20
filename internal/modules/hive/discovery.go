package hive

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

// Peer represents a discovered node in the Hive.
type Peer struct {
	IP       string    `json:"ip"`
	Hostname string    `json:"hostname"`
	Port     int       `json:"port"`
	LastSeen time.Time `json:"-"`
}

// DiscoveryService handles UDP multicast for peer discovery.
type DiscoveryService struct {
	Port      int
	Hostname  string
	StopChan  chan struct{}
	Peers     map[string]*Peer // Key: IP
	NewPeerCh chan *Peer
}

func NewDiscoveryService(port int, hostname string) *DiscoveryService {
	return &DiscoveryService{
		Port:      port,
		Hostname:  hostname,
		StopChan:  make(chan struct{}),
		Peers:     make(map[string]*Peer),
		NewPeerCh: make(chan *Peer, 10),
	}
}

func (ds *DiscoveryService) Start() {
	go ds.listen()
	go ds.broadcast()
}

func (ds *DiscoveryService) Stop() {
	close(ds.StopChan)
}

func (ds *DiscoveryService) broadcast() {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:31337") // Broadcast addr
	if err != nil {
		log.Printf("[HIVE] Failed to resolve broadcast address: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("[HIVE] Failed to dial UDP: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	me := Peer{
		IP:       "", // Receiver determines IP
		Hostname: ds.Hostname,
		Port:     ds.Port,
	}

	for {
		select {
		case <-ds.StopChan:
			return
		case <-ticker.C:
			data, _ := json.Marshal(me)
			conn.Write(data)
		}
	}
}

func (ds *DiscoveryService) listen() {
	addr, err := net.ResolveUDPAddr("udp", ":31337")
	if err != nil {
		log.Printf("[HIVE] Failed to resolve UDP listen address: %v", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("[HIVE] Failed to listen UDP: %v", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ds.StopChan:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			var p Peer
			if err := json.Unmarshal(buf[:n], &p); err != nil {
				continue
			}

			// Don't discover self (simple check, improve later)
			if p.Hostname == ds.Hostname {
				continue
			}

			p.IP = remoteAddr.IP.String()
			p.LastSeen = time.Now()

			ds.handlePeer(p)
		}
	}
}

func (ds *DiscoveryService) handlePeer(p Peer) {
	if existing, exists := ds.Peers[p.IP]; exists {
		existing.LastSeen = time.Now()
	} else {
		log.Printf("[HIVE] Discovered new peer: %s (%s)", p.Hostname, p.IP)
		ds.Peers[p.IP] = &p
		// Non-blocking send
		select {
		case ds.NewPeerCh <- &p:
		default:
		}
	}
}
