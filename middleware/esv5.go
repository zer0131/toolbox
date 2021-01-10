package middleware

// esv5的封装并没有是否用ns，因为es每个节点内置sniff功能，用于探测整个集群node变化，
// 相当于sdk内置dns功能，只是不是domain name，而是ip:port。
// 如果我们通过ns接管这部分工作，缺点有两个：
// 1. es和ns服务之间node变化数据要打通，这块就有运维成本
// 2. sdk不关心sniff得到的节点，只通过ns sdk拿到的列表访问，不直观
// es和redis/mysql或者http不一样，后面这几种没有内置ns功能。

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v5"

	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
)

type ESV5 struct {
	*elastic.Client
}

type DpESV5 interface {
	// String returns a string representation of the client status.
	String() string
	// IsRunning returns true if the background processes of the client are
	// running, false otherwise.
	IsRunning() bool
	// Start starts the background processes like sniffing the cluster and
	// periodic health checks. You don't need to run Start when creating a
	// client with NewClient; the background processes are run by default.
	//
	// If the background processes are already running, this is a no-op.
	Start()
	// Stop stops the background processes that the client is running,
	// i.e. sniffing the cluster periodically and running health checks
	// on the nodes.
	//
	// If the background processes are not running, this is a no-op.
	Stop()
	// PerformRequest does a HTTP request to Elasticsearch.
	// See PerformRequestWithContentType for details.
	PerformRequest(ctx context.Context, method, path string, params url.Values, body interface{}, ignoreErrors ...int) (*elastic.Response, error)
	// PerformRequestWithContentType executes a HTTP request with a specific content type.
	// It returns a response (which might be nil) and an error on failure.
	//
	// Optionally, a list of HTTP error codes to ignore can be passed.
	// This is necessary for services that expect e.g. HTTP status 404 as a
	// valid outcome (Exists, IndicesExists, IndicesTypeExists).
	PerformRequestWithContentType(ctx context.Context, method, path string, params url.Values, body interface{}, contentType string, ignoreErrors ...int) (*elastic.Response, error)
	// PerformRequestWithOptions executes a HTTP request with the specified options.
	// It returns a response (which might be nil) and an error on failure.
	PerformRequestWithOptions(ctx context.Context, opt elastic.PerformRequestOptions) (*elastic.Response, error)
	// Index a document.
	Index() *elastic.IndexService
	// Get a document.
	Get() *elastic.GetService
	// MultiGet retrieves multiple documents in one roundtrip.
	MultiGet() *elastic.MgetService
	// Mget retrieves multiple documents in one roundtrip.
	Mget() *elastic.MgetService
	// Delete a document.
	Delete() *elastic.DeleteService
	// DeleteByQuery deletes documents as found by a query.
	DeleteByQuery(indices ...string) *elastic.DeleteByQueryService
	// Update a document.
	Update() *elastic.UpdateService
	// UpdateByQuery performs an update on a set of documents.
	UpdateByQuery(indices ...string) *elastic.UpdateByQueryService
	// Bulk is the entry point to mass insert/update/delete documents.
	Bulk() *elastic.BulkService
	// BulkProcessor allows setting up a concurrent processor of bulk requests.
	BulkProcessor() *elastic.BulkProcessorService
	// Reindex copies data from a source index into a destination index.
	//
	// See https://www.elastic.co/guide/en/elasticsearch/reference/5.2/docs-reindex.html
	// for details on the Reindex API.
	Reindex() *elastic.ReindexService
	// TermVectors returns information and statistics on terms in the fields
	// of a particular document.
	TermVectors(index, typ string) *elastic.TermvectorsService
	// MultiTermVectors returns information and statistics on terms in the fields
	// of multiple documents.
	MultiTermVectors() *elastic.MultiTermvectorService
	// Search is the entry point for searches.
	Search(indices ...string) *elastic.SearchService
	// Suggest returns a service to return suggestions.
	Suggest(indices ...string) *elastic.SuggestService
	// MultiSearch is the entry point for multi searches.
	MultiSearch() *elastic.MultiSearchService
	// Count documents.
	Count(indices ...string) *elastic.CountService
	// Explain computes a score explanation for a query and a specific document.
	Explain(index, typ, id string) *elastic.ExplainService
	// Validate allows a user to validate a potentially expensive query without executing it.
	Validate(indices ...string) *elastic.ValidateService
	// SearchShards returns statistical information about nodes and shards.
	SearchShards(indices ...string) *elastic.SearchShardsService
	// FieldCaps returns statistical information about fields in indices.
	FieldCaps(indices ...string) *elastic.FieldCapsService
	// FieldStats returns statistical information about fields in indices.
	FieldStats(indices ...string) *elastic.FieldStatsService
	// Exists checks if a document exists.
	Exists() *elastic.ExistsService
	// Scroll through documents. Use this to efficiently scroll through results
	// while returning the results to a client.
	Scroll(indices ...string) *elastic.ScrollService
	// ClearScroll can be used to clear search contexts manually.
	ClearScroll(scrollIds ...string) *elastic.ClearScrollService
	// CreateIndex returns a service to create a new index.
	CreateIndex(name string) *elastic.IndicesCreateService
	// DeleteIndex returns a service to delete an index.
	DeleteIndex(indices ...string) *elastic.IndicesDeleteService
	// IndexExists allows to check if an index exists.
	IndexExists(indices ...string) *elastic.IndicesExistsService
	// ShrinkIndex returns a service to shrink one index into another.
	ShrinkIndex(source, target string) *elastic.IndicesShrinkService
	// RolloverIndex rolls an alias over to a new index when the existing index
	// is considered to be too large or too old.
	RolloverIndex(alias string) *elastic.IndicesRolloverService
	// TypeExists allows to check if one or more types exist in one or more indices.
	TypeExists() *elastic.IndicesExistsTypeService
	// IndexStats provides statistics on different operations happining
	// in one or more indices.
	IndexStats(indices ...string) *elastic.IndicesStatsService
	// OpenIndex opens an index.
	OpenIndex(name string) *elastic.IndicesOpenService
	// CloseIndex closes an index.
	CloseIndex(name string) *elastic.IndicesCloseService
	// IndexGet retrieves information about one or more indices.
	// IndexGet is only available for Elasticsearch 1.4 or later.
	IndexGet(indices ...string) *elastic.IndicesGetService
	// IndexGetSettings retrieves settings of all, one or more indices.
	IndexGetSettings(indices ...string) *elastic.IndicesGetSettingsService
	// IndexPutSettings sets settings for all, one or more indices.
	IndexPutSettings(indices ...string) *elastic.IndicesPutSettingsService
	// IndexSegments retrieves low level segment information for all, one or more indices.
	IndexSegments(indices ...string) *elastic.IndicesSegmentsService
	// IndexAnalyze performs the analysis process on a text and returns the
	// token breakdown of the text.
	IndexAnalyze() *elastic.IndicesAnalyzeService
	// Forcemerge optimizes one or more indices.
	// It replaces the deprecated Optimize API.
	Forcemerge(indices ...string) *elastic.IndicesForcemergeService
	// Refresh asks Elasticsearch to refresh one or more indices.
	Refresh(indices ...string) *elastic.RefreshService
	// Flush asks Elasticsearch to free memory from the index and
	// flush data to disk.
	Flush(indices ...string) *elastic.IndicesFlushService
	// Alias enables the caller to add and/or remove aliases.
	Alias() *elastic.AliasService
	// Aliases returns aliases by index name(s).
	Aliases() *elastic.AliasesService
	// GetTemplate gets a search template.
	// Use IndexXXXTemplate funcs to manage index templates.
	GetTemplate() *elastic.GetTemplateService
	// PutTemplate creates or updates a search template.
	// Use IndexXXXTemplate funcs to manage index templates.
	PutTemplate() *elastic.PutTemplateService
	// DeleteTemplate deletes a search template.
	// Use IndexXXXTemplate funcs to manage index templates.
	DeleteTemplate() *elastic.DeleteTemplateService
	// IndexGetTemplate gets an index template.
	// Use XXXTemplate funcs to manage search templates.
	IndexGetTemplate(names ...string) *elastic.IndicesGetTemplateService
	// IndexTemplateExists gets check if an index template exists.
	// Use XXXTemplate funcs to manage search templates.
	IndexTemplateExists(name string) *elastic.IndicesExistsTemplateService
	// IndexPutTemplate creates or updates an index template.
	// Use XXXTemplate funcs to manage search templates.
	IndexPutTemplate(name string) *elastic.IndicesPutTemplateService
	// IndexDeleteTemplate deletes an index template.
	// Use XXXTemplate funcs to manage search templates.
	IndexDeleteTemplate(name string) *elastic.IndicesDeleteTemplateService
	// GetMapping gets a mapping.
	GetMapping() *elastic.IndicesGetMappingService
	// PutMapping registers a mapping.
	PutMapping() *elastic.IndicesPutMappingService
	// GetFieldMapping gets mapping for fields.
	GetFieldMapping() *elastic.IndicesGetFieldMappingService
	// IngestPutPipeline adds pipelines and updates existing pipelines in
	// the cluster.
	IngestPutPipeline(id string) *elastic.IngestPutPipelineService
	// IngestGetPipeline returns pipelines based on ID.
	IngestGetPipeline(ids ...string) *elastic.IngestGetPipelineService
	// IngestDeletePipeline deletes a pipeline by ID.
	IngestDeletePipeline(id string) *elastic.IngestDeletePipelineService
	// IngestSimulatePipeline executes a specific pipeline against the set of
	// documents provided in the body of the request.
	IngestSimulatePipeline() *elastic.IngestSimulatePipelineService
	// ClusterHealth retrieves the health of the cluster.
	ClusterHealth() *elastic.ClusterHealthService
	// ClusterState retrieves the state of the cluster.
	ClusterState() *elastic.ClusterStateService
	// ClusterStats retrieves cluster statistics.
	ClusterStats() *elastic.ClusterStatsService
	// NodesInfo retrieves one or more or all of the cluster nodes information.
	NodesInfo() *elastic.NodesInfoService
	// NodesStats retrieves one or more or all of the cluster nodes statistics.
	NodesStats() *elastic.NodesStatsService
	// TasksCancel cancels tasks running on the specified nodes.
	TasksCancel() *elastic.TasksCancelService
	// TasksList retrieves the list of tasks running on the specified nodes.
	TasksList() *elastic.TasksListService
	// TasksGetTask retrieves a task running on the cluster.
	TasksGetTask() *elastic.TasksGetTaskService
	// SnapshotCreate creates a snapshot.
	SnapshotCreate(repository string, snapshot string) *elastic.SnapshotCreateService
	// SnapshotCreateRepository creates or updates a snapshot repository.
	SnapshotCreateRepository(repository string) *elastic.SnapshotCreateRepositoryService
	// SnapshotDeleteRepository deletes a snapshot repository.
	SnapshotDeleteRepository(repositories ...string) *elastic.SnapshotDeleteRepositoryService
	// SnapshotGetRepository gets a snapshot repository.
	SnapshotGetRepository(repositories ...string) *elastic.SnapshotGetRepositoryService
	// SnapshotVerifyRepository verifies a snapshot repository.
	SnapshotVerifyRepository(repository string) *elastic.SnapshotVerifyRepositoryService
	// ElasticsearchVersion returns the version number of Elasticsearch
	// running on the given URL.
	ElasticsearchVersion(url string) (string, error)
	// IndexNames returns the names of all indices in the cluster.
	IndexNames() ([]string, error)
	// Ping checks if a given node in a cluster exists and (optionally)
	// returns some basic information about the Elasticsearch server,
	// e.g. the Elasticsearch version number.
	//
	// Notice that you need to specify a URL here explicitly.
	Ping(url string) *elastic.PingService
	// WaitForStatus waits for the cluster to have the given status.
	// This is a shortcut method for the ClusterHealth service.
	//
	// WaitForStatus waits for the specified timeout, e.g. "10s".
	// If the cluster will have the given state within the timeout, nil is returned.
	// If the request timed out, ErrTimeout is returned.
	WaitForStatus(status string, timeout string) error
	// WaitForGreenStatus waits for the cluster to have the "green" status.
	// See WaitForStatus for more details.
	WaitForGreenStatus(timeout string) error
	// WaitForYellowStatus waits for the cluster to have the "yellow" status.
	// See WaitForStatus for more details.
	WaitForYellowStatus(timeout string) error
}

