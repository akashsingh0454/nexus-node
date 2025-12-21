package wazuh

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// Client sends events to a Wazuh Manager (or any Syslog receiver).
type Client struct {
	Protocol string
	Address  string
	conn     net.Conn
}

func NewClient(host string, port int, protocol string) *Client {
	if protocol == "" {
		protocol = "udp"
	}
	return &Client{
		Protocol: protocol,
		Address:  fmt.Sprintf("%s:%d", host, port),
	}
}

// SendEvent formats the event as JSON and sends it over the wire.
// Format: <134>1 TIMESTAMP HOSTNAME TAG - - - JSON
func (c *Client) SendEvent(event interface{}) error {
	// 1. Connect (Ephemeral for UDP, typically persistent for TCP but we keep it simple)
	conn, err := net.Dial(c.Protocol, c.Address)
	if err != nil {
		return fmt.Errorf("wazuh connect failed: %w", err)
	}
	defer conn.Close()

	// 2. Marshal Event
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	// 3. Construct RFC 5424 Syslog Header
	// PRI: <14> (User-level, Info)
	// Timestamp: RFC3339
	// Hostname: os.Hostname
	// AppName: nexus-node
	hostname, _ := os.Hostname()
	timestamp := time.Now().Format(time.RFC3339)

	// We wrap the JSON in a format Wazuh loves: just the JSON if using "json" decoder,
	// or standard syslog. Let's stick to standard syslog with JSON payload.
	// Format: <PRI>TIMESTAMP HOSTNAME APPNAME[PID]: MESSAGE
	msg := fmt.Sprintf("<14>%s %s nexus-node[%d]: %s\n",
		timestamp, hostname, os.Getpid(), string(data))

	// 4. Send
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("wazuh write failed: %w", err)
	}

	return nil
}
