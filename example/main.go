package main

import (
	"docapi"
	"encoding/json"
	"net/http"

	"docapi/openapi"

	"github.com/go-chi/chi/v5"
)

type Model struct {
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
			ParamQuery("id", openapi.DataTypeInteger, true).
			ResponseObjectBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			HandlerFn(),
	)

	//Auth JWT
	router = doc.NewRouterSecurityBearer()
	r.MethodFunc(
		router.Get("/get-bearer", controller).
			Tag("Generic").
			Description("Method Get").
			ResponseObjectBodyJson(http.StatusOK, http.StatusText(http.StatusOK), Model{}).
			Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).
			HandlerFn(),
	)

	// Endpoint swagger
	r.HandleFunc(doc.HandlerFn())
	http.ListenAndServe(":8080", r)
}
