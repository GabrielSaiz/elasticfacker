# elasticfacker

Elasticsearch simulator for API endpoint with the proposal of make the ES Golang Client testable

```
Disclaimer: This is a work in progress...

Please note that this is a very simplified approach to mocking and may not be suitable for all use cases. 
In particular, this approach only allows you to define one mock response for each method, 
regardless of the arguments passed to it. If you need more complex behaviour, 
such as returning different responses depending on the arguments or changing responses over time, 
you would need a more sophisticated solution.


```

## Operations

Implemented operations:


- HEAD /{indexName} -> esapi.IndicesExistsRequest
- GET /_cat/indices/{indexNamePattern} -> esapi.CatIndicesRequest
- GET /{indexName}/_alias -> esapi.IndicesGetAliasRequest
- GET /_alias/{aliasName} -> esapi.IndicesGetAliasRequest
- PUT /{indexName} -> esapi.IndicesCreateRequest
- PUT /{indexName}/_aliases/{aliasName} -> esapi.IndicesPutAliasRequest
- DELETE /{indexName} -> esapi.IndicesDeleteRequest
- DELETE /{indexName}/_aliases/{aliasName} -> esapi.IndicesDeleteAliasRequest

- 

## How to use

There are 2 options to use this library:
- Using the Elasticfacker as a in memory server
- Using the Elasticfaceker with mock data

### Starting the server

Starting the server is very simple:

- Add this import to your code:

```go
import "github.com/elasticfacker/elasticfacker"
```

- Then you can start the server with the following code into your test:

```go
esFacker := elasticfacker.NewInMemoryElasticsearch()
esFacker.Start("localhost:9200")
defer esFacker.Stop()
```

- Then you can use the Elasticsearch client as you would normally do:

```go
    esClient, error := elasticsearch.NewDefaultClient()

    req := esapi.IndicesExistsRequest{
        Index: []string{indexName},
    }

    res, err := req.Do(context.Background(), esClient)
    if err != nil {
        log.Get().Infof("Error when checking if index exists '%s': %s", indexName, err)
        return 
    }
    defer res.Body.Close()

    if res.StatusCode == 200 {
        log.Get().Infof("Index '%s' already exists", indexName)
        return
    } else {
		log.Get().Infof("Index '%s' does not exists", indexName)
		return
    }
}
```

### Using as a memory server

Now you can use the server as a memory server:

After starting the elasticfaker server, you can use it as an Elasticsearch index, a facade, currently it store the name of the index, the alias and the relation between index and aliases... 

So it's just an index to add or to get this info... for a better answer, you can use the mock data.   

#### Memory server example

```go
func TestElasticApiClient(t *testing.T) {
    subtests := []struct {
        name       string
		indexName  string
        expected   int
    }{
        {
            name: "Index not found",
			indexName: "products-test-not-found",
            expected: 404,
        },
        {
            name: "Index found",
			indexName: "products-test",
            expected: 200,
        },
	}

    esClient, error := elasticsearch.NewDefaultClient()
	if error != nil {
        t.Errorf("Error when creating the Elasticsearch client: %s", error)
    }

    esFacker := elasticfacker.NewInMemoryElasticsearch()
    esFacker.Start("localhost:9200")
    defer esFacker.Stop()
    time.Sleep(2 * time.Second)

    for _, subtest := range subtests {
    
        t.Run(subtest.name, func(t *testing.T) {
            req1 := esapi.IndicesCreateRequest{
                Index: subtest.indexName,
            }

            res1, err1 := req1.Do(context.Background(), esClient)
            assert.Nil(t, err1)
            defer res1.Body.Close()
            

            if res1.IsError() || res1.StatusCode != 200 {
                t.Errorf("Error when creating the index '%s': %s", subtest.indexName, res1.String())
            }
			
            req2 := esapi.IndicesExistsRequest{
                Index: []string{subtest.indexName},
            }
            
            res2, err2 := req2.Do(context.Background(), esClient)
            assert.Nil(t, err2)
			defer res2.Body.Close()
            
            assert.Equal(t, subtest.expected, res2.StatusCode)
        })
    }
}

```


### Using the mock data

Elasticseach API uses an HTTP Rest API, so you can use the mock data to simulate the responses of the API.

The MockMethod is an struct that contains the following fields:

```go
type MockMethods struct {
    StatusCode   int
    Status       string
    BodyAsString string
}
```

#### Mock data example

```go

func TestElasticApiClient(t *testing.T) {
	subtests := []struct {
		name       string
		indexName  string
		expected   bool
		mockMethod elasticfacker.MockMethods
	}{
		{
			name: "Index not found",
			indexName: "products-test-not-found",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 404,
				Status:     "Not Found",
			},
			expected: false,
		},
		{
			name: "Index is a teapot",
			indexName: "products-test-teapot",
			mockMethod: elasticfacker.MockMethods{
				StatusCode: 418,
				Status:     "I'm a teapot",
			},
			expected: false,
		},
		{
			name: "Index found",
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

```