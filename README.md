# DocApi

<p align="justify"> DocApi is a library that facilitates API documentation 📑 for the GO language </p>

### Install
 - go get -u github.com/trs-source/docapi

### Example:
```go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/trs-source/docapi"
)

type Model struct {
	ID     int64    `json:"id" docapi:"example:1"`
	Name   string   `json:"name" docapi:"example:model;required:true"`
	Model2 []Model2 `json:"model2"`
	Type   []int8   `json:"type" docapi:"enum:1,2,3,4;example:1;required:true"`
}

type Model2 struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func controllerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(Model{
		ID:   1,
		Name: "Name",
	})
}

func controllerPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func main() {
	doc := docapi.NewDocApi("http://localhost:8080/swagger/")
	doc.Info("DocApi swagger documentation", "Lib docapi", "1.0").
		Contact("Test", docapi.WithContactEmail("email@email.com.br"), docapi.WithContactWebSite("https://www.example.com/support")).
		ExternalDocs("Help", "https://test.com/").
		License("Test", "https://www.test.com.br").
		Server("http://localhost:8080")

	r := http.NewServeMux()

	//No auth
	router := doc.NewRouter()
	r.HandleFunc(
		router.Get("/get", controllerGet).
			Tag("Generic").
			Description("Method Get").
			Summary("Summary method get").
			ParamQuery("id", docapi.DataTypeInteger, docapi.WithParamRequired()).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), []Model2{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			HandleFunc(),
	)

	//router = doc.NewRouterSecurityApiKeyHeader()
	//router = doc.NewRouterSecurityApiKeyQuery()
	//router = doc.NewRouterSecurityBasic()
	//router = doc.NewRouterSecurityOAuth2Client("http://localhost:8080/auth")
	//router = doc.NewRouterSecurityOAuth2Password("http://localhost:8080/auth")

	//Auth JWT
	router = doc.NewRouterSecurityBearer()
	r.HandleFunc(
		router.Post("/post", controllerPost).
			Tag("Generic").
			Description("Method Post").
			RequestBodyJson(Model{}, docapi.WithDescription("Request POST"), docapi.WithRequired()).
			Response(http.StatusCreated, http.StatusText(http.StatusCreated)).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			Response(http.StatusForbidden, http.StatusText(http.StatusForbidden)).
			HandleFunc(),
	)

	//Endpoint swagger
	r.HandleFunc(doc.HandlerFuncNetHttp())
	http.ListenAndServe(":8080", r)
}


```

