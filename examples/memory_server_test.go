package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gabrielsaiz/elasticfacker"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestIndicesRequest(t *testing.T) {
	time.Sleep(1 * time.Second)
	subtests := []struct {
		name      string
		indexName string
	}{
		{
			name:      "CreatedIndex",
			indexName: "products-test",
		},
		{
			name:      "IndexAlreadyExists",
			indexName: "products-test",
		},
		{
			name:      "IndexDoesNotExist",
			indexName: "products-test-2",
		},
		{
			name:      "IndexExists",
			indexName: "products-test",
		},
		{
			name:      "IndexNotFound",
			indexName: "products-test-2*",
		},
		{
			name:      "IndexFound",
			indexName: "products-test",
		},
		{
			name:      "IndexDeletedNotFound",
			indexName: "products-test-2",
		},
		{
			name:      "IndexDeleted",
			indexName: "products-test",
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

			switch subtest.name {
			case "CreatedIndex":
				req := esapi.IndicesCreateRequest{
					Index: subtest.indexName,
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)
			case "IndexAlreadyExists":
				req := esapi.IndicesCreateRequest{
					Index: subtest.indexName,
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 409)
			case "IndexDoesNotExist":
				req := esapi.IndicesExistsRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "IndexExists":
				req := esapi.IndicesExistsRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)
			case "IndexNotFound":
				req := esapi.CatIndicesRequest{
					Index:  []string{subtest.indexName},
					Format: "json",
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "IndexFound":
				req := esapi.CatIndicesRequest{
					Index:  []string{subtest.indexName},
					Format: "json",
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)

				var indices []map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&indices)

				assert.Nil(t, err)
				assert.NotNil(t, indices)
				for _, indexInfo := range indices {
					if indexName, ok := indexInfo["index"].(string); ok {
						assert.True(t, ok)
						assert.NotNil(t, indexName)
					}
				}
			case "IndexDeletedNotFound":
				req := esapi.IndicesDeleteRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "IndexDeleted":
				req := esapi.IndicesDeleteRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)
			}

		})
	}
}

func TestIndicesAliasRequest(t *testing.T) {
	time.Sleep(1 * time.Second)
	subtests := []struct {
		name      string
		indexName string
		aliasName string
	}{
		{
			name:      "CreatedAliasIndexNotFound",
			indexName: "products-test-2",
			aliasName: "products-test-alias",
		},
		{
			name:      "CreatedAlias",
			indexName: "products-test",
			aliasName: "products-test-alias",
		},
		{
			name:      "CreateAliasAlreadyExists",
			indexName: "products-test",
			aliasName: "products-test-alias",
		},
		{
			name:      "AliasNotFound",
			indexName: "products-test",
			aliasName: "products-test-alias-2",
		},
		{
			name:      "AliasFound",
			indexName: "products-test",
			aliasName: "products-test-alias",
		},
		{
			name:      "GetAliasByIndexNotFound",
			indexName: "products-test-2",
		},
		{
			name:      "GetAliasByIndexFound",
			indexName: "products-test",
		},
		{
			name:      "DeleteAliasIndexNotFound",
			indexName: "products-test-2",
			aliasName: "products-test-alias",
		},
		{
			name:      "DeleteAliasNotFound",
			indexName: "products-test",
			aliasName: "products-test-alias-2",
		},
		{
			name:      "DeleteAlias",
			indexName: "products-test",
			aliasName: "products-test-alias",
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	req := esapi.IndicesCreateRequest{
		Index: "products-test",
	}

	res, err := req.Do(context.Background(), esClient)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 200)

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "CreatedAliasIndexNotFound":
				req := esapi.IndicesPutAliasRequest{
					Index: []string{subtest.indexName},
					Name:  subtest.aliasName,
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 500)
			case "CreatedAlias":
				req := esapi.IndicesPutAliasRequest{
					Index: []string{subtest.indexName},
					Name:  subtest.aliasName,
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)

			case "CreateAliasAlreadyExists":
				req := esapi.IndicesPutAliasRequest{
					Index: []string{subtest.indexName},
					Name:  subtest.aliasName,
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 409)

			case "AliasNotFound":
				req := esapi.IndicesGetAliasRequest{
					Name: []string{subtest.aliasName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "AliasFound":
				req := esapi.IndicesGetAliasRequest{
					Name: []string{subtest.aliasName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)

			case "GetAliasByIndexNotFound":
				req := esapi.IndicesGetAliasRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "GetAliasByIndexFound":
				req := esapi.IndicesGetAliasRequest{
					Index: []string{subtest.indexName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)
			case "DeleteAliasIndexNotFound":
				req := esapi.IndicesDeleteAliasRequest{
					Index: []string{subtest.indexName},
					Name:  []string{subtest.aliasName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 500)
			case "DeleteAliasNotFound":
				req := esapi.IndicesDeleteAliasRequest{
					Index: []string{subtest.indexName},
					Name:  []string{subtest.aliasName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 404)
			case "DeleteAlias":
				req := esapi.IndicesDeleteAliasRequest{
					Index: []string{subtest.indexName},
					Name:  []string{subtest.aliasName},
				}

				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.True(t, res.StatusCode == 200)
			}

		})
	}
}

func TestSearchTemplateRequest(t *testing.T) {
	time.Sleep(1 * time.Second)
	subtests := []struct {
		name      string
		indexName string
		body      *strings.Reader
	}{
		{
			name:      "BadRequest",
			indexName: "products-test",
			body:      strings.NewReader("{badRequest}"),
		},
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			body:      buildBody("templateId", "test", 10),
		},
		{
			name:      "SearchTemplateEmpty",
			indexName: "products-test",
			body:      buildBody("templateId", "test", 10),
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	req := esapi.IndicesCreateRequest{
		Index: "products-test",
	}

	res, err := req.Do(context.Background(), esClient)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 200)

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {
			req := esapi.SearchTemplateRequest{
				Index: []string{subtest.indexName},
				Body:  subtest.body,
			}

			switch subtest.name {
			case "BadRequest":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 400, res.StatusCode)
			case "IndexNotFound":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 404, res.StatusCode)
			case "SearchTemplateEmpty":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 200, res.StatusCode)
			}

		})
	}
}

