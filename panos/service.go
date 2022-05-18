package panos

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (interface{}, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "panos"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData, nil
	}

	// Default to using env vars
	hostname := os.Getenv("PANOS_HOSTNAME")
	username := os.Getenv("PANOS_USERNAME")
	password := os.Getenv("PANOS_PASSWORD")
	apiKey := os.Getenv("PANOS_API_KEY")

	// But prefer the config
	panosConfig := GetConfig(d.Connection)

	if panosConfig.Hostname != nil {
		hostname = *panosConfig.Hostname
	}
	if panosConfig.APIKey != nil {
		apiKey = *panosConfig.APIKey
	}
	if panosConfig.Username != nil {
		username = *panosConfig.Username
	}
	if panosConfig.Password != nil {
		password = *panosConfig.Password
	}

	if len(hostname) == 0 {
		// Credentials not set
		return nil, errors.New("hostname must be configured")
	}

	if !isAPIKeyDefined(apiKey) && !isUsernamePasswordDefined(username, password) {
		// Credentials not set
		return nil, errors.New("either 'api_key' or 'username-password' must be configured")
	}

	if isAPIKeyDefined(apiKey) {
		// if the api key is defined,
		// then choose those over the username/password combo
		username = ""
		password = ""
	}

	requestTimeout := 10
	if panosConfig.Timeout != nil && *panosConfig.Timeout > 0 {
		requestTimeout = *panosConfig.Timeout
	}

	conn, err := pango.Connect(
		pango.Client{
			Hostname: hostname,
			ApiKey:   apiKey,
			Username: username,
			Password: password,
			Timeout:  requestTimeout,
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
	return strings.HasSuffix(err.Error(), "not found")
}

func isAPIKeyDefined(apiKey string) bool {
	return apiKey != ""
}

func isUsernamePasswordDefined(username string, password string) bool {
	return username != "" && password != ""
}
