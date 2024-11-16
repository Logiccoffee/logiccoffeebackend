package ghupload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-github/v59/github"

	"golang.org/x/oauth2"
)

// Function to calculate the SHA-256 hash of a file's content
func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Function to upload file to GitHub with hashed filename
func GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail string, fileContent []byte, githubOrg string, githubRepo string, pathFile string, replace bool) (content *github.RepositoryContentResponse, response *github.Response, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Upload file"),
		Content: fileContent,
		Branch:  github.String("main"),
		Author: &github.CommitAuthor{
			Name:  github.String(GitHubAuthorName),
			Email: github.String(GitHubAuthorEmail),
		},
	}

	content, response, err = client.Repositories.CreateFile(ctx, githubOrg, githubRepo, pathFile, opts)
	if (err != nil) && (replace) {
		currentContent, _, _, _ := client.Repositories.GetContents(ctx, githubOrg, githubRepo, pathFile, nil)
		opts.SHA = github.String(currentContent.GetSHA())
		content, response, err = client.Repositories.UpdateFile(ctx, githubOrg, githubRepo, pathFile, opts)
		return
	}

	return
}

// Function to get file content from GitHub repository
// Set header untuk mendownload file
// w.Header().Set("Content-Disposition", "attachment; filename=\"file.ext\"")
// w.Header().Set("Content-Type", "application/octet-stream")
// w.Header().Set("Content-Length", fmt.Sprint(len(fileContent)))
// // Tulis konten file ke response writer
// w.Write(fileContent)
func GithubGetFile(GitHubAccessToken, githubOrg, githubRepo, pathFile string) (fileContent []byte, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get file content from the repository
	downloadResponse, _, err := client.Repositories.DownloadContents(ctx, githubOrg, githubRepo, pathFile, nil)
	if err != nil {
		err = errors.New("error GetContents " + err.Error())
		return
	}
	defer downloadResponse.Close()

	// Read the binary content
	fileContent, err = io.ReadAll(downloadResponse)
	if err != nil {
		return nil, fmt.Errorf("error reading binary content: %w", err)
	}

	return
}
