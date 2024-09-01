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
	key  string
	doc  *DocJson
}

// NewDocApi responsável por iniciar o processo de configuração.
//
// swaggerURL: URL que irá acessar a página da documentação. Ex.: http://localhost:8080/swagger/
func NewDocApi(swaggerURL string) *StartDocApi {
	if strings.TrimSpace(swaggerURL) == "" {
		log.Fatal("empty url")
	}

	if !strings.Contains(swaggerURL, "http://") && !strings.Contains(swaggerURL, "https://") {
		swaggerURL = "http://" + swaggerURL
	}

	parse, err := url.Parse(swaggerURL)
	if err != nil {
		log.Fatal(err)
	}

	if len(parse.Path) > 0 {
		if string(parse.Path[len(parse.Path)-1]) != "/" {
			parse.Path += "/"
			swaggerURL += "/"
		}
	}

	doc := &StartDocApi{
		url:  swaggerURL,
		path: parse.Path,
		key:  parse.Path + "doc.json",
	}

	doc.doc = Session().NewDoc(doc.key)
	return doc
}

func (d *StartDocApi) Info(title, description, version string) *StartDocApi {
	d.doc.Info.Title = title
	d.doc.Info.Description = description
	d.doc.Info.Version = version
	return d
}

func (d *StartDocApi) Contact(name, url, email string) *StartDocApi {
	d.doc.Info.Contact = &Contact{
		Name:  name,
		URL:   url,
		Email: email,
	}
	return d
}

func (d *StartDocApi) ContactOptions(opts ...OptsContact) *StartDocApi {
	if d.doc.Info.Contact == nil {
		d.doc.Info.Contact = &Contact{}
	}

	for _, fn := range opts {
		fn(d.doc.Info.Contact)
	}

	return d
}

func (d *StartDocApi) License(name, url string) *StartDocApi {
	d.doc.Info.License = &License{
		Name: name, Url: url,
	}
	return d
}

func (d *StartDocApi) ExternalDocs(description, url string) *StartDocApi {
	d.doc.ExternalDocs = &ExternalDocs{
		Description: description,
		Url:         url,
	}
	return d
}

func (d *StartDocApi) Servers(url ...string) *StartDocApi {
	d.doc.AddServer(url...)
	return d
}

// NewRouter para iniciar a configuração de endpoint.
func (d *StartDocApi) NewRouter() *Router {
	return newRouter(d.doc, SecurityNone)
}

// NewRouterSecurityBasic para iniciar a configuração de endpoint com autenticação basic.
func (d *StartDocApi) NewRouterSecurityBasic() *Router {
	ss := NewSecurityShemes(SecurityHttp)
	ss.TypeName = SecurityBasic.String()
	ss.Schema = SecurityBasic.String()
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.doc, SecurityBasic)
}

// NewRouterSecurityBearer para iniciar a configuração de endpoint com autenticação bearer token.
func (d *StartDocApi) NewRouterSecurityBearer() *Router {
	ss := NewSecurityShemes(SecurityHttp)
	ss.TypeName = SecurityBearer.String()
	ss.Schema = SecurityBearer.String()
	ss.Format = "JWT"
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.doc, SecurityBearer)
}

// NewRouterSecurityApiKeyHeader para iniciar a configuração de endpoint com autenticação api key header.
func (d *StartDocApi) NewRouterSecurityApiKeyHeader() *Router {
	return d.newRouterSecurityApiKey(ApiKeyHeader)
}

// NewRouterSecurityApiKeyQuery para iniciar a configuração de endpoint com autenticação api key query.
func (d *StartDocApi) NewRouterSecurityApiKeyQuery() *Router {
	return d.newRouterSecurityApiKey(ApiKeyQuery)
}

func (d *StartDocApi) newRouterSecurityApiKey(in string) *Router {
	ss := NewSecurityShemes(SecurityApiKey)
	ss.In = in
	ss.Name = "apiKey"
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.doc, SecurityApiKey)
}

// NewRouterSecurityOAuth2Password para iniciar a configuração de endpoint com autenticação oauth2 password.
func (d *StartDocApi) NewRouterSecurityOAuth2Password(tokenUrl string) *Router {
	return d.newRouterSecurityOAuth2(tokenUrl, OAuth2Password)
}

// NewRouterSecurityOAuth2Password para iniciar a configuração de endpoint com autenticação oauth2 client.
func (d *StartDocApi) NewRouterSecurityOAuth2Client(tokenUrl string) *Router {
	return d.newRouterSecurityOAuth2(tokenUrl, OAuth2ClientCredentials)
}

func (d *StartDocApi) newRouterSecurityOAuth2(tokenUrl, in string) *Router {
	ss := NewSecurityShemes(SecurityOAuth2)
	ss.Flows = &SecurityFlows{
		in: &SecurityClient{
			TokenUrl: tokenUrl,
			Scopes:   &SecurityScope{},
		},
	}
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.doc, SecurityOAuth2)
}

// HandlerFn responsável por retornar o endereço do swagger e a função do controller.
func (d *StartDocApi) HandlerFunc() (pattern string, controller http.HandlerFunc) {
	slog.Info("DocApi", "URL", d.url)
	return d.path + "*", HandlerFunc(d.url + "doc.json")
}

// HandlerFunc responsável por retornar o endereço do swagger e a função do controller.
//
// Uso no net/http
func (d *StartDocApi) HandlerFuncNetHttp() (pattern string, controller http.HandlerFunc) {
	slog.Info("DocApi", "URL", d.url)
	return d.path, HandlerFunc(d.url + "doc.json")
}
