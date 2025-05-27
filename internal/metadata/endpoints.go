package metadata

import (
	"regexp"
	"strings"
)

type AuthLevel int

const (
	AuthNone AuthLevel = iota
	AuthUser
	AuthAdmin
)

type Endpoint struct {
	HTTPMethod  string
	PathPattern string         // Исходный шаблон (например, "/v1/hotels/{id}")
	Regex       *regexp.Regexp // Скомпилированное регулярное выражение
	GRPCMethod  string
	Level       AuthLevel
}

type EndpointConfig struct {
	endpoints      []Endpoint
	swaggerEnabled bool
}

func NewEndpointConfig() *EndpointConfig {
	return &EndpointConfig{
		endpoints:      []Endpoint{},
		swaggerEnabled: true,
	}
}

func (c *EndpointConfig) AddEndpoint(method, pathPattern, grpcMethod string, level AuthLevel) {
	// Replace { with named group pattern and } with the regex part
	regexPattern := "^" + strings.ReplaceAll(pathPattern, "{", "(?P<") + "$"
	regexPattern = strings.ReplaceAll(regexPattern, "}", ">[^/]+)")

	regex := regexp.MustCompile(regexPattern)

	c.endpoints = append(c.endpoints, Endpoint{
		HTTPMethod:  strings.ToUpper(method),
		PathPattern: pathPattern,
		Regex:       regex,
		GRPCMethod:  grpcMethod,
		Level:       level,
	})
}

func (c *EndpointConfig) MatchEndpoint(method, path string) (Endpoint, bool) {
	method = strings.ToUpper(method)
	for _, ep := range c.endpoints {
		if ep.HTTPMethod == method && ep.Regex.MatchString(path) {
			return ep, true
		}
	}
	return Endpoint{}, false
}

func (c *EndpointConfig) SetSwaggerEnabled(enabled bool) {
	c.swaggerEnabled = enabled
}

func (c *EndpointConfig) IsSwaggerEnabled() bool {
	return c.swaggerEnabled
}
