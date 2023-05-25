package elasticfacker

import "fmt"

func (es *InMemoryElasticsearch) GetAlias(alias string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	_, exists := es.aliases[alias]
	var responseStatusCode int
	if exists {
		responseStatusCode = 200
	} else {
		responseStatusCode = 404
	}

	var bodyResponse string
	for index, alias := range es.aliases {
		if alias == alias {
			bodyResponse = index
		}
	}
	return &MockMethods{
		StatusCode:   responseStatusCode,
		BodyAsString: bodyResponse,
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
