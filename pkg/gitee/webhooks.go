package gitee

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-github/v45/github"
)

type Webhook struct {
	Id                  *int64  `json:"id,omitempty"`
	CreatedAt           *string `json:"created_at,omitempty"`
	IssuesEvents        *bool   `json:"issues_events,omitempty"`
	MergeRequestsEvents *bool   `json:"merge_requests_events,omitempty"`
	NoteEvents          *bool   `json:"note_events,omitempty"`
	Password            *string `json:"password,omitempty"`
	ProjectId           *int64  `json:"project_id,omitempty"`
	PushEvents          *bool   `json:"push_events,omitempty"`
	Result              *string `json:"result,omitempty"`
	ResultCode          *int64  `json:"result_code,omitempty"`
	TagPushEvents       *bool   `json:"tag_push_events,omitempty"`
	Title               *string `json:"title,omitempty"`
	URL                 *string `json:"url,omitempty"`
}

type WebhooksService struct {
	*github.RepositoriesService
	client *CustomClient
}

type WebhookListOptions struct {
	Owner      string
	Repository string
	Page       int
	PerPage    int
}

func (o *WebhookListOptions) Validate() error {
	if o.Owner == "" {
		return fmt.Errorf("missing owner")
	}
	if o.Repository == "" {
		return fmt.Errorf("missing repository")
	}
	return nil
}

func (o *WebhookListOptions) ConstructEndpointUrl() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("repos/%v/%v/hooks", o.Owner, o.Repository))
	param := u.Query()
	var paramVal string
	if o.Page > 0 {
		paramVal = strconv.FormatInt(int64(o.Page), 10)
		param.Set("page", paramVal)
	}
	if o.Page > 0 {
		paramVal = strconv.FormatInt(int64(o.Page), 10)
		param.Set("page", paramVal)
	}
	u.RawQuery = param.Encode()
	return u
}

func (o *WebhookListOptions) RequestMethod() string {
	return http.MethodGet
}

func (s *WebhooksService) List(ctx context.Context, opts *WebhookListOptions) (*Webhook, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	endpointUrl := opts.ConstructEndpointUrl()
	param := endpointUrl.Query()
	param.Set("access_token", s.client.accessToken)
	endpointUrl.RawQuery = param.Encode()
	req, err := s.client.NewRequest(opts.RequestMethod(), endpointUrl.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	webhook := new(Webhook)
	resp, err := s.client.Do(ctx, req, webhook)
	if err != nil {
		return nil, resp, err
	}

	return webhook, resp, nil
}

type WebhookCreationOptions struct {
	Owner               string
	Repository          string
	URL                 string
	Title               string
	EncryptionType      string
	Password            string
	PushEvents          bool
	TagPushEvents       bool
	IssuesEvents        bool
	NoteEvents          bool
	MergeRequestsEvents bool
}

const MaxWebhookTitleLength = 191

func (o *WebhookCreationOptions) Validate() error {
	if o.Owner == "" {
		return fmt.Errorf("missing owner")
	}
	if o.Repository == "" {
		return fmt.Errorf("missing repository")
	}
	if o.URL == "" {
		return fmt.Errorf("missing url")
	}
	if o.EncryptionType != "0" && o.EncryptionType != "1" {
		return fmt.Errorf("encryption type '%v' is invalid", o.EncryptionType)
	}
	if len(o.Title) > MaxWebhookTitleLength {
		return fmt.Errorf("title is too long which exceeds 191 characters")
	}
	return nil
}

func (o *WebhookCreationOptions) ConstructEndpointUrl() *url.URL {
	u, _ := url.Parse(fmt.Sprintf("repos/%v/%v/hooks", o.Owner, o.Repository))
	return u
}

func (o *WebhookCreationOptions) RequestMethod() string {
	return http.MethodPost
}

func (o *WebhookCreationOptions) MakeRequestData() map[string]string {
	data := map[string]string{"url": o.URL}
	if o.Title != "" {
		data["title"] = o.Title
	}
	if o.Password != "" {
		data["password"] = o.Password
	}
	data["encryption_type"] = o.EncryptionType
	data["push_events"] = strconv.FormatBool(o.PushEvents)
	data["tag_push_events"] = strconv.FormatBool(o.TagPushEvents)
	data["issues_events"] = strconv.FormatBool(o.IssuesEvents)
	data["note_events"] = strconv.FormatBool(o.NoteEvents)
	data["merge_requests_events"] = strconv.FormatBool(o.MergeRequestsEvents)
	return data
}

func (s *WebhooksService) Create(ctx context.Context, opts *WebhookCreationOptions) (*Webhook, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	endpointUrl := opts.ConstructEndpointUrl()
	requestData := opts.MakeRequestData()
	requestData["access_token"] = s.client.accessToken

	req, err := s.client.NewRequest(opts.RequestMethod(), endpointUrl.String(), requestData)
	if err != nil {
		return nil, nil, err
	}
	webhook := new(Webhook)
	resp, err := s.client.Do(ctx, req, webhook)
	if err != nil {
		return nil, resp, err
	}
	return webhook, resp, nil
}
