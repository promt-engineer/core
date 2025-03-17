package ip2country

import (
	"errors"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"io"
	"net/http"
	"strings"
	"time"
)

const BaseURL = "https://ipinfo.io/"
const EndURL = "/country"
const IP2CountryName = "IP2Country"

var (
	ErrWrongStatusCode = errors.New("wrong status code")
)

type ClientWithCache struct {
	cacheLifetime    time.Duration
	cache            *ttlcache.Cache[string, string]
	countryTransform func(string) string
	httpClient       *http.Client
}

func NewClientWithCache(cacheLifetime time.Duration, countryTransform func(string) string) *ClientWithCache {
	return &ClientWithCache{
		cacheLifetime:    cacheLifetime,
		countryTransform: countryTransform,
		cache: ttlcache.New[string, string](
			ttlcache.WithTTL[string, string](cacheLifetime),
		),

		httpClient: http.DefaultClient,
	}
}

func (c *ClientWithCache) Get(ip string) (string, error) {
	result := c.cache.Get(ip)

	var country string

	if result == nil {
		country, err := c.askServer(ip)

		if err != nil {
			return "", err
		}

		country = strings.TrimSpace(country)
		country = c.countryTransform(country)

		c.cache.Set(ip, country, ttlcache.DefaultTTL)
	} else {
		country = result.Value()
	}

	return country, nil
}

func (c *ClientWithCache) askServer(ip string) (string, error) {
	resp, err := c.httpClient.Get(BaseURL + ip + EndURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v, %v", ErrWrongStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *ClientWithCache) Start() {
	go c.cache.Start()
}

func (c *ClientWithCache) Stop() {
	c.cache.Stop()
}
