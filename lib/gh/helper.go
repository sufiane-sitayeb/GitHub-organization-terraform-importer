package gh

import (
    "context"
    "github.com/google/go-github/v29/github"
    "golang.org/x/oauth2"
)

type Helper struct {
    Client*     github.Client
    Ctx         context.Context
    Org         string
    Token       string

    Users       UserCollection
    Teams       TeamCollection
    Repos       RepoCollection
}

var Github Helper

func Initialize(org string, token string) {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    Github = Helper{
        client,
        ctx,
        org,
        token,
        UserCollection{},
        TeamCollection{},
        RepoCollection{},
    }
}
