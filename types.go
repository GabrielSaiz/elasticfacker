package elasticfacker

import "net/http"

type MockMethods struct {
	StatusCode   int
	Status       string
	BodyAsString string
}
type InMemoryElasticsearch struct {
	indicesAlias     map[string]map[string]interface{}
	indicesDocuments map[string][]Document
	aliases          map[string]interface{}
	mock             *MockMethods
	server           *http.Server
}

type IndexFake struct {
	Health string `json:"health"`
	Status string `json:"status"`
	Index  string `json:"index"`
}

type ProductIndexFake struct {
	Aliases map[string]interface{} `json:"aliases"`
}

type IndexMapFake map[string]ProductIndexFake

type ElasticSearchRequest struct {
	Id     string                     `json:"id"`
	Params ElasticSearchRequestParams `json:"params"`
}

type ElasticSearchRequestScriptQuery struct {
	Size   string      `json:"size"`
	Source []string    `json:"_source,omitempty"`
	Query  interface{} `json:"query"`
}

type ElasticSearchRequestParams struct {
	SearchTerm string `json:"search_term"`
	Size       string `json:"size"`
}

type Document struct {
	Index  string                 `json:"_index"`
	Id     string                 `json:"_id"`
	Score  string                 `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type ElasticSearchResponseFake struct {
	Took   int                             `json:"took"`
	Shards ElasticSearchResponseFakeShards `json:"_shards"`
	Hits   ElasticSearchResponseFakeHits   `json:"hits"`
}

type ElasticSearchCountResponseFake struct {
	Count  int                             `json:"count"`
	Shards ElasticSearchResponseFakeShards `json:"_shards"`
}

type ElasticSearchResponseFakeShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type ElasticSearchResponseFakeHits struct {
	Total ElasticSearchResponseFakeHitsTotal `json:"total"`
	Hits  []Document                         `json:"hits"`
}

type ElasticSearchResponseFakeHitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