func (esv5 *ESV5) String() string {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.String()
}

func (esv5 *ESV5) IsRunning() bool {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IsRunning()
}

func (esv5 *ESV5) Start() {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	esv5.Client.Start()
}

func (esv5 *ESV5) Stop() {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	esv5.Client.Stop()
}

func (esv5 *ESV5) PerformRequest(ctx context.Context, method string, path string, params url.Values, body interface{}, ignoreErrors ...int) (*elastic.Response, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.PerformRequest(ctx, method, path, params, body, ignoreErrors...)
}

func (esv5 *ESV5) PerformRequestWithContentType(ctx context.Context, method string, path string, params url.Values, body interface{}, contentType string, ignoreErrors ...int) (*elastic.Response, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.PerformRequestWithContentType(ctx, method, path, params, body, contentType, ignoreErrors...)
}

func (esv5 *ESV5) PerformRequestWithOptions(ctx context.Context, opt elastic.PerformRequestOptions) (*elastic.Response, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.PerformRequestWithOptions(ctx, opt)
}

func (esv5 *ESV5) Index() *elastic.IndexService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Index()
}

func (esv5 *ESV5) Get() *elastic.GetService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Get()
}

func (esv5 *ESV5) MultiGet() *elastic.MgetService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.MultiGet()
}

