package gitee

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
)

type RepositoryFile struct {
	FileName string
	Content  string
	Length   int64
}

type RepositoryFilesService struct {
	client *CustomClient
}

// GetRawFile fetches a file of a specific revision from repository specified by argument owner and repo.
// Argument revision can be a branch, tag and commit, and even an empty string which means fetch file from the default branch.
// Reference: https://gitee.com/api/v5/swagger#/getV5ReposOwnerRepoRawPath
func (s *RepositoryFilesService) GetRawFile(ctx context.Context, owner, repo, filename, revision string) (*RepositoryFile, *github.Response, error) {
	queryParams := "access_token=" + url.QueryEscape(s.client.accessToken)
	if revision != "" {
		queryParams = fmt.Sprintf("%v&ref=%v", queryParams, url.QueryEscape(revision))
	}
	endpointUrl := fmt.Sprintf("repos/%v/%v/raw/%v?%v", owner, repo, filename, queryParams)
	req, err := s.client.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	buf := new(strings.Builder)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return nil, resp, err
	}
	var headerVal string
	headerVal = resp.Header.Get("content-disposition")
	fileName := strings.TrimPrefix(headerVal, "attachment; ")
	fileName = strings.Split(fileName, "=")[1]
	headerVal = resp.Header.Get("content-length")
	contentLength, err := strconv.ParseInt(headerVal, 10, 64)
	if err != nil {
		return nil, resp, errors.New("invalid content-length: " + headerVal)
	}
	repoFile := &RepositoryFile{
		FileName: fileName,
		Content:  buf.String(),
		Length:   contentLength,
	}
	return repoFile, resp, nil
}

type SlimLinks struct {
	Self string `json:"self,omitempty"`
	Html string `json:"html,omitempty"`
}

type Person struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Date  *string `json:"date,omitempty"`
}

type Tree struct {
	Sha *string `json:"sha,omitempty"`
	Url *string `json:"url,omitempty"`
}

type Content struct {
	Name    *string    `json:"name,omitempty"`
	Path    *string    `json:"path,omitempty"`
	Size    int64      `json:"size,omitempty"`
	Sha     *string    `json:"sha,omitempty"`
	Type    *string    `json:"type,omitempty"`
	URL     *string    `json:"url,omitempty"`
	HtmlURL *string    `json:"html_url,omitempty"`
	Links   *SlimLinks `json:"_links,omitempty"`
}

type CreationCommit struct {
	Sha       *string `json:"sha,omitempty"`
	Message   *string `json:"message,omitempty"`
	Author    *Person `json:"author,omitempty"`
	Committer *Person `json:"committer,omitempty"`
	Tree      *Tree   `json:"tree,omitempty"`
}

type FileOperationResult struct {
	Content *Content        `json:"content,omitempty"`
	Commit  *CreationCommit `json:"commit,omitempty"`
}

type FileOperationOptions struct {
	Owner          string
	Repository     string
	Path           string
	Branch         string
	Message        string
	CommitterName  string
	CommitterEmail string
	AuthorName     string
	AuthorEmail    string
}

func (o *FileOperationOptions) Validate() error {
	if o.Owner == "" {
		return fmt.Errorf("missing owner")
	}
	if o.Repository == "" {
		return fmt.Errorf("missing repository")
	}
	if o.Path == "" {
		return fmt.Errorf("missing path")
	}
	if o.Message == "" {
		return fmt.Errorf("missing commit message")
	}
	return nil
}

func (o *FileOperationOptions) MakeRequestData() map[string]string {
	data := map[string]string{
		"message": o.Message,
	}
	if o.Branch != "" {
		data["branch"] = o.Branch
	}
	if o.CommitterName != "" {
		data["committer[name]"] = o.CommitterName
	}
	if o.CommitterEmail != "" {
		data["committer[email]"] = o.CommitterEmail
	}
	if o.AuthorName != "" {
		data["author[name]"] = o.AuthorName
	}
	if o.AuthorEmail != "" {
		data["author[email]"] = o.AuthorEmail
	}
	return data
}

func (o *FileOperationOptions) ConstructEndpointUrl() string {
	return fmt.Sprintf("repos/%s/%s/contents/%s", o.Owner, o.Repository, url.PathEscape(o.Path))
}

func (o *FileOperationOptions) RequestMethod() string {
	return http.MethodPost
}

type FileCreationOptions struct {
	*FileOperationOptions

	Content string
}

func (o *FileCreationOptions) Validate() error {
	if err := o.FileOperationOptions.Validate(); err != nil {
		return err
	}
	if o.Content == "" {
		return fmt.Errorf("missing content")
	}
	return nil
}

func (o *FileCreationOptions) MakeRequestData() map[string]string {
	data := o.FileOperationOptions.MakeRequestData()
	contentInBase64 := base64.StdEncoding.EncodeToString([]byte(o.Content))
	data["content"] = contentInBase64
	return data
}

// CreateFile creates a file at specified revision.
// Reference: https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoContentsPath
func (s *RepositoryFilesService) CreateFile(ctx context.Context, opts *FileCreationOptions) (*FileOperationResult, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	return s.operateFile(ctx, "create", opts)
}

type FileUpdateOptions struct {
	*FileCreationOptions

	Sha string
}

func (o *FileUpdateOptions) Validate() error {
	if err := o.FileCreationOptions.Validate(); err != nil {
		return err
	}
	if o.Sha == "" {
		return fmt.Errorf("missing sha of the original file content")
	}
	return nil
}

func (o *FileUpdateOptions) MakeRequestData() map[string]string {
	data := o.FileCreationOptions.MakeRequestData()
	data["sha"] = o.Sha
	return data
}

func (o *FileUpdateOptions) RequestMethod() string {
	return http.MethodPut
}

// UpdateFile updates a file.
func (s *RepositoryFilesService) UpdateFile(ctx context.Context, opts *FileUpdateOptions) (*FileOperationResult, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	return s.operateFile(ctx, "update", opts)
}

type FileDeletionOptions struct {
	*FileOperationOptions

	Sha string
}

func (o *FileDeletionOptions) Validate() error {
	if err := o.FileOperationOptions.Validate(); err != nil {
		return err
	}
	if o.Sha == "" {
		return fmt.Errorf("missing sha of the original file content")
	}
	return nil
}

func (o *FileDeletionOptions) MakeRequestData() map[string]string {
	data := o.FileOperationOptions.MakeRequestData()
	data["sha"] = o.Sha
	return data
}

func (o *FileDeletionOptions) RequestMethod() string {
	return http.MethodDelete
}

// DeleteFile deletes a file.
func (s *RepositoryFilesService) DeleteFile(ctx context.Context, opts *FileDeletionOptions) (*FileOperationResult, *github.Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	return s.operateFile(ctx, "delete", opts)
}

func (s *RepositoryFilesService) operateFile(ctx context.Context, operation string, opts any) (*FileOperationResult, *github.Response, error) {
	operationRequest, ok := opts.(RequestDataGenerator)
	if !ok {
		return nil, nil, fmt.Errorf("unknown operation '%v'", operation)
	}

	endpointUrl := operationRequest.ConstructEndpointUrl()
	data := operationRequest.MakeRequestData()
	requestMethod := operationRequest.RequestMethod()

	data["access_token"] = s.client.accessToken
	req, err := s.client.NewRequest(requestMethod, endpointUrl, data)
	if err != nil {
		return nil, nil, err
	}

	fci := new(FileOperationResult)
	resp, err := s.client.Do(ctx, req, fci)
	if err != nil {
		return nil, resp, err
	}
	return fci, resp, nil
}
