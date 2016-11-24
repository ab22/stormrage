package mikrotik

import (
	"fmt"

	routeros "github.com/jda/routeros-api-go"
)

func (s *service) connectToRouter() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.client != nil {
		return nil
	}

	var (
		addr        = fmt.Sprintf("%s:%s", s.cfg.PrivateRouter.Address, s.cfg.PrivateRouter.Port)
		client, err = routeros.New(addr)
	)

	if err != nil {
		return fmt.Errorf("error parsing address: %v", err)
	}

	err = client.Connect(s.cfg.PrivateRouter.User, s.cfg.PrivateRouter.Password)

	if err != nil {
		return fmt.Errorf("error connecting to device: %v", err)
	}

	s.client = client
	return nil
}

func (s *service) closeConn() {
	s.client = nil
}

func (s *service) queryRouter(query string) (*routeros.Reply, error) {
	if err := s.connectToRouter(); err != nil {
		return nil, err
	}

	res, err := s.client.Call(query, nil)
	if err != nil {
		s.closeConn()
		return nil, err
	}

	return &res, nil
}
