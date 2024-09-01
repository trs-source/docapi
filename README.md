# DocApi

<p align="justify"> DocApi is a library that facilitates API documentation 📑 for the GO language </p>

### Install
 - go get -u github.com/trs-source/docapi

### Example:
```
package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trs-source/docapi"
)

type Model struct {
	ID     int64    `json:"id" docapi:"examples:1"`
	Name   string   `json:"name"`
	Model2 []Model2 `json:"model2"`
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
		Contact("Test", "https://www.example.com/support", "email@email.com.br").
		ExternalDocs("Help", "https://test.com/").
		License("Test", "https://www.test.com.br").
		Servers("http://localhost:8080")

	r := chi.NewRouter()

	//No auth
	router := doc.NewRouter()
	r.MethodFunc(
		router.Get("/get", controllerGet).
			Tag("Generic").
			Description("Method Get").
			Summary("Summary method get").
			ParamQuery("id", docapi.DataTypeInteger, true).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), []Model2{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			MethodFunc(),
	)

	//Auth JWT
	router = doc.NewRouterSecurityBearer()
	r.MethodFunc(
		router.Post("/post", controllerPost).
			Tag("Generic").
			Description("Method Post").
			RequestBodyJson(Model{}, docapi.WithDescription("Request POST"), docapi.WithRequired(true)).
			Response(http.StatusCreated, http.StatusText(http.StatusCreated)).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			Response(http.StatusForbidden, http.StatusText(http.StatusForbidden)).
			MethodFunc(),
	)

	//Endpoint swagger
	r.Get(doc.HandlerFunc())
	http.ListenAndServe(":8080", r)
}

```