func (esv5 *ESV5) Mget() *elastic.MgetService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Mget()
}

func (esv5 *ESV5) Delete() *elastic.DeleteService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Delete()
}

func (esv5 *ESV5) DeleteByQuery(indices ...string) *elastic.DeleteByQueryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.DeleteByQuery(indices...)
}

func (esv5 *ESV5) Update() *elastic.UpdateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Update()
}

func (esv5 *ESV5) UpdateByQuery(indices ...string) *elastic.UpdateByQueryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.UpdateByQuery(indices...)
}

func (esv5 *ESV5) Bulk() *elastic.BulkService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Bulk()
}

func (esv5 *ESV5) BulkProcessor() *elastic.BulkProcessorService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.BulkProcessor()
}

func (esv5 *ESV5) Reindex() *elastic.ReindexService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Reindex()
}

func (esv5 *ESV5) TermVectors(index string, typ string) *elastic.TermvectorsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.TermVectors(index, typ)
}

func (esv5 *ESV5) MultiTermVectors() *elastic.MultiTermvectorService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.MultiTermVectors()
}

func (esv5 *ESV5) Search(indices ...string) *elastic.SearchService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Search(indices...)
}

func (esv5 *ESV5) Suggest(indices ...string) *elastic.SuggestService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Suggest(indices...)
}