func TestSearchRequest(t *testing.T) {
	time.Sleep(1 * time.Second)
	subtests := []struct {
		name      string
		indexName string
		body      *strings.Reader
	}{
		{
			name:      "BadRequest",
			indexName: "products-test",
			body:      strings.NewReader("{badRequest}"),
		},
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			body:      buildScriptQueryBody("test", 10),
		},
		{
			name:      "SearchEmpty",
			indexName: "products-test",
			body:      buildScriptQueryBody("test", 10),
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	req := esapi.IndicesCreateRequest{
		Index: "products-test",
	}

	res, err := req.Do(context.Background(), esClient)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 200)

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {
			req := esapi.SearchRequest{
				Index: []string{subtest.indexName},
				Body:  subtest.body,
			}

			switch subtest.name {
			case "BadRequest":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 400, res.StatusCode)
			case "IndexNotFound":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 404, res.StatusCode)
			case "SearchEmpty":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 200, res.StatusCode)
			}

		})
	}
}

func TestCountRequest(t *testing.T) {
	time.Sleep(1 * time.Second)
	subtests := []struct {
		name      string
		indexName string
		body      *strings.Reader
	}{
		{
			name:      "BadRequest",
			indexName: "products-test",
			body:      strings.NewReader("{badRequest}"),
		},
		{
			name:      "IndexNotFound",
			indexName: "products-test-not-found",
			body:      buildScriptQueryBody("test", 10),
		},
		{
			name:      "CountZero",
			indexName: "products-test",
			body:      buildScriptQueryBody("test", 10),
		},
	}

	esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
		t.Errorf("Error when creating the Elasticsearch client: %s", error)
	}

	esFacker := elasticfacker.NewInMemoryElasticsearch()
	esFacker.Start("localhost:9200")
	defer esFacker.Stop()

	req := esapi.IndicesCreateRequest{
		Index: "products-test",
	}

	fmt.Printf("Req: %v \n", req)
	fmt.Printf("Context: %v \n", context.Background())
	fmt.Printf("esClient: %v \n", esClient)

	res, err := req.Do(context.Background(), esClient)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.True(t, res.StatusCode == 200)

	for _, subtest := range subtests {
		time.Sleep(1 * time.Second)

		t.Run(subtest.name, func(t *testing.T) {
			req := esapi.CountRequest{
				Index: []string{subtest.indexName},
				Body:  subtest.body,
			}

			switch subtest.name {
			case "BadRequest":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 400, res.StatusCode)
			case "IndexNotFound":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 404, res.StatusCode)
			case "CountZero":
				res, err := req.Do(context.Background(), esClient)
				assert.Nil(t, err)
				defer res.Body.Close()

				assert.Equal(t, 200, res.StatusCode)
			}

		})
	}
}
