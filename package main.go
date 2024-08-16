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

func copyTopicToCustomProperties(owner, repo, newTopic string) error {
    ctx := context.Background()
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        return fmt.Errorf("GITHUB_TOKEN environment variable not set")
    }

    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // Fetch current topics
    topics, _, err := client.Repositories.ListAllTopics(ctx, owner, repo)
    if err != nil {
        return fmt.Errorf("error fetching topics: %v", err)
    }

    // Add new topic to custom properties (assuming custom properties are stored in a specific format)
    customProperties := map[string]string{
        "custom_topic": newTopic,
    }

    // Update repository with new custom properties
    repoRequest := &github.Repository{
        Topics: append(topics, newTopic),
    }
    _, _, err = client.Repositories.Edit(ctx, owner, repo, repoRequest)
    if err != nil {
        return fmt.Errorf("error updating repository: %v", err)
    }

    log.Printf("Successfully copied topic '%s' to custom properties", newTopic)
    return nil
}
    fmt.Printf("Copying topic to custom properties for repository %s: %s\n", repo, topic)
}