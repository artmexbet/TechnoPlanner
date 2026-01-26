package publisher

import (
	"context"
	"time"

	"github.com/artmexbet/TechnoPlanner/services/gateway/internal/domain"
)

// вот тут наверное можно сделать отдельный стракт для того чтобы весь Poster не пихать
// если я вообще нормально сделал, не уверен
// type TechCaller struct {
// 	Publisher *Publisher
// }

func (s *Publisher) AddTechnic(ctx context.Context, technic domain.Technic) (domain.Technic, error) {
	jsoned, err := technic.MarshalJSON()
	if err != nil {
		return domain.Technic{}, err
	}
	msg, err := s.nc.Request("tech.add", jsoned, time.Second)

	var respondTech domain.Technic
	err = respondTech.UnmarshalJSON(msg.Data)
	if err != nil {
		return domain.Technic{}, err
	}

	return respondTech, nil
}
