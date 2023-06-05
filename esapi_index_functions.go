package elasticfacker

import (
	"encoding/json"
	"regexp"
)

func (es *InMemoryElasticsearch) IndexExists(index string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}
	_, exists := es.indicesAlias[index]

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
	indicesStruct := make([]IndexFake, 0)
	for indexName, index := range es.indicesAlias {
		re := regexp.MustCompile(indexPattern)

		if re.MatchString(indexName) {
			indices[indexName] = index
			indicesStruct = append(indicesStruct, IndexFake{
				Health: "yellow",
				Status: "open",
				Index:  indexName,
			})
		}
	}

	if len(indicesStruct) == 0 {
		return &MockMethods{
			StatusCode: 404,
			Status:     "Not Found",
		}
	}

	// Converter toJSON
	jsonData, _ := json.Marshal(indicesStruct)

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

	_, exists := es.indicesAlias[index]

	var responseStatusCode int
	var responseStatus string
	if exists {
		responseStatusCode = 409
		responseStatus = "Conflict"
	} else {
		es.indicesAlias[index] = make(map[string]interface{})
		es.indicesDocuments[index] = make([]Document, 0)

		responseStatusCode = 200
		responseStatus = "OK"
	}
	return &MockMethods{
		StatusCode: responseStatusCode,
		Status:     responseStatus,
	}
}

func (es *InMemoryElasticsearch) DeleteIndex(index string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	_, exists := es.indicesAlias[index]
	if !exists {
		return &MockMethods{
			StatusCode: 404,
			Status:     "Not Found",
		}
	}

	delete(es.indicesAlias, index)
	delete(es.indicesDocuments, index)

	return &MockMethods{
		StatusCode: 200,
		Status:     "OK",
	}
}
