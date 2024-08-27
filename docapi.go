package docapi

import (
	"docapi/internal"
	"docapi/openapi"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type DocApi struct {
	url     string
	pattern string
	key     string
	doc     *openapi.Doc
}

func NewDocApi(url string) *DocApi {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("Problem in DocApi", "error", err)
		}
	}()

	if strings.TrimSpace(url) == "" {
		log.Fatal("empty url")
	}

	d := &DocApi{url: url}

	count := 0
	for _, v := range d.url {
		if string(v) == "/" {
			count++
		}
		if count > 2 {
			d.pattern += string(v)
		}
	}

	d.key = d.pattern + "doc.json"
	d.pattern += "*"

	d.doc = openapi.Mapping().NewDoc(d.key)
	return d
}

func (d *DocApi) Info(title, description, version string) *DocApi {
	d.doc.Info.Title = title
	d.doc.Info.Description = description
	d.doc.Info.Version = version
	return d
}

func (d *DocApi) Contact(name, email string) *DocApi {
	d.doc.Info.Contact = &openapi.Contact{
		Name:  name,
		Email: email,
	}
	return d
}

func (d *DocApi) License(name, url string) *DocApi {
	d.doc.Info.License = &openapi.License{
		Name: name, Url: url,
	}
	return d
}

func (d *DocApi) ExternalDocs(description, url string) *DocApi {
	d.doc.ExternalDocs = &openapi.ExternalDocs{
		Description: description,
		Url:         url,
	}
	return d
}

func (d *DocApi) Servers(url ...string) *DocApi {
	d.doc.AddServer(url...)
	return d
}

func (d *DocApi) NewRouter() *Router {
	return newRouter(d.key, openapi.SecurityNone)
}

func (d *DocApi) NewRouterSecurityBasic() *Router {
	ss := openapi.NewSecurityShemes(openapi.SecurityHttp)
	ss.TypeName = openapi.SecurityBasic.String()
	ss.Schema = openapi.SecurityBasic.String()
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.key, openapi.SecurityBasic)
}

func (d *DocApi) NewRouterSecurityBearer() *Router {
	ss := openapi.NewSecurityShemes(openapi.SecurityHttp)
	ss.TypeName = openapi.SecurityBearer.String()
	ss.Schema = openapi.SecurityBearer.String()
	ss.Format = "JWT"
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.key, openapi.SecurityBearer)

}

func (d *DocApi) NewRouterSecurityApiKeyHeader() *Router {
	return d.newRouterSecurityApiKey(openapi.ApiKeyHeader)
}

func (d *DocApi) NewRouterSecurityApiKeyQuery() *Router {
	return d.newRouterSecurityApiKey(openapi.ApiKeyQuery)
}

func (d *DocApi) newRouterSecurityApiKey(in string) *Router {
	ss := openapi.NewSecurityShemes(openapi.SecurityApiKey)
	ss.In = in
	ss.Name = "apiKey"
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.key, openapi.SecurityApiKey)
}

func (d *DocApi) NewRouterSecurityOAuth2Password(tokenUrl string) *Router {
	return d.newRouterSecurityOAuth2(tokenUrl, openapi.OAuth2Password)
}

func (d *DocApi) NewRouterSecurityOAuth2Client(tokenUrl string) *Router {
	return d.newRouterSecurityOAuth2(tokenUrl, openapi.OAuth2ClientCredentials)
}

func (d *DocApi) newRouterSecurityOAuth2(tokenUrl, in string) *Router {
	ss := openapi.NewSecurityShemes(openapi.SecurityOAuth2)
	ss.Flows = &openapi.SecurityFlows{
		in: &openapi.SecurityClient{
			TokenUrl: tokenUrl,
			Scopes:   &openapi.SecurityScope{},
		},
	}
	d.doc.AddComponentSecurity(ss)
	return newRouter(d.key, openapi.SecurityOAuth2)
}

// HandlerFn responsável por retornar o endereço do swagger e a função do controller.
func (d *DocApi) HandlerFn() (pattern string, controller http.HandlerFunc) {
	slog.Info("DocApi", "URL", d.url)
	return d.pattern, internal.Handle(d.url + "doc.json")
}
