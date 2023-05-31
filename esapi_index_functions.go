package elasticfacker

import (
	"encoding/json"
	"regexp"
)

func (es *InMemoryElasticsearch) IndexExists(index string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}
	_, exists := es.indices[index]

	var responseStatusCode int
	if exists {
		responseStatusCode = 200
	} else {
		responseStatusCode = 404
	}
	return &MockMethods{
		StatusCode: responseStatusCode,
	}
}

func (es *InMemoryElasticsearch) GetIndex(indexPattern string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	indices := make(map[string]interface{})
	for indexName, index := range es.indices {
		re := regexp.MustCompile(indexPattern)

		if re.MatchString(indexName) {
			indices[indexName] = index
		}
	}

	// Converter tpoJSON
	jsonData, err := json.Marshal(indices)
	if err != nil {
		return &MockMethods{
			StatusCode: 500,
			Status:     "Internal Server Error",
		}
	}

	return &MockMethods{
		StatusCode:   200,
		Status:       "OK",
		BodyAsString: string(jsonData),
	}

}

func (es *InMemoryElasticsearch) CreateIndex(index string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}
	es.indices[index] = make(map[string]interface{})

	return &MockMethods{
		StatusCode: 200,
	}
}

func (es *InMemoryElasticsearch) DeleteIndex(index string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}
	delete(es.indices, index)

	return &MockMethods{
		StatusCode: 200,
	}
}
