package grpc_gateway

type EndpointConfig struct {
	PublicPaths    map[string]struct{}
	AdminPrefixes  []string
	SwaggerEnabled bool
}

func NewEndpointConfig() *EndpointConfig {
	return &EndpointConfig{
		PublicPaths:    make(map[string]struct{}),
		AdminPrefixes:  []string{"/v1/admin"},
		SwaggerEnabled: true,
	}
}

func (c *EndpointConfig) AddPublic(path string) {
	c.PublicPaths[path] = struct{}{}
}

func (c *EndpointConfig) AddAdminPrefix(prefix string) {
	c.AdminPrefixes = append(c.AdminPrefixes, prefix)
}
