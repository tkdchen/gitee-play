package gitee

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-github/v45/github"
)

type Sha struct {
	Sha *string `json:"sha,omitempty"`
	Url *string `json:"url,omitempty"`
}

type UserContact struct {
	Name  *string `json:"name,omitempty"`
	Date  *string `json:"date,omitempty"`
	Email *string `json:"email,omitempty"`
}

type UserProfile struct {
	Id                *string `json:"id,omitempty"`
	Login             *string `json:"login,omitempty"`
	Name              *string `json:"name,omitempty"`
	AvatarURL         *string `json:"avatar_url,omitempty"`
	Url               *string `json:"url,omitempty"`
	HtmlURL           *string `json:"html_url,omitempty"`
	Remark            *string `json:"remark,omitempty"`
	FollowersURL      *string `json:"followers_url,omitempty"`
	FollowingURL      *string `json:"following_url,omitempty"`
	GistsURL          *string `json:"gists_url,omitempty"`
	StarredURL        *string `json:"starred_url,omitempty"`
	SubscriptionsURL  *string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  *string `json:"organizations_url,omitempty"`
	RepositoriesURL   *string `json:"repos_url,omitempty"`
	EventsURL         *string `json:"events_url,omitempty"`
	ReceivedEventsURL *string `json:"received_events_url,omitempty"`
	Type              *string `json:"type,omitempty"`
}

type NestedCommit struct {
	URL       *string      `json:"url,omitempty"`
	Message   *string      `json:"message,omitempty"`
	Author    *UserContact `json:"author,omitempty"`
	Committer *UserContact `json:"committer,omitempty"`
	Tree      *Sha         `json:"tree,omitempty"`
}

type Commit struct {
	Sha          *string       `json:"sha,omitempty"`
	URL          *string       `json:"url,omitempty"`
	HtmlURL      *string       `json:"html_url,omitempty"`
	CommentsURL  *string       `json:"comments_url,omitempty"`
	NestedCommit *NestedCommit `json:"commit,omitempty"`
	Author       *UserProfile  `json:"author,omitempty"`
	Committer    *UserProfile  `json:"committer,omitempty"`
	Parents      *[]Sha        `json:"parents,omitempty"`
}

type Branch struct {
	Name          *string    `json:"name,omitempty"`
	Commit        *Commit    `json:"commit,omitempty"`
	Links         *SlimLinks `json:"_links,omitempty"`
	Protected     *bool      `json:"protected,omitempty"`
	ProtectionURL *string    `json:"protection_url,omitempty"`
}

type BranchesService struct {
	client *CustomClient
}

// Reference: https://gitee.com/api/v5/swagger#/getV5ReposOwnerRepoBranchesBranch
func (s *BranchesService) Get(ctx context.Context, owner, repo, branch string) (*Branch, *github.Response, error) {
	endpointUrl := fmt.Sprintf("repos/%v/%v/branches/%v", owner, repo, branch)
	req, err := s.client.NewRequest("GET", endpointUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	branchObj := new(Branch)
	resp, err := s.client.Do(ctx, req, branchObj)
	if err != nil {
		return nil, resp, err
	}
	return branchObj, resp, nil
}

type BranchCreationOptions struct {
	Owner      string
	Repository string
	Reference  string
	BranchName string
}

func (o *BranchCreationOptions) ConstructEndpointUrl() *url.URL {
	endpointUrl, _ := url.Parse(fmt.Sprintf("repos/%v/%v/branches", o.Owner, o.Reference))
	return endpointUrl
}

func (o *BranchCreationOptions) MakeRequestData() map[string]string {
	return map[string]string{
		"refs":        o.Reference,
		"branch_name": o.BranchName,
	}
}

func (o *BranchCreationOptions) RequestMethod() string {
	return http.MethodPost
}

// Create creates a new branch based on reference specified by ref.
// Argument ref can be a branch name e.g. main.
// Reference: https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoBranches
func (s *BranchesService) Create(ctx context.Context, opts *BranchCreationOptions) (*Branch, *github.Response, error) {
	endpointUrl := opts.ConstructEndpointUrl()
	requestData := opts.MakeRequestData()
	requestData["access_token"] = s.client.accessToken

	req, err := s.client.NewRequest("GET", endpointUrl.String(), requestData)
	if err != nil {
		return nil, nil, err
	}
	branchObj := new(Branch)
	resp, err := s.client.Do(ctx, req, branchObj)
	if err != nil {
		return nil, resp, err
	}
	return branchObj, resp, nil
}

func (s *BranchesService) Delete(ctx context.Context, branch string) (*Branch, *github.Response, error) {
	return nil, nil, fmt.Errorf("gitee service does not have endpoint to delete a branch")
}
