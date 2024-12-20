package registry

import (
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/account"
	"github.com/mandarine-io/Backend/internal/domain/service/auth"
	"github.com/mandarine-io/Backend/internal/domain/service/geocoding"
	"github.com/mandarine-io/Backend/internal/domain/service/health"
	masterprofile "github.com/mandarine-io/Backend/internal/domain/service/master/profile"
	"github.com/mandarine-io/Backend/internal/domain/service/resource"
	"github.com/mandarine-io/Backend/internal/domain/service/ws"
	geocoding2 "github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/rs/zerolog/log"
)

type Services struct {
	Account       service.AccountService
	Auth          service.AuthService
	Geocoding     service.GeocodingService
	Health        service.HealthService
	MasterProfile service.MasterProfileService
	Resource      service.ResourceService
	Websocket     service.WebsocketService
}

func setupServices(c *Container) {
	log.Debug().Msg("setup services")

	geocodingProviders := make([]geocoding2.Provider, 0, len(c.Config.GeocodingClients))
	for _, provider := range c.GeocodingProviders {
		geocodingProviders = append(geocodingProviders, provider)
	}

	c.SVCs = &Services{
		Account: account.NewService(
			c.Repos.User,
			c.Cache.Manager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Auth: auth.NewService(
			c.Repos.User,
			c.Repos.BannedToken,
			c.OauthProviders,
			c.Cache.Manager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Geocoding: geocoding.NewService(
			geocodingProviders,
			c.Cache.Manager,
		),
		Health: health.NewService(
			c.DB,
			c.Cache.RDB,
			c.PubSub.RDB,
			c.S3.Minio,
			c.SmtpSender,
		),
		MasterProfile: masterprofile.NewService(
			c.Repos.MasterProfile,
		),
		Resource: resource.NewService(
			c.S3.Client,
		),
		Websocket: ws.NewService(
			c.WebsocketPool,
		),
	}
}
