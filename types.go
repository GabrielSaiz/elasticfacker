package elasticfacker

type MockMethods struct {
	StatusCode   int
	Status       string
	BodyAsString string
}
type InMemoryElasticsearch struct {
	indices map[string]map[string]interface{}
	aliases map[string]string
	mock    *MockMethods
}