func (esv5 *ESV5) MultiSearch() *elastic.MultiSearchService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.MultiSearch()
}

func (esv5 *ESV5) Count(indices ...string) *elastic.CountService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Count(indices...)
}

func (esv5 *ESV5) Explain(index string, typ string, id string) *elastic.ExplainService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Explain(index, typ, id)
}

func (esv5 *ESV5) Validate(indices ...string) *elastic.ValidateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Validate(indices...)
}

func (esv5 *ESV5) SearchShards(indices ...string) *elastic.SearchShardsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SearchShards(indices...)
}

func (esv5 *ESV5) FieldCaps(indices ...string) *elastic.FieldCapsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.FieldCaps(indices...)
}

func (esv5 *ESV5) FieldStats(indices ...string) *elastic.FieldStatsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.FieldStats(indices...)
}

func (esv5 *ESV5) Exists() *elastic.ExistsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Exists()
}

func (esv5 *ESV5) Scroll(indices ...string) *elastic.ScrollService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Scroll(indices...)
}

func (esv5 *ESV5) ClearScroll(scrollIds ...string) *elastic.ClearScrollService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ClearScroll(scrollIds...)
}

func (esv5 *ESV5) CreateIndex(name string) *elastic.IndicesCreateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.CreateIndex(name)
}

func (esv5 *ESV5) DeleteIndex(indices ...string) *elastic.IndicesDeleteService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.DeleteIndex(indices...)
}

