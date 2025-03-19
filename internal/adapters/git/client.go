// Package git defines an adapter for git commands natively within Butler.
//
// # Copyright (c) 2025, The Butler Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"go.uber.org/zap"
)

// GitClient implements the GitAdapter interface
type GitClient struct {
	RepoURL  string
	Username string
	Token    string
	Logger   *zap.Logger
}

// NewGitClient initializes a Git client
func NewGitClient(repoURL, username, token string, logger *zap.Logger) *GitClient {
	return &GitClient{
		RepoURL:  repoURL,
		Username: username,
		Token:    token,
		Logger:   logger,
	}
}

// CloneRepo clones a Git repository to the specified local path
func (g *GitClient) CloneRepo(ctx context.Context, localPath string) error {
	g.Logger.Info("Cloning Git repository", zap.String("repo", g.RepoURL), zap.String("path", localPath))

	// Check if the repo already exists
	if _, err := os.Stat(filepath.Join(localPath, ".git")); err == nil {
		g.Logger.Info("Repository already cloned, pulling latest changes...")
		repo, err := git.PlainOpen(localPath)
		if err != nil {
			return fmt.Errorf("failed to open existing repo: %w", err)
		}
		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: "oauth2",
				Password: g.Token,
			},
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull latest changes: %w", err)
		}
		g.Logger.Info("Git repository is up-to-date")
		return nil
	}

	// Clone the repo if it doesn't exist
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL: g.RepoURL,
		Auth: &http.BasicAuth{
			Username: "oauth2",
			Password: g.Token,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	g.Logger.Info("Git repository cloned successfully", zap.String("path", localPath))
	return nil
}

// CommitAndPush stages, commits, and pushes changes
func (g *GitClient) CommitAndPush(ctx context.Context, localPath, commitMessage string) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Stage changes
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Check if there are actually changes to commit
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %w", err)
	}

	if status.IsClean() {
		g.Logger.Info("No changes detected, skipping commit.")
		return nil
	}

	// Commit changes
	_, err = w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Butler Automation",
			Email: "butler@butler.dev",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Push changes to remote
	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: g.Username,
			Password: g.Token,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	g.Logger.Info("Git changes pushed successfully")
	return nil
}