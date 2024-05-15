package service

import (
	"github.com/go-chi/chi"
  	"github.com/zepif/Test-service/internal/config"
	"github.com/zepif/Test-service/internal/data/pg"
	"github.com/zepif/EtherUSDC/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
            handlers.CtxDB(pg.NewStorage(cfg.DB())),
		),
	)
	r.Route("/integrations/EtherUSDC", func(r chi.Router) {
		r.Post()
        r.Get()
        r.Get()
	})

	return r
}
