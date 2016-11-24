package mikrotik

import "github.com/ab22/stormrage/models"

func (s *service) RequestClients() ([]models.Client, error) {
	res, err := s.queryRouter("/queue/simple/print")
	if err != nil {
		return nil, err
	}

	clients := make([]models.Client, 0, len(res.SubPairs))

	for _, pair := range res.SubPairs {
		clients = append(clients, models.Client{
			ID:             pair["id"],
			Name:           pair["name"],
			Target:         pair["target"],
			MaxLimit:       pair["max-limit"],
			BurstLimit:     pair["burst-limit"],
			BurstThreshold: pair["burst-threshold"],
			BurstTime:      pair["burst-time"],
		})
	}

	return clients, nil
}
