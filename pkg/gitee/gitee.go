package gitee

// A package interacting with Gitee Git hosting service.

import (
	"net/http"
	"net/url"

	"github.com/google/go-github/v45/github"
)

const DefaultBaseURL = "https://gitee.com/api/v5/"

type RequestDataGenerator interface {
	Validate() error
	ConstructEndpointUrl() string
	MakeRequestData() map[string]string
	RequestMethod() string
}

type CustomClient struct {
	*github.Client // Drop this one day

	accessToken string
	client      *http.Client

	BranchesService        *BranchesService
	RepositoriesService    *MyRepositoriesService
	RepositoryFilesService *RepositoryFilesService
	PullRequestsService    *PullRequestsService
	WebhooksService        *WebhooksService
}

type Opt struct {
	HttpClient  *http.Client
	AccessToken string
}

func NewClient(opt Opt) *CustomClient {
	httpClient := opt.HttpClient
	if opt.HttpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &CustomClient{
		Client:      github.NewClient(httpClient),
		client:      httpClient,
		accessToken: opt.AccessToken,
	}
	c.BaseURL = baseURL
	c.UploadURL = baseURL
	c.RepositoriesService = &MyRepositoriesService{client: c}
	c.BranchesService = &BranchesService{client: c}
	c.RepositoryFilesService = &RepositoryFilesService{client: c}
	c.PullRequestsService = &PullRequestsService{client: c}
	c.WebhooksService = &WebhooksService{client: c}
	return c
}
