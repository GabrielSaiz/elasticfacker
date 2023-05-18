package elasticfacker

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
