package panos

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

const defaultLimit uint64 = 1000

func connect(ctx context.Context, d *plugin.QueryData) (interface{}, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "panos"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(interface{}), nil
	}

	// Default to using env vars
	hostname := os.Getenv("PANOS_HOSTNAME")
	apiKey := os.Getenv("PANOS_API_KEY")

	// But prefer the config
	panosConfig := GetConfig(d.Connection)
	if &panosConfig != nil {
		if panosConfig.Hostname != nil {
			hostname = *panosConfig.Hostname
		}
		if panosConfig.APIKey != nil {
			apiKey = *panosConfig.APIKey
		}
	}

	if hostname == "" || apiKey == "" {
		// Credentials not set
		return nil, errors.New("hostname and api_key must be configured")
	}

	conn, err := pango.Connect(
		pango.Client{
			Hostname: hostname,
			Username: *panosConfig.Username,
			Password: *panosConfig.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn, nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "Resource not found")
}
