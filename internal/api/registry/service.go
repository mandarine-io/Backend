package registry

import (
	"github.com/mandarine-io/Backend/internal/api/service/account"
	"github.com/mandarine-io/Backend/internal/api/service/auth"
	"github.com/mandarine-io/Backend/internal/api/service/resource"
	"github.com/mandarine-io/Backend/internal/api/service/ws"
	"github.com/rs/zerolog/log"
)

type Services struct {
	Auth     *auth.Service
	Account  *account.Service
	Resource *resource.Service
	WS       *ws.Service
}

func newServices(c *Container) *Services {
	log.Debug().Msg("setup services")
	return &Services{
		Auth: auth.NewService(
			c.Repositories.User,
			c.Repositories.BannedToken,
			c.OauthProviders,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Account: account.NewService(
			c.Repositories.User,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Resource: resource.NewService(c.S3Client),
		WS:       ws.NewService(c.WebsocketPool),
	}
}
