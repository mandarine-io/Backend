package registry

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/chained"
	"github.com/mandarine-io/Backend/pkg/geocoding/factory"
	"github.com/rs/zerolog/log"
)

func setupGeocodingClients(c *Container) {
	log.Debug().Msg("setup geocoding providers")

	providers := make([]geocoding.Provider, 0, len(c.Config.GeocodingClients))
	c.GeocodingProviders = make(map[string]geocoding.Provider)

	for k, v := range c.Config.GeocodingClients {
		log.Debug().Msgf("setup geocoding provider: %s", k)
		provider := factory.NewProviderByKey(k, v.APIKey)
		if provider == nil {
			log.Warn().Msgf("unknown geocoding provider: %s", k)
			continue
		}

		c.GeocodingProviders[k] = provider
		providers = append(providers, provider)
	}

	c.GeocodingChainedProvider = chained.NewProvider(providers...)
}