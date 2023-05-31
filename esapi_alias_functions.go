package elasticfacker

import (
	"encoding/json"
	"fmt"
)

func (es *InMemoryElasticsearch) GetAlias(aliasName string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	_, exists := es.aliases[aliasName]
	if exists {
		for index, alias := range es.aliases {
			if alias == aliasName {
				jsonData, err := json.Marshal(index)
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
		}
	}

	return &MockMethods{
		StatusCode: 404,
		Status:     "Not Found",
	}
}

func (es *InMemoryElasticsearch) GetAliasFromIndex(indexName string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	index, exists := es.indices[indexName]
	if exists {
		jsonData, err := json.Marshal(index)
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

	return &MockMethods{
		StatusCode: 404,
		Status:     "Not Found",
	}
}

func (es *InMemoryElasticsearch) PutAlias(index string, alias string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	indexExistsResponse := es.IndexExists(index)
	if indexExistsResponse.StatusCode != 200 {
		return &MockMethods{
			StatusCode:   500,
			Status:       "Internal Server Error",
			BodyAsString: fmt.Sprintf("{\"error\":\"Index %s does not exist\"}", index),
		}
	}

	es.aliases[alias] = make(map[string]interface{})
	es.indices[index][alias] = es.aliases[alias]

	return &MockMethods{
		StatusCode: 200,
	}
}

func (es *InMemoryElasticsearch) DeleteAlias(index string, alias string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	indexExistsResponse := es.IndexExists(index)
	if indexExistsResponse.StatusCode != 200 {
		return &MockMethods{
			StatusCode:   500,
			Status:       "Internal Server Error",
			BodyAsString: fmt.Sprintf("{\"error\":\"Index %s does not exist\"}", index),
		}
	}

	aliasExistsResponse := es.GetAlias(alias)
	if aliasExistsResponse.StatusCode != 200 {
		return &MockMethods{
			StatusCode:   500,
			Status:       "Internal Server Error",
			BodyAsString: fmt.Sprintf("{\"error\":\"Alias %s does not exist\"}", alias),
		}
	}

	indexAliases, _ := es.indices[index]
	_, existAliasInIndex := indexAliases[alias]
	if !existAliasInIndex {
		return &MockMethods{
			StatusCode:   404,
			Status:       "Not Found",
			BodyAsString: fmt.Sprintf("{\"error\":\"Alias %s does not exist assicated to index %s\"}", alias, index),
		}
	}

	delete(es.aliases, alias)
	delete(indexAliases, alias)

	return &MockMethods{
		StatusCode: 200,
		Status:     "OK",
	}
}
