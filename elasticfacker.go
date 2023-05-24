package elasticfacker

import (
	"log"
	"net/http"
	"strings"
)

func NewInMemoryElasticsearch(mock *MockMethods) *InMemoryElasticsearch {
	return &InMemoryElasticsearch{
		indices: make(map[string]map[string]interface{}),
		aliases: make(map[string]string),
		mock:    mock,
		server:  &http.Server{},
	}
}

func (es *InMemoryElasticsearch) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", es.handleRequest)

	es.server = &http.Server{
		Addr:    ":9200",
		Handler: mux,
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

func (es *InMemoryElasticsearch) handleRequest(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	indexOrAlias := splitPath[1]
	var response *MockMethods
	switch {
	case r.Method == http.MethodHead && len(splitPath) == 2:
		// IndicesExistsRequest
		response = es.IndexExists(indexOrAlias)
	case r.Method == http.MethodPut && len(splitPath) == 2:
		// IndicesCreateRequest
		response = es.CreateIndex(indexOrAlias)
	case r.Method == http.MethodGet && len(splitPath) == 3 && splitPath[2] == "_alias":
		// IndicesGetAliasRequest (by index)
		response = es.GetAlias(indexOrAlias)
	case r.Method == http.MethodGet && len(splitPath) == 3 && splitPath[1] == "_alias":
		// IndicesGetAliasRequest (by alias)
		response = es.GetAlias(splitPath[2])
	case r.Method == http.MethodGet && len(splitPath) == 4 && splitPath[2] == "_cat" && splitPath[3] == "indices":
		// CatIndicesRequest
		response = es.IndexExists(indexOrAlias)
	case r.Method == http.MethodDelete && len(splitPath) == 2:
		// IndicesDeleteRequest
		response = es.DeleteIndex(indexOrAlias)
	case r.Method == http.MethodDelete && len(splitPath) == 4 && splitPath[2] == "_aliases":
		// IndicesDeleteAliasRequest
		response = es.DeleteAlias(indexOrAlias, splitPath[3])
	case r.Method == http.MethodPut && len(splitPath) == 4 && splitPath[2] == "_aliases":
		// IndicesPutAliasRequest
		response = es.PutAlias(indexOrAlias, splitPath[3])
	default:
		response = &MockMethods{
			StatusCode:   400,
			BodyAsString: `{"error": "Invalid request"}`,
		}
	}

	w.WriteHeader(response.StatusCode)
	_, _ = w.Write([]byte(response.BodyAsString))
}
