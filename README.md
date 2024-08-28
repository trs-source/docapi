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

	"github.com/trs-source/docapi"

	"github.com/go-chi/chi/v5"
)

type Model struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Model2 []Model2 `json:"model2"`
}

type Model2 struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(Model{
		ID:   1,
		Name: "Name",
	})

}

func main() {
	doc := docapi.NewDocApi("http://localhost:8080/swagger/")
	doc.Info("DocApi swagger documentation", "Lib docapi", "1.0").
		Contact("Test", "email@email.com.br").
		ExternalDocs("Help", "https://test.com/").
		License("Test", "https://www.test.com.br").
		Servers("http://localhost:8080")

	r := chi.NewRouter()

	//No auth
	router := doc.NewRouter()
	r.MethodFunc(
		router.Get("/get", controller).
			Tag("Generic").
			Description("Method Get").
			Summary("Summary method get").
			ParamQuery("id", docapi.DataTypeInteger, true).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), []Model2{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			HandlerFn(),
	)

	//Auth JWT
	router = doc.NewRouterSecurityBearer()
	r.MethodFunc(
		router.Get("/get-bearer", controller).
			Tag("Generic").
			Description("Method Get").
			ResponseBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			HandlerFn(),
	)

	// Endpoint swagger
	r.HandleFunc(doc.HandlerFn())
	http.ListenAndServe(":8080", r)
}


```

