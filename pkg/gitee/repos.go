package gitee

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-github/v45/github"
)

type Repository struct {
	Id        *int64  `json:"id,omitempty"`
	FullName  *string `json:"full_name,omitempty"`
	HumanName *string `json:"human_name,omitempty"`
	URL       *string `json:"url,omitempty"`
	Path      *string `json:"path,omitempty"`
	Name      *string `json:"name,omitempty"`

	Description         *string   `json:"description,omitempty"`
	Private             *bool     `json:"private,omitempty"`
	Public              *bool     `json:"public,omitempty"`
	Internal            *bool     `json:"internal,omitempty"`
	Fork                *bool     `json:"fork,omitempty"`
	HtmlURL             *string   `json:"html_url,omitempty"`
	SshURL              *string   `json:"ssh_url,omitempty"`
	ForksURL            *string   `json:"forks_url,omitempty"`
	KeysURL             *string   `json:"keys_url,omitempty"`
	CollaboratorsURL    *string   `json:"collaborators_url,omitempty"`
	HooksURL            *string   `json:"hooks_url,omitempty"`
	BranchesURL         *string   `json:"branches_url,omitempty"`
	TagsURL             *string   `json:"tags_url,omitempty"`
	BlobsURL            *string   `json:"blobs_url,omitempty"`
	StargazersURL       *string   `json:"stargazers_url,omitempty"`
	ContributorsURL     *string   `json:"contributors_url,omitempty"`
	CommitsURL          *string   `json:"commits_url,omitempty"`
	CommentsURL         *string   `json:"comments_url,omitempty"`
	IssueCommentURL     *string   `json:"issue_comment_url,omitempty"`
	IssuesURL           *string   `json:"issues_url,omitempty"`
	PullsURL            *string   `json:"pulls_url,omitempty"`
	MilestonesURL       *string   `json:"milestones_url,omitempty"`
	NotificationsURL    *string   `json:"notifications_url,omitempty"`
	LabelsURL           *string   `json:"labels_url,omitempty"`
	ReleasesURL         *string   `json:"releases_url,omitempty"`
	Recommend           *bool     `json:"recommend,omitempty"`
	GVP                 *bool     `json:"gvp,omitempty"`
	Homepage            *string   `json:"homepage,omitempty"`
	Language            *string   `json:"language,omitempty"`
	ForksCount          *int      `json:"forks_count,omitempty"`
	StargazersCount     *int      `json:"wstargazers_count,omitempty"`
	WatchersCount       *int      `json:"watchers_count,omitempty"`
	DefaultBranch       *string   `json:"default_branch,omitempty"`
	OpenIssuesCount     *int      `json:"open_issues_count,omitempty"`
	HasIssues           *bool     `json:"has_issues,omitempty"`
	HasWiki             *bool     `json:"has_wiki,omitempty"`
	IssueComment        *bool     `json:"issue_comment,omitempty"`
	CanComment          *bool     `json:"can_comment,omitempty"`
	PullRequestEnabled  *bool     `json:"pull_request_enabled,omitempty"`
	HasPage             *bool     `json:"has_page,omitempty"`
	License             *string   `json:"license,omitempty"`
	OutSourced          *bool     `json:"outsourced,omitempty"`
	ProjectCreator      *string   `json:"project_creator,omitempty"`
	PushedAt            *string   `json:"pushed_at,omitempty"`
	CreatedAt           *string   `json:"created_at,omitempty"`
	UpdatedAt           *string   `json:"updated_at,omitempty"`
	Parent              *string   `json:"parent,omitempty"` // TBD
	Paas                *string   `json:"paas,omitempty"`   // TBD
	AssigneesNumber     *int      `json:"assignees_number,omitempty"`
	TestersNumber       *int      `json:"testers_number,omitempty"`
	Status              *string   `json:"status,omitempty"`
	Programs            *[]string `json:"programs,omitempty"`
	Enterprise          *string   `json:"enterprise,omitempty"`
	ProjectLabels       *[]Label  `json:"project_labels,omitempty"`
	IssueTemplateSource *string   `json:"issue_template_source,omitempty"`

	Members *[]string `json:"members,omitempty"`

	Namespace *Namespace     `json:"namespace,omitempty"`
	Owner     *UserProfile   `json:"owner,omitempty"`
	Assigner  *UserProfile   `json:"assigner,omitempty"`
	Assignee  *[]UserProfile `json:"assignee,omitempty"`
	Testers   *[]UserProfile `json:"testers,omitempty"`
}

type MyRepositoriesService struct {
	*github.RepositoriesService
	client *CustomClient
}

func (s *MyRepositoriesService) Get(ctx context.Context, owner, repo string) (*Repository, *github.Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repository := new(Repository)
	resp, err := s.client.Do(ctx, req, repository)
	if err != nil {
		return nil, resp, err
	}

	return repository, resp, nil
}

type CompareOptions struct {
	Owner      string
	Repository string
	Base       string
	Head       string
	Straight   bool
	Suffix     string
}

func (o *CompareOptions) Validate() error {
	if o.Owner == "" {
		return fmt.Errorf("missing owner")
	}
	if o.Repository == "" {
		return fmt.Errorf("missing repository")
	}
	if o.Base == "" {
		return fmt.Errorf("missing base")
	}
	if o.Head == "" {
		return fmt.Errorf("missing head")
	}
	return nil
}

func (o *CompareOptions) MakeRequestData() map[string]string {
	return nil
}

func (o *CompareOptions) ConstructEndpointUrl() *url.URL {
	endpointUrl, _ := url.Parse(
		fmt.Sprintf("repos/%v/%v/compare/%v...%v", o.Owner, o.Repository, o.Base, o.Head))
	q := endpointUrl.Query()
	if o.Straight {
		q.Set("straight", "true")
	}
	if o.Suffix != "" {
		q.Set("suffix", o.Suffix)
	}
	return endpointUrl
}

func (o *CompareOptions) RequestMethod() string {
	return http.MethodGet
}

type File struct {
	Sha        *string `json:"sha,omitempty"`
	Filename   *string `json:"filename,omitempty"`
	Status     *string `json:"status,omitempty"`
	Additions  *int    `json:"additions,omitempty"`
	Deletions  *int    `json:"deletions,omitempty"`
	Changes    *int    `json:"changes,omitempty"`
	BlobURL    *string `json:"blob_url,omitempty"`
	RawURL     *string `json:"raw_url,omitempty"`
	ContentURL *string `json:"content_url,omitempty"`
	Patch      *string `json:"patch,omitempty"`
	Truncated  *bool   `json:"truncated,omitempty"`
}

type CompareResult struct {
	BaseCommit      *Commit   `json:"base_commit,omitempty"`
	MergeBaseCommit *Commit   `json:"merge_base_commit,omitempty"`
	Commits         *[]Commit `json:"commits,omitempty"`
	Files           *[]File   `json:"files,omitempty"`
}

// Compare compares changes between specified head and base branches.
// Reference: https://gitee.com/api/v5/swagger#/getV5ReposOwnerRepoCompareBase...Head
func (s *MyRepositoriesService) Compare(ctx context.Context, opts *CompareOptions) (*CompareResult, *github.Response, error) {
	endpointUrl := opts.ConstructEndpointUrl()
	requestMethod := opts.RequestMethod()
	req, err := s.client.NewRequest(requestMethod, endpointUrl.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(CompareResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