func (esv5 *ESV5) IndexExists(indices ...string) *elastic.IndicesExistsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexExists(indices...)
}

func (esv5 *ESV5) ShrinkIndex(source string, target string) *elastic.IndicesShrinkService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ShrinkIndex(source, target)
}

func (esv5 *ESV5) RolloverIndex(alias string) *elastic.IndicesRolloverService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.RolloverIndex(alias)
}

func (esv5 *ESV5) TypeExists() *elastic.IndicesExistsTypeService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.TypeExists()
}

func (esv5 *ESV5) IndexStats(indices ...string) *elastic.IndicesStatsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexStats(indices...)
}

func (esv5 *ESV5) OpenIndex(name string) *elastic.IndicesOpenService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.OpenIndex(name)
}

func (esv5 *ESV5) CloseIndex(name string) *elastic.IndicesCloseService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.CloseIndex(name)
}

func (esv5 *ESV5) IndexGet(indices ...string) *elastic.IndicesGetService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexGet(indices...)
}

func (esv5 *ESV5) IndexGetSettings(indices ...string) *elastic.IndicesGetSettingsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexGetSettings(indices...)
}

func (esv5 *ESV5) IndexPutSettings(indices ...string) *elastic.IndicesPutSettingsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexPutSettings(indices...)
}

func (esv5 *ESV5) IndexSegments(indices ...string) *elastic.IndicesSegmentsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexSegments(indices...)
}

func (esv5 *ESV5) IndexAnalyze() *elastic.IndicesAnalyzeService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexAnalyze()
}

func (esv5 *ESV5) Forcemerge(indices ...string) *elastic.IndicesForcemergeService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Forcemerge(indices...)
}

func (esv5 *ESV5) Refresh(indices ...string) *elastic.RefreshService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Refresh(indices...)
}

func (esv5 *ESV5) Flush(indices ...string) *elastic.IndicesFlushService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Flush(indices...)
}

func (esv5 *ESV5) Alias() *elastic.AliasService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Alias()
}

func (esv5 *ESV5) Aliases() *elastic.AliasesService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Aliases()
}

func (esv5 *ESV5) GetTemplate() *elastic.GetTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.GetTemplate()
}

func (esv5 *ESV5) PutTemplate() *elastic.PutTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.PutTemplate()
}

func (esv5 *ESV5) DeleteTemplate() *elastic.DeleteTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.DeleteTemplate()
}

func (esv5 *ESV5) IndexGetTemplate(names ...string) *elastic.IndicesGetTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexGetTemplate(names...)
}

func (esv5 *ESV5) IndexTemplateExists(name string) *elastic.IndicesExistsTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexTemplateExists(name)
}

func (esv5 *ESV5) IndexPutTemplate(name string) *elastic.IndicesPutTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexPutTemplate(name)
}

func (esv5 *ESV5) IndexDeleteTemplate(name string) *elastic.IndicesDeleteTemplateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexDeleteTemplate(name)
}

func (esv5 *ESV5) GetMapping() *elastic.IndicesGetMappingService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.GetMapping()
}

func (esv5 *ESV5) PutMapping() *elastic.IndicesPutMappingService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.PutMapping()
}

func (esv5 *ESV5) GetFieldMapping() *elastic.IndicesGetFieldMappingService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.GetFieldMapping()
}

func (esv5 *ESV5) IngestPutPipeline(id string) *elastic.IngestPutPipelineService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IngestPutPipeline(id)
}

func (esv5 *ESV5) IngestGetPipeline(ids ...string) *elastic.IngestGetPipelineService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IngestGetPipeline(ids...)
}

