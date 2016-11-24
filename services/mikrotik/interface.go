package mikrotik

import (
	"sync"

	"github.com/ab22/stormrage/config"
	"github.com/ab22/stormrage/models"
	routeros "github.com/jda/routeros-api-go"
)

// Service interface describes all functions that must be implemented.
type Service interface {
	RequestClients() ([]models.Client, error)
}

type service struct {
	cfg    *config.Config
	client *routeros.Client
	mutex  sync.Mutex
}

// NewService initialization.
func NewService(cfg *config.Config) Service {
	s := &service{
		client: nil,
		cfg:    cfg,
		mutex:  sync.Mutex{},
	}

	s.connectToRouter()

	return s
}
