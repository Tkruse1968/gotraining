package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/google/go-github/v39/github"
    "golang.org/x/oauth2"
)

func main() {
    // Example usage
    err := copyTopicToCustomProperties("https://github.example.com/api/v3/", "owner", "repo", "new-topic")
    if err != nil {
        log.Fatalf("Error copying topic to custom properties: %v", err)
    }
}

func copyTopicToCustomProperties(baseURL, owner, repo, newTopic string) error {
    ctx := context.Background()
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        return fmt.Errorf("GITHUB_TOKEN environment variable not set")
    }

    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    
    // Create a new GitHub client with the custom API endpoint
    client, err := github.NewEnterpriseClient(baseURL, baseURL, tc)
    if err != nil {
        return fmt.Errorf("error creating GitHub client: %v", err)
    }

    // Fetch current topics with pagination
    var allTopics []string
    opt := &github.ListOptions{PerPage: 100}
    for {
        topics, resp, err := client.Repositories.ListAllTopics(ctx, owner, repo, opt)
        if err != nil {
            return fmt.Errorf("error fetching topics: %v", err)
        }
        allTopics = append(allTopics, topics...)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }

    // Add new topic to custom properties (assuming custom properties are stored in a specific format)
    customProperties := map[string]string{
        "custom_topic": newTopic,
    }

    // Update repository with new custom properties
    repoRequest := &github.Repository{
        Topics: append(allTopics, newTopic),
    }
    _, _, err = client.Repositories.Edit(ctx, owner, repo, repoRequest)
    if err != nil {
        return fmt.Errorf("error updating repository: %v", err)
    }

    log.Printf("Successfully copied topic '%s' to custom properties", newTopic)
    return nil
}