func (esv5 *ESV5) IngestDeletePipeline(id string) *elastic.IngestDeletePipelineService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IngestDeletePipeline(id)
}

func (esv5 *ESV5) IngestSimulatePipeline() *elastic.IngestSimulatePipelineService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IngestSimulatePipeline()
}

func (esv5 *ESV5) ClusterHealth() *elastic.ClusterHealthService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ClusterHealth()
}

func (esv5 *ESV5) ClusterState() *elastic.ClusterStateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ClusterState()
}

func (esv5 *ESV5) ClusterStats() *elastic.ClusterStatsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ClusterStats()
}

func (esv5 *ESV5) NodesInfo() *elastic.NodesInfoService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.NodesInfo()
}

func (esv5 *ESV5) NodesStats() *elastic.NodesStatsService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.NodesStats()
}

func (esv5 *ESV5) TasksCancel() *elastic.TasksCancelService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.TasksCancel()
}

func (esv5 *ESV5) TasksList() *elastic.TasksListService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.TasksList()
}

func (esv5 *ESV5) TasksGetTask() *elastic.TasksGetTaskService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.TasksGetTask()
}

func (esv5 *ESV5) SnapshotCreate(repository string, snapshot string) *elastic.SnapshotCreateService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SnapshotCreate(repository, snapshot)
}

func (esv5 *ESV5) SnapshotCreateRepository(repository string) *elastic.SnapshotCreateRepositoryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SnapshotCreateRepository(repository)
}

func (esv5 *ESV5) SnapshotDeleteRepository(repositories ...string) *elastic.SnapshotDeleteRepositoryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SnapshotDeleteRepository(repositories...)
}

func (esv5 *ESV5) SnapshotGetRepository(repositories ...string) *elastic.SnapshotGetRepositoryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SnapshotGetRepository(repositories...)
}

func (esv5 *ESV5) SnapshotVerifyRepository(repository string) *elastic.SnapshotVerifyRepositoryService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.SnapshotVerifyRepository(repository)
}

func (esv5 *ESV5) ElasticsearchVersion(url string) (string, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.ElasticsearchVersion(url)
}

func (esv5 *ESV5) IndexNames() ([]string, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.IndexNames()
}

func (esv5 *ESV5) Ping(url string) *elastic.PingService {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.Ping(url)
}

func (esv5 *ESV5) WaitForStatus(status string, timeout string) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.WaitForStatus(status, timeout)
}

func (esv5 *ESV5) WaitForGreenStatus(timeout string) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.WaitForGreenStatus(timeout)
}

func (esv5 *ESV5) WaitForYellowStatus(timeout string) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.ESV5), startTime)
	return esv5.Client.WaitForYellowStatus(timeout)
}

func InitESV5(opt ...Esv5OptionsFunc) (DpESV5, error) {
	opts := defaultESV5Options
	for _, o := range opt {
		o(&opts)
	}

	if !strings.HasPrefix(opts.addr, "http://") {
		return nil, errors.New("addr should start with http://")
	}

	// http这里一些参数暂时写死，保证长连接
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   opts.connTimeout,
				KeepAlive: opts.keepalive,
				DualStack: true}).DialContext,

			IdleConnTimeout:       opts.idleTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 0,

			MaxIdleConns:        opts.maxIdleConnCount,
			MaxIdleConnsPerHost: opts.maxIdleConnCount,
		},
		Timeout: opts.timeout,
	}

	client, err := elastic.NewClient(
		elastic.SetURL(opts.addr),
		elastic.SetHttpClient(httpClient),
		elastic.SetTraceLog(elastic.Logger(log.New(os.Stdout, "[es] ", log.Ldate|log.Ltime|log.Lshortfile))),
		elastic.SetInfoLog(elastic.Logger(log.New(os.Stdout, "[es] ", log.Ldate|log.Ltime|log.Lshortfile))),
		elastic.SetErrorLog(elastic.Logger(log.New(os.Stdout, "[es] ", log.Ldate|log.Ltime|log.Lshortfile))))

	if err != nil {
		return nil, err
	}
	return &ESV5{client}, nil
}
