package docapi

type SecurityType string

const (
	SecurityNone          SecurityType = ""
	SecurityHttp          SecurityType = "http"
	SecurityApiKey        SecurityType = "apiKey"
	SecurityOAuth2        SecurityType = "oauth2"
	SecurityOpenIdConnect SecurityType = "openIdConnect"
	SecurityBasic         SecurityType = "basic"
	SecurityBearer        SecurityType = "bearer"
)

func (s SecurityType) String() string {
	return string(s)
}

const (
	ApiKeyHeader = "header"
	ApiKeyQuery  = "query"
)

const (
	OAuth2ClientCredentials = "clientCredentials"
	OAuth2Password          = "password"
)

type SecuritySchemes struct {
	Type     SecurityType   `json:"type,omitempty"`
	TypeName string         `json:"-"`
	In       string         `json:"in,omitempty"`
	Name     string         `json:"name,omitempty"`
	Schema   string         `json:"scheme,omitempty"`
	Format   string         `json:"bearerFormat,omitempty"`
	Flows    *SecurityFlows `json:"flows,omitempty"`
}

type SecurityFlows map[string]*SecurityClient

type SecurityClient struct {
	TokenUrl string         `json:"tokenUrl,omitempty"`
	Scopes   *SecurityScope `json:"scopes"`
}

type SecurityScope struct{}

func NewSecurityShemes(sType SecurityType) *SecuritySchemes {
	return &SecuritySchemes{Type: sType, TypeName: sType.String()}
}
