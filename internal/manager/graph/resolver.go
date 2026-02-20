package graph

import (
	"nexus/internal/manager/graph/model"
	"sync"
)

type Resolver struct {
	mu           sync.RWMutex
	Agents       []*model.Agent
	Destinations []*model.DestinationConfig
}
