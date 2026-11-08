package main

import (
	"context"
	"fmt"

	"gitee-play.com/build-service/pkg/git/gitee"
	go_gitee "gitee-play.com/pkg/gitee"
)

const BuildServiceUserAgent = "Gitee-play/0.0 (main)"

func main() {
	token := ""

	giteeClient := go_gitee.NewClient(go_gitee.Opt{})
	ctx := context.TODO()
	repo, _, err := giteeClient.RepositoriesService.Get(ctx, "qicx", "open-ai-window")
	if err != nil {
		panic(err)
	}
	fmt.Println("Is public:", !*repo.Private)
	fmt.Println("License:", *repo.License)

	client := gitee.NewGiteeClient(token)

	client.EnsurePaCMergeRequest("https://gitee.com/qicx/open-ai-window", nil)
}
