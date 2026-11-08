package gitee

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
)

type PullRequestsService struct {
	client *CustomClient
}

type Reference struct {
	Label      *string         `json:"label,omitempty"`
	Ref        *string         `json:"ref,omitempty"`
	Sha        *string         `json:"sha,omitempty"`
	User       *User           `json:"user,omitempty"`
	Repository *SlimRepository `json:"repo,omitempty"`
}

type PullRequestLinks struct {
	Self           *Link `json:"self,omitempty"`
	Html           *Link `json:"html,omitempty"`
	Issue          *Link `json:"issue,omitempty"`
	Comments       *Link `json:"comments,omitempty"`
	ReviewComments *Link `json:"review_comments,omitempty"`
	ReviewComment  *Link `json:"review_comment,omitempty"`
	Commits        *Link `json:"commits,omitempty"`
}

type ReferencePullRequest struct {
	Number *string `json:"number,omitempty"`
	Title  *string `json:"title,omitempty"`
	State  *string `json:"state,omitempty"`
}

type PullRequest struct {
	Id                *int64                  `json:"id,omitempty"`
	Title             *string                 `json:"title,omitempty"`
	Content           *string                 `json:"content,omitempty"`
	Url               *string                 `json:"url,omitempty"`
	HtmlUrl           *string                 `json:"html_url,omitempty"`
	DiffUrl           *string                 `json:"diff_url,omitempty"`
	PatchUrl          *string                 `json:"patch_url,omitempty"`
	IssueUrl          *string                 `json:"issue_url,omitempty"`
	CommitsUrl        *string                 `json:"commits_url,omitempty"`
	ReviewCommentsUrl *string                 `json:"review_comments_url,omitempty"`
	ReviewCommentUrl  *string                 `json:"review_comment_url,omitempty"`
	CommentsUrl       *string                 `json:"comments_url,omitempty"`
	Number            *int                    `json:"number,omitempty"`
	CloseRelatedIssue *int                    `json:"close_related_issue,omitempty"`
	PruneBranch       *bool                   `json:"prune_branch,omitempty"`
	State             *string                 `json:"state,omitempty"`
	AssigneesNumber   *int                    `json:"assignees_number,omitempty"`
	TestersNumber     *int                    `json:"testers_number,omitempty"`
	Locked            *bool                   `json:"locked,omitempty"`
	CreatedAt         *string                 `json:"created_at,omitempty"`
	UpdatedAt         *string                 `json:"updated_at,omitempty"`
	ClosedAt          *string                 `json:"closed_at,omitempty"`
	Draft             *bool                   `json:"draft,omitempty"`
	MergedAt          *string                 `json:"merged_at,omitempty"`
	SecurityHole      *bool                   `json:"security_hole,omitempty"`
	Mergeable         *bool                   `json:"mergeable,omitempty"`
	CanMergeCheck     *bool                   `json:"can_merge_check,omitempty"`
	Assignees         *[]User                 `json:"assignees,omitempty"`
	Testers           *[]User                 `json:"testers,omitempty"`
	Labels            *[]Label                `json:"labels,omitempty"`
	User              *User                   `json:"user,omitempty"`
	Links             *PullRequestLinks       `json:"_links,omitempty"`
	RefPullRequests   *[]ReferencePullRequest `json:"ref_pull_requests,omitempty"`
	Base              *Reference              `json:"base,omitempty"`
	Head              *Reference              `json:"head,omitempty"`
	Milestone         *Milestone              `json:"milestone,omitempty"`
}

type PullRequestsListOptions struct {
	// FIXME: define constants
	State string // State, open, closed, merged and all.
	Head  string // Source branch name in forms branch or username:branch
	Base  string // Target branch name in forms branch or username:branch
	// FIXME: define constants
	Sort            string // Sort by this field. Valid names are created, updated, popularity and long-running. If omitted, sort by creation time.
	Since           string // List pull requests since this datetime which is in ISO8601 format, e.g. 2026-01-01T12:00:00+08:00
	Direction       string // Sort in ASC or DESC. Valid values are asc and desc.
	MilestoneNumber string // Empty means omitted. It will be converted into an integer during method execution.
	Labels          []string
	Page            string // Empty means omitted. It will be converted into an integer during method execution.
	PerPage         string // Empty means omitted. It will be converted into an integer during method execution.
	Author          string
	Assignee        string
	Tester          string
}

func (o *PullRequestsListOptions) encodeToQueryString() (url.Values, error) {
	params, _ := url.ParseQuery("")
	if o.State != "" {
		params.Set("state", o.State)
	}
	if o.Head != "" {
		params.Set("head", o.Head)
	}
	if o.Base != "" {
		params.Set("base", o.Base)
	}
	if o.Sort != "" {
		params.Set("sort", o.Sort)
	}
	if o.Since != "" {
		params.Set("since", o.Since)
	}
	if o.Direction != "" {
		params.Set("direction", o.Direction)
	}
	if o.MilestoneNumber != "" {
		_, err := strconv.Atoi(o.MilestoneNumber)
		if err != nil {
			return nil, err
		}
		params.Set("milestone_number", o.MilestoneNumber)
	}
	joinedLabels := strings.Join(o.Labels, ",")
	if joinedLabels != "" {
		params.Set("labels", joinedLabels)
	}
	if o.Page != "" {
		_, err := strconv.Atoi(o.Page)
		if err != nil {
			return nil, err
		}
		params.Set("page", o.Page)
	}
	if o.PerPage != "" {
		_, err := strconv.Atoi(o.PerPage)
		if err != nil {
			return nil, err
		}
		params.Set("per_page", o.PerPage)
	}
	if o.Author != "" {
		params.Set("author", o.Author)
	}
	if o.Assignee != "" {
		params.Set("assignee", o.Assignee)
	}
	if o.Tester != "" {
		params.Set("tester", o.Tester)
	}
	return params, nil
}

