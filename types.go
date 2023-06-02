package elasticfacker

import "net/http"

type MockMethods struct {
	StatusCode   int
	Status       string
	BodyAsString string
}
type InMemoryElasticsearch struct {
	indices map[string]map[string]interface{}
	aliases map[string]interface{}
	mock    *MockMethods
	server  *http.Server
}

type IndexFake struct {
	Health string `json:"health"`
	Status string `json:"status"`
	Index  string `json:"index"`
}

type ProductIndexFake struct {
	Aliases map[string]interface{} `json:"aliases"`
}

type IndexMapFake map[string]ProductIndexFake
