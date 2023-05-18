package elasticfacker

func (es *InMemoryElasticsearch) GetAlias(indexOrAlias string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	_, exists := es.aliases[indexOrAlias]
	var responseStatusCode int
	if exists {
		responseStatusCode = 200
	} else {
		responseStatusCode = 404
	}

	var bodyResponse string
	for index, alias := range es.aliases {
		if alias == indexOrAlias {
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
	es.aliases[index] = alias

	return &MockMethods{
		StatusCode: 200,
	}
}

func (es *InMemoryElasticsearch) DeleteAlias(index string, alias string) *MockMethods {
	if es.mock != nil {
		return es.mock
	}

	var responseStatusCode int
	if es.aliases[index] == alias {
		delete(es.aliases, index)
		responseStatusCode = 200
	} else {
		responseStatusCode = 404
	}

	return &MockMethods{
		StatusCode: responseStatusCode,
	}
}
