package gitee

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitee-play.com/pkg/gitee"
)

func getOwnerAndRepoFromUrl(repoUrl string) (owner string, repository string) {
	// https://github.com/owner/repository
	gitSourceUrlParts := strings.Split(strings.TrimSuffix(repoUrl, ".git"), "/")
	owner = gitSourceUrlParts[3]
	repository = gitSourceUrlParts[4]
	return owner, repository
}

type GiteeClient struct {
	ctx    context.Context
	client *gitee.CustomClient
}

type MergeRequestData struct {
	CommitMessage string
	SignedOff     bool
	BranchName    string
}

type MergeRequest struct {
	Id        int64
	CreatedAt *time.Time
	WebUrl    string
	Title     string
}

// GitProviderClient interface implementation

func (c *GiteeClient) EnsurePaCMergeRequest(repoUrl string, d *MergeRequestData) (webUrl string, err error) {
	prs, resp, err := c.client.PullRequestsService.List(c.ctx, "qicx", "open-ai-window", &gitee.PullRequestsListOptions{})
	fmt.Println(prs)
	fmt.Println("---")
	fmt.Println(resp)
	fmt.Println("---")
	fmt.Println(err)
	fmt.Println("---")
	for _, v := range *prs {
		fmt.Println(*v.Title)
	}
	return "", nil
}

func (c *GiteeClient) UndoPaCMergeRequest(repoUrl string, d *MergeRequestData) (webUrl string, err error) {
	return "", nil
}

func (c *GiteeClient) FindUnmergedPaCMergeRequest(repoUrl string, d *MergeRequestData) (*MergeRequest, error) {
	return nil, nil
}

func (c *GiteeClient) SetupPaCWebhook(repoUrl, webhookUrl, webhookSecret string) error {
	return nil
}

func (c *GiteeClient) DeletePaCWebhook(repoUrl, webhookUrl string) error {
	return nil
}

func (c *GiteeClient) GetDefaultBranchWithChecks(repoUrl string) (string, error) {
	return "", nil
}

func (c *GiteeClient) DeleteBranch(repoUrl, branchName string) (bool, error) {
	return false, errors.New("haven't found out a Gitee API for branch deletion")
}

func (c *GiteeClient) GetBranchSha(repoUrl, branchName string) (string, error) {
	owner, repository := getOwnerAndRepoFromUrl(repoUrl)

	// If branch is not specified, use default branch
	if branchName == "" {
		defaultBranchName, err := c.getDefaultBranch(owner, repository)
		if err != nil {
			return "", err
		}
		branchName = defaultBranchName
	}

	branch, err := c.getBranch(owner, repository, branchName)
	if err != nil {
		return "", err
	}
	if branch == nil {
		return "", nil
	}
	return *branch.Commit.Sha, nil
}

func (c *GiteeClient) IsFileExist(repoUrl, branchName, filePath string) (bool, error) {
	return false, nil
}

func (c *GiteeClient) DownloadFileContent(repoUrl, branchName, filePath string) ([]byte, error) {
	repoFile, _, err := c.client.RepositoryFilesService.GetRawFile(c.ctx, "qicx", "open-ai-window", filePath, branchName)
	if err != nil {
		return nil, err
	}
	fmt.Println("Filename:", repoFile.FileName)
	fmt.Println("Content length:", repoFile.Length)
	fmt.Println("Content:", repoFile.Content)
	return nil, nil
}

func (c *GiteeClient) IsRepositoryPublic(repoUrl string) (bool, error) {
	owner, repository := getOwnerAndRepoFromUrl(repoUrl)
	repoInfo, err := c.getRepository(owner, repository)
	if err != nil {
		return false, err
	}
	if repoInfo == nil {
		// There is no difference between private and not found for GitHub unless the user is the owner.
		return false, nil
	}
	return !*repoInfo.Private, nil
}

func (c *GiteeClient) GetBrowseRepositoryAtShaLink(repoUrl, sha string) string {
	return ""
}

func (c *GiteeClient) GetAppUserId(userName string) (int64, error) {
	return 0, errors.New("not use application for Gitee")
}

// end of GitProviderClient interface implementation

func NewGiteeClient(accessToken string) *GiteeClient {
	client := &GiteeClient{ctx: context.TODO()}
	client.client = gitee.NewClient(gitee.Opt{AccessToken: accessToken})
	client.client.UserAgent = "SealVis/0.0 (build-service)"

	return client
}
