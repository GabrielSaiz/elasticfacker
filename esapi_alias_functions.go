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
		for indexName, index := range es.indicesAlias {
			for aliasNameInIndex, _ := range index {
				if aliasNameInIndex == aliasName {
					indices := IndexMapFake{
						indexName: ProductIndexFake{
							Aliases: map[string]interface{}{
								aliasName: struct{}{},
							},
						},
					}

					jsonData, _ := json.Marshal(indices)

					return &MockMethods{
						StatusCode:   200,
						Status:       "OK",
						BodyAsString: string(jsonData),
					}
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

	index, exists := es.indicesAlias[indexName]
	if exists {
		indices := IndexMapFake{
			indexName: ProductIndexFake{
				Aliases: index,
			},
		}
		jsonData, _ := json.Marshal(indices)

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

func (es *InMemoryElasticsearch) PutAlias(indexName string, aliasName string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	indexExistsResponse := es.IndexExists(indexName)
	if indexExistsResponse.StatusCode != 200 {
		return &MockMethods{
			StatusCode:   500,
			Status:       "Internal Server Error",
			BodyAsString: fmt.Sprintf("{\"error\":\"Index %s does not exist\"}", indexName),
		}
	}

	_, exists := es.aliases[aliasName]
	if exists {
		return &MockMethods{
			StatusCode: 409,
			Status:     "Conflict",
		}
	}

	es.aliases[aliasName] = make(map[string]interface{})
	es.indicesAlias[indexName][aliasName] = es.aliases[aliasName]

	return &MockMethods{
		StatusCode: 200,
		Status:     "OK",
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
			StatusCode:   404,
			Status:       "Not Found",
			BodyAsString: fmt.Sprintf("{\"error\":\"Alias %s does not exist\"}", alias),
		}
	}

	indexAliases, _ := es.indicesAlias[index]

	delete(es.aliases, alias)
	delete(indexAliases, alias)

	return &MockMethods{
		StatusCode: 200,
		Status:     "OK",
	}
}
