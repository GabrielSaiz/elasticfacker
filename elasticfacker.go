package elasticfacker

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	HeaderContentType     = "Content-Type"
	HeaderXElasticProduct = "X-Elastic-Product"
)

func NewInMemoryElasticsearch() *InMemoryElasticsearch {
	return &InMemoryElasticsearch{
		indices: make(map[string]map[string]interface{}),
		aliases: make(map[string]interface{}),
		//mock:    mock,
		//server:  &http.Server{},
	}
}

func (es *InMemoryElasticsearch) SetMockMethods(mock *MockMethods) {
	es.mock = mock
}

func (es *InMemoryElasticsearch) Start(address string) {
	go es.startServer(address)
}

func (es *InMemoryElasticsearch) startServer(address string) {
	r := mux.NewRouter()
	r.HandleFunc("/", es.handleRoot).Methods("GET")
	r.HandleFunc("/{indexName}", es.handleIndicesExists).Methods("HEAD")                             //esapi.IndicesExistsRequest
	r.HandleFunc("/{indexName}", es.handleIndicesCreate).Methods("PUT")                              //esapi.IndicesCreateRequest
	r.HandleFunc("/_cat/indices/{indexNamePattern}", es.handleCatIndices).Methods("GET")             //esapi.CatIndicesRequest
	r.HandleFunc("/{indexName}/_alias", es.handleIndicesGetAlias).Methods("GET")                     //esapi.IndicesGetAliasRequest
	r.HandleFunc("/{indexName}", es.handleIndicesDelete).Methods("DELETE")                           //esapi.IndicesDeleteRequest
	r.HandleFunc("/_alias/{aliasName}", es.handleIndicesGetAlias).Methods("GET")                     //esapi.IndicesGetAliasRequest
	r.HandleFunc("/{indexName}/_aliases/{aliasName}", es.handleIndicesDeleteAlias).Methods("DELETE") //esapi.IndicesDeleteAliasRequest
	r.HandleFunc("/{indexName}/_aliases/{aliasName}", es.handleIndicesPutAlias).Methods("PUT")       //esapi.IndicesPutAliasRequest

	es.server = &http.Server{
		Addr:    address,
		Handler: r,
	}

	log.Println("Starting the server on port 9200")

	err := es.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func (es *InMemoryElasticsearch) Stop() {
	if es.server != nil {
		es.server.Close()
	}
}

func (es *InMemoryElasticsearch) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, "application/json")
	w.Header().Set(HeaderXElasticProduct, "Elasticsearch")
	json.NewEncoder(w).Encode(map[string]string{
		"name":         "elasticsearch-simulator",
		"version":      "8.0.0",
		"cluster_name": "elasticsearch-simulator-cluster",
		"cluster_uuid": "fKg7K_YTQH6pG5-VzF7nZQ",
		"tagline":      "You Know, for Search",
	})
}

func (es *InMemoryElasticsearch) handleIndicesExists(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	response := es.IndexExists(indexName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleIndicesCreate(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	response := es.CreateIndex(indexName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleCatIndices(w http.ResponseWriter, r *http.Request) {
	indexNamePattern := mux.Vars(r)["indexNamePattern"]
	response := es.IndexExists(indexNamePattern)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleIndicesGetAlias(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	response := es.GetAlias(indexName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleIndicesDelete(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	response := es.DeleteIndex(indexName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleIndicesDeleteAlias(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	aliasName := mux.Vars(r)["aliasName"]
	response := es.DeleteAlias(indexName, aliasName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) handleIndicesPutAlias(w http.ResponseWriter, r *http.Request) {
	indexName := mux.Vars(r)["indexName"]
	aliasName := mux.Vars(r)["aliasName"]
	response := es.PutAlias(indexName, aliasName)
	es.writeResponse(w, response)
}

func (es *InMemoryElasticsearch) writeResponse(w http.ResponseWriter, response *MockMethods) {
	w.Header().Set(HeaderXElasticProduct, "Elasticsearch")
	w.WriteHeader(response.StatusCode)
	_, _ = w.Write([]byte(response.BodyAsString))
}
