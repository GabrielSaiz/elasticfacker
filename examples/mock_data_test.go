package examples

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/assert"

	"context"
	"github.com/gabrielsaiz/elasticfacker"
	"testing"
	"time"
)

func TestInfoRequest(t *testing.T) {
	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	req := esapi.InfoRequest{}

	res, err := req.Do(context.Background(), esClient)
	assert.Nil(t, err)
	defer res.Body.Close()

	var info map[string]string
	err = json.NewDecoder(res.Body).Decode(&info)

	assert.Nil(t, err)
	assert.NotNil(t, info)
}

func TestIndicesExistsRequest(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "Index not found",
			indexName: "products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "Index is a teapot",
			indexName: "products-test-teapot",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 418,
				Status:     "I'm a teapot",
			},
			expected: false,
		},
		{
			name:      "Index found",
			indexName: "products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesExistsRequest{
				Index: []string{subtest.indexName},
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)
		})
	}
}

func TestIndicesCreateRequest(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "IndexAlreadyExists",
			indexName: "products-test-already-exists",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 409,
				Status:     "Already Exists",
			},
			expected: false,
		},
		{
			name:      "CreatedIndex",
			indexName: "products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesCreateRequest{
				Index: subtest.indexName,
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)
		})
	}
}

func TestCatIndicesRequest(t *testing.T) {
	subtests := []struct {
		name             string
		indexNamePattern string
		expected         bool
		mockMethod       elasticfacker.MockMethods
	}{
		{
			name:             "IndexNotFound",
			indexNamePattern: "products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:             "IndexIsATeapot",
			indexNamePattern: "products-test-teapot",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 418,
				Status:     "I'm a teapot",
			},
			expected: false,
		},
		{
			name:             "IndexFound",
			indexNamePattern: "products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
				BodyAsString: `[
					{
						"health": "yellow",
						"status": "open",
						"index": "products-test-000001",
						"uuid": "u8FNjxh8Rfy_awN11oDKYQ",
						"pri": "1",
						"rep": "1",
						"docs.count": "1200",
						"docs.deleted": "0",
						"store.size": "88.1kb",
						"pri.store.size": "88.1kb"
					},
					{
						"health": "green",
						"status": "open",
						"index": "products-test-000002",
						"uuid": "nYFWZEO7TUiOjLQXBaYJpA",
						"pri": "1",
						"rep": "0",
						"docs.count": "0",
						"docs.deleted": "0",
						"store.size": "260b",
						"pri.store.size": "260b"
					}
				]`,
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.CatIndicesRequest{
				Index:  []string{subtest.indexNamePattern},
				Format: "json",
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)

			switch subtest.name {
			case "IndexFound":
				var indices []map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&indices)

				assert.Nil(t, err)
				assert.NotNil(t, indices)
				indexNames := make([]string, 0)
				for _, indexInfo := range indices {
					if indexName, ok := indexInfo["index"].(string); ok {
						indexNames = append(indexNames, indexName)
					}
				}
				assert.Len(t, indexNames, 2)
			}
		})
	}
}

func TestIndicesGetAliasRequestByIndex(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "IndexFound",
			indexName: "products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
				BodyAsString: `{
  					"products-test": {
    					"aliases": {
      						"products-test-alias1": {},
							"products-test-alias2": {}
    					}
  					}
				}`,
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesGetAliasRequest{
				Index: []string{subtest.indexName},
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)

			switch subtest.name {
			case "IndexFound":
				var indices map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&indices)

				assert.Nil(t, err)
				assert.NotNil(t, indices)
				aliasNames := make([]string, 0)
				for _, indexInfo := range indices {
					indexInfo := indexInfo.(map[string]interface{})
					aliases, ok := indexInfo["aliases"].(map[string]interface{})
					assert.True(t, ok)

					for aliasName, _ := range aliases {
						aliasNames = append(aliasNames, aliasName)
					}
				}
				assert.Len(t, aliasNames, 2)
			}
		})
	}
}

func TestIndicesDeleteRequest(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "DeletedIndex",
			indexName: "products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesDeleteRequest{
				Index: []string{subtest.indexName},
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)
		})
	}
}

func TestIndicesGetAliasRequestByAlias(t *testing.T) {
	subtests := []struct {
		name       string
		aliasName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "AliasNotFound",
			aliasName: "alias-products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "AliasFound",
			aliasName: "alias-products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
				BodyAsString: `{
  					"products-test": {
    					"aliases": {
      						"alias-products-test": {},
    					}
  					}
				}`,
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesGetAliasRequest{
				Name: []string{subtest.aliasName},
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)

			switch subtest.name {
			case "IndexFound":
				var indices map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&indices)

				assert.Nil(t, err)
				assert.NotNil(t, indices)
				aliasNames := make([]string, 0)
				for _, indexInfo := range indices {
					indexInfo := indexInfo.(map[string]interface{})
					aliases, ok := indexInfo["aliases"].(map[string]interface{})
					assert.True(t, ok)

					for aliasName, _ := range aliases {
						aliasNames = append(aliasNames, aliasName)
					}
				}
				assert.Len(t, aliasNames, 1)
			}
		})
	}
}

func TestIndicesDeleteAliasRequest(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		aliasName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			aliasName: "alias-products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "AliasNotFound",
			indexName: "products-test",
			aliasName: "alias-products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "DeletedAlias",
			indexName: "products-test",
			aliasName: "alias-products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesDeleteAliasRequest{
				Index: []string{subtest.indexName},
				Name:  []string{subtest.aliasName},
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)
		})
	}
}

func TestIndicesPutAliasRequest(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		aliasName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			aliasName: "alias-products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name:      "AliasAlreadyExists",
			indexName: "products-test",
			aliasName: "alias-products-test-already-exists",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 409,
				Status:     "Already Exists",
			},
			expected: false,
		},
		{
			name:      "CreateAlias",
			indexName: "products-test",
			aliasName: "alias-products-test",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 200,
				Status:     "OK",
			},
			expected: true,
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {

			esFacker.SetMockMethods(&subtest.mockMethod)

			req := esapi.IndicesPutAliasRequest{
				Index: []string{subtest.indexName},
				Name:  subtest.aliasName,
			}

			res, err := req.Do(context.Background(), esClient)
			assert.Nil(t, err)
			defer res.Body.Close()

			assert.Equal(t, subtest.expected, res.StatusCode == 200)
		})
	}
}
