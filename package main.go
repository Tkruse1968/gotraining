package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"

    go get "github.com/google/go-github/v39/github"
    go get "golang.org/x/oauth2"
)

func main() {
    // Replace with your GitHub access token
    token := "YOUR_GITHUB_ACCESS_TOKEN"

    // Replace with your GitHub Enterprise organization name
    org := "YOUR_GITHUB_ENTERPRISE_ORG"

    // Replace with your internal address
    internalAddress := "http://your-internal-address"

    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)

    // Create a custom HTTP client with a custom transport
    client := github.NewClient(&http.Client{
        Transport: &http.Transport{
            DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                // Allow connections to internal addresses
                if addr == internalAddress {
                    return net.Dial(network, addr)
                }
                // For other addresses, use the default dialer
                return (&net.Dialer{}).DialContext(ctx, network, addr)
            },
        },
    })

    // List all repositories for the organization
    opt := &github.RepositoryListByOrgOptions{
        ListOptions: github.ListOptions{PerPage: 10}, // Adjust the number of repositories per page as needed
    }
    repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)
    if err != nil {
        log.Fatal(err)
    }

    // Iterate over the repositories
    for _, repo := range repos {
        // Get the repository topics
        topics, _, err := client.Repositories.ListAllTopics(ctx, org, *repo.Name)
        if err != nil {
            log.Fatal(err)
        }

        // Copy the topics to custom properties
        for _, topic := range topics {
            fmt.Printf("Copying topic: %s\n", *topic.Name)
            // Add your custom logic here to copy the topic to custom properties
            copyToCustomProperties(*repo.Name, *topic.Name)
        }
    }
}

func copyToCustomProperties(repo, topic string) {
    // Add your logic here to copy the topic to custom properties in GitHub Enterprise
    // This function will be triggered by an event of adding a new topic to a GitHub Enterprise repository
    // You can use the GitHub Enterprise API to update the custom properties
    // Make sure to handle any errors and log them if necessary
    fmt.Printf("Copying topic to custom properties for repository %s: %s\n", repo, topic)
}