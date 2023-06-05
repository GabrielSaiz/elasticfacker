package elasticfacker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func (es *InMemoryElasticsearch) SearchTemplate(indexName string, body []byte) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	var searchTemplateRequest ElasticSearchRequest
	err := json.Unmarshal(body, &searchTemplateRequest)
	if err != nil {
		return &MockMethods{
			StatusCode:   400,
			Status:       "Bad Request",
			BodyAsString: fmt.Sprintf("{\"error\":\"%s\"}", err.Error()),
		}
	}

	indexDocuments, exists := es.indicesDocuments[indexName]
	if !exists {
		return &MockMethods{
			StatusCode:   404,
			Status:       "Not Found",
			BodyAsString: fmt.Sprintf("{\"error\":\"Index %s does not exist\"}", indexName),
		}
	}

	var total int
	if len(indexDocuments) > 10 {
		total = 10
	} else {
		total = len(indexDocuments)
	}

	searchTemplateResponse := ElasticSearchResponseFake{
		Took: rand.New(rand.NewSource(time.Now().UnixNano())).Intn(20),
		Shards: ElasticSearchResponseFakeShards{
			Total:      total,
			Successful: 1,
			Skipped:    0,
			Failed:     0,
		},
		Hits: ElasticSearchResponseFakeHits{
			Total: ElasticSearchResponseFakeHitsTotal{
				Value:    len(indexDocuments),
				Relation: "eq",
			},
			Hits: indexDocuments,
		},
	}

	jsonData, _ := json.Marshal(searchTemplateResponse)
	return &MockMethods{
		StatusCode:   200,
		Status:       "OK",
		BodyAsString: string(jsonData),
	}

}
