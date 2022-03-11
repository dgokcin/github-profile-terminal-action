package profile

import (
    "context"
    "github.com/google/go-github/v43/github"
    "github.com/liamg/github-profile-magic-action/config"
    "golang.org/x/oauth2"
    "net/http"
    "os"
)

func newGithubClient(conf *config.Config) *github.Client {
    var tc *http.Client
    token := conf.Context.Token
    if token == "" {
        token = os.Getenv("GITHUB_TOKEN")
    }
    if token != "" {
        ctx := context.Background()
        ts := oauth2.StaticTokenSource(
            &oauth2.Token{AccessToken: token},
        )
        tc = oauth2.NewClient(ctx, ts)
    }
    return github.NewClient(tc)
}
