package docapi

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type StartDocApi struct {
	url  string
	path string
	doc  *Doc
}

// NewDocApi responsável por iniciar o processo de configuração.
//
// docURL: URL que irá acessar a página da docação. Ex.: http://localhost:8080/swagger/
func NewDocApi(docURL string) *StartDocApi {
	if strings.TrimSpace(docURL) == "" {
		log.Fatal("empty url")
	}

	if !strings.Contains(docURL, "http://") && !strings.Contains(docURL, "https://") {
		docURL = "http://" + docURL
	}

	parse, err := url.Parse(docURL)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasSuffix(parse.Path, "/") {
		parse.Path += "/"
		docURL += "/"
	}

	swagger := &StartDocApi{
		url:  docURL,
		path: parse.Path,
	}

	swagger.doc = GetDocs().NewDoc(parse.Path + "doc.json")
	return swagger
}

func (s *StartDocApi) Info(title, description, version string) *StartDocApi {
	s.doc.Info.Title = title
	s.doc.Info.Description = description
	s.doc.Info.Version = version
	return s
}

func (s *StartDocApi) Contact(name string, opts ...OptsContact) *StartDocApi {
	c := &Contact{Name: name}
	for _, fn := range opts {
		fn(c)
	}
	s.doc.Info.Contact = c
	return s
}

func (s *StartDocApi) License(name, url string) *StartDocApi {
	s.doc.Info.License = &License{
		Name: name, Url: url,
	}
	return s
}

func (s *StartDocApi) ExternalDocs(description, helpURL string) *StartDocApi {
	s.doc.ExternalDocs = &ExternalDocs{
		Description: description,
		URL:         helpURL,
	}
	return s
}

func (s *StartDocApi) Server(url ...string) *StartDocApi {
	s.doc.AddServer(url...)
	return s
}

// NewRouter para iniciar a configuração de endpoint.
func (s *StartDocApi) NewRouter() Router {
	return newRouter(s.doc, SecurityNone)
}

// NewRouterSecurityBasic para iniciar a configuração de endpoint com autenticação basic.
func (s *StartDocApi) NewRouterSecurityBasic() Router {
	ss := NewSecurityShemes(SecurityHttp)
	ss.TypeName = SecurityBasic.String()
	ss.Schema = SecurityBasic.String()
	s.doc.Components.AddSecurity(ss)
	return newRouter(s.doc, SecurityBasic)
}

// NewRouterSecurityBearer para iniciar a configuração de endpoint com autenticação bearer token.
func (s *StartDocApi) NewRouterSecurityBearer() Router {
	ss := NewSecurityShemes(SecurityHttp)
	ss.TypeName = SecurityBearer.String()
	ss.Schema = SecurityBearer.String()
	ss.Format = "JWT"
	s.doc.Components.AddSecurity(ss)
	return newRouter(s.doc, SecurityBearer)
}

// NewRouterSecurityApiKeyHeader para iniciar a configuração de endpoint com autenticação api key header.
func (s *StartDocApi) NewRouterSecurityApiKeyHeader(key string) Router {
	return s.newRouterSecurityApiKey(key, ApiKeyHeader)
}

// NewRouterSecurityApiKeyQuery para iniciar a configuração de endpoint com autenticação api key query.
func (s *StartDocApi) NewRouterSecurityApiKeyQuery(key string) Router {
	return s.newRouterSecurityApiKey(key, ApiKeyQuery)
}

func (s *StartDocApi) newRouterSecurityApiKey(key, in string) Router {
	ss := NewSecurityShemes(SecurityApiKey)
	ss.In = in
	ss.Name = key
	if key == "" {
		ss.Name = "apiKey"
	}
	s.doc.Components.AddSecurity(ss)
	return newRouter(s.doc, SecurityApiKey)
}

// NewRouterSecurityOAuth2Password para iniciar a configuração de endpoint com autenticação oauth2 passwors.
func (s *StartDocApi) NewRouterSecurityOAuth2Password(tokenUrl string) Router {
	return s.newRouterSecurityOAuth2(tokenUrl, OAuth2Password)
}

// NewRouterSecurityOAuth2Password para iniciar a configuração de endpoint com autenticação oauth2 client.
func (s *StartDocApi) NewRouterSecurityOAuth2Client(tokenUrl string) Router {
	return s.newRouterSecurityOAuth2(tokenUrl, OAuth2ClientCredentials)
}

func (s *StartDocApi) newRouterSecurityOAuth2(tokenUrl, in string) Router {
	ss := NewSecurityShemes(SecurityOAuth2)
	ss.Flows = &SecurityFlows{
		in: &SecurityClient{
			TokenUrl: tokenUrl,
			Scopes:   &SecurityScope{},
		},
	}
	s.doc.Components.AddSecurity(ss)
	return newRouter(s.doc, SecurityOAuth2)
}

// HandlerFn responsável por retornar o endereço do swagger e a função do controller.
func (s *StartDocApi) HandlerFunc() (pattern string, controller http.HandlerFunc) {
	slog.Info("DocApi", "URL", s.url)
	return s.path + "*", HandlerFunc(s.url + "doc.json")
}

// HandlerFunc responsável por retornar o endereço do swagger e a função do controller.
//
// Uso no net/http
func (s *StartDocApi) HandlerFuncNetHttp() (pattern string, controller http.HandlerFunc) {
	slog.Info("DocApi", "URL", s.url)
	return s.path, HandlerFunc(s.url + "doc.json")
}
