package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitee-play.com/build-service/pkg/git/gitee"
	go_gitee "gitee-play.com/pkg/gitee"
)

const BuildServiceUserAgent = "Gitee-play/0.0 (main)"

func main() {
	token := os.Getenv("GITEE_TOKEN")
	repoUrl := os.Getenv("REPO_URL")

	if token == "" {
		panic("GITEE_TOKEN environment variable not set")
	}

	if repoUrl == "" {
		panic("REPO_URL environment variable not set")
	}

	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		panic(err)
	}

	pos := strings.LastIndex(parsedUrl.Path, "/")
	if pos == -1 {
		panic(fmt.Errorf("URL path does not have /."))
	}

	owner := parsedUrl.Path[0:pos]
	repoName := parsedUrl.Path[pos+1 : len(parsedUrl.Path)]

	giteeClient := go_gitee.NewClient(go_gitee.Opt{})
	ctx := context.TODO()
	repo, _, err := giteeClient.RepositoriesService.Get(ctx, owner, repoName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Is public:", !*repo.Private)
	fmt.Println("License:", *repo.License)

	client := gitee.NewGiteeClient(token)

	client.EnsurePaCMergeRequest(repoUrl, nil)
}