func (s *PullRequestsService) List(ctx context.Context, owner, repo string, opts *PullRequestsListOptions) (*[]PullRequest, *github.Response, error) {
	endpointUrl, err := url.Parse(fmt.Sprintf("repos/%v/%v/pulls", owner, repo))
	if err != nil {
		return nil, nil, err
	}

	queryParams, err := opts.encodeToQueryString()
	if err != nil {
		return nil, nil, err
	}
	queryParams.Set("access_token", s.client.accessToken)
	endpointUrl.RawQuery = queryParams.Encode()

	req, err := s.client.NewRequest(http.MethodGet, endpointUrl.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	pullRequests := new([]PullRequest)
	resp, err := s.client.Do(ctx, req, pullRequests)
	if err != nil {
		return nil, resp, err
	}
	return pullRequests, resp, nil
}

type PullRequestCreationOptions struct {
	Owner                 string
	Repository            string
	Title                 string
	Head                  string
	Base                  string
	Body                  string
	MilestoneNumber       int
	Labels                []string
	Issue                 string
	Assignees             []string
	Testers               []string
	AssigneesNumber       int
	TestersNumber         int
	RefPullRequestNumbers []int
	PruneSourceBranch     bool
	CloseRelatedIssue     bool
	Draft                 bool
	Squash                bool
	SecurityHole          bool
}

func (o *PullRequestCreationOptions) Validate() error {
	if o.Owner == "" {
		return fmt.Errorf("missing owner")
	}
	if o.Repository == "" {
		return fmt.Errorf("missing repository")
	}
	if o.Title == "" {
		return fmt.Errorf("missing title")
	}
	if o.Head == "" {
		return fmt.Errorf("missing head")
	}
	if o.Base == "" {
		return fmt.Errorf("missing base")
	}
	return nil
}

func (o *PullRequestCreationOptions) ConstructEndpointUrl() *url.URL {
	endpointUrl, _ := url.Parse(fmt.Sprintf("repos/%v/%v/pulls", o.Owner, o.Repository))
	return endpointUrl
}

func (o *PullRequestCreationOptions) MakeRequestData() map[string]string {
	data := map[string]string{}
	if o.Body != "" {
		data["body"] = o.Body
	}
	if o.MilestoneNumber > 0 {
		data["milestone_number"] = strconv.FormatInt(int64(o.MilestoneNumber), 10)
	}
	if len(o.Labels) > 0 {
		data["labels"] = strings.Join(o.Labels, ",")
	}
	if o.Issue != "" {
		data["issue"] = o.Issue
	}
	if len(o.Assignees) > 0 {
		data["assignees"] = strings.Join(o.Assignees, ",")
	}
	if len(o.Testers) > 0 {
		data["testers"] = strings.Join(o.Testers, ",")
	}
	if o.AssigneesNumber > 0 {
		data["assignees_number"] = strconv.FormatInt(int64(o.AssigneesNumber), 10)
	}
	if o.TestersNumber > 0 {
		data["testers_number"] = strconv.FormatInt(int64(o.TestersNumber), 10)
	}
	if len(o.RefPullRequestNumbers) > 0 {
		var numbers []string
		for _, v := range o.RefPullRequestNumbers {
			numbers = append(numbers, strconv.FormatInt(int64(v), 10))
		}
		data["ref_pull_request_numbers"] = strings.Join(numbers, ",")
	}
	if o.PruneSourceBranch {
		data["prune_source_branch"] = "true"
	}
	if o.CloseRelatedIssue {
		data["close_related_issue"] = "true"
	}
	if o.Draft {
		data["draft"] = "true"
	}
	if o.Squash {
		data["squash"] = "true"
	}
	if o.SecurityHole {
		data["security_hole"] = "true"
	}
	return data
}

func (o *PullRequestCreationOptions) RequestMethod() string {
	return http.MethodPost
}

func (s *PullRequestsService) Create(ctx context.Context, opts *PullRequestCreationOptions) (*PullRequest, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	endpointUrl := opts.ConstructEndpointUrl()
	queryParams := endpointUrl.Query()
	queryParams.Set("access_token", s.client.accessToken)
	endpointUrl.RawQuery = queryParams.Encode()

	requestData := opts.MakeRequestData()

	req, err := s.client.NewRequest(opts.RequestMethod(), endpointUrl.String(), requestData)
	if err != nil {
		return nil, nil, err
	}
	pullRequest := new(PullRequest)
	resp, err := s.client.Do(ctx, req, pullRequest)
	if err != nil {
		return nil, resp, err
	}
	return pullRequest, resp, nil
}
