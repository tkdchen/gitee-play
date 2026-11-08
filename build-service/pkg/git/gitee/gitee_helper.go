package gitee

import (
	"fmt"

	"gitee-play.com/pkg/gitee"
	"github.com/google/go-github/v45/github"
)

func (c *GiteeClient) getRepository(owner, name string) (*gitee.Repository, error) {
	repo, resp, err := c.client.RepositoriesService.Get(c.ctx, owner, name)
	if err == nil {
		return repo, nil
	}
	if resp == nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		// Do not treat 404 as failure operation. It just means the queried repository does not
		// exist which is represented as nil.
		return nil, nil
	}
	return nil, err
}

func (c *GiteeClient) getDefaultBranch(owner, repository string) (string, error) {
	repositoryInfo, resp, err := c.client.RepositoriesService.Get(c.ctx, owner, repository)
	if err != nil {
		//return "", refineGitHostingServiceError(resp.Response, err)
		// FIXME: enable above line when ready
		fmt.Println(resp)
		return "", err
	}
	if repositoryInfo == nil {
		return "", fmt.Errorf("repository info is empty in GitHub API response")
	}
	return *repositoryInfo.DefaultBranch, nil
}

func (c *GiteeClient) getBranch(owner, repository, branch string) (*gitee.Branch, error) {
	branchObj, resp, err := c.client.BranchesService.Get(c.ctx, owner, repository, branch)
	if resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		//return nil, refineGitHostingServiceError(resp.Response, err)
		// FIXME: restore the above line
		return nil, err
	}
	return branchObj, nil
}

func (c *GiteeClient) getWebhookByTargetUrl(owner, repository, webhookTargetUrl string) (*github.Hook, error) {
	// Suppose that the repository does not have more than 100 webhooks
	listOpts := &github.ListOptions{PerPage: 100}
	webhooks, resp, err := c.client.Repositories.ListHooks(c.ctx, owner, repository, listOpts)
	if err != nil {
		//return nil, refineGitHostingServiceError(resp.Response, err)
		// FIXME: restore the above line
		fmt.Println(resp)
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Config["url"] == webhookTargetUrl {
			return webhook, nil
		}
	}
	// Webhook with the given URL not found
	return nil, nil
}

func (c *GiteeClient) branchExist(owner, repository, branchName string) (bool, error) {
	_, resp, err := c.client.BranchesService.Get(c.ctx, owner, repository, branchName)
	if err == nil {
		return true, nil
	}
	if resp == nil {
		return false, err
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return false, err
}
