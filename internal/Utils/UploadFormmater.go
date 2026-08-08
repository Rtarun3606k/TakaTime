package utils

import (
	"fmt"
	"strings"

	"github.com/Rtarun3606k/TakaTime/internal/types"
)

func FormmatUpload(token string, repoURL string, path string, branch string, commitMsg string) (types.UploadStruct, error) {

	if len(commitMsg) == 0 {
		commitMsg = "Adding toadys stats"
	}
	splitRepoName := strings.Split(repoURL, "/")

	if len(splitRepoName) < 2 {
		return types.UploadStruct{}, fmt.Errorf("invalid repo URL format, expected 'owner/repo', got: %s", repoURL)
	}

	uploadStruct := types.UploadStruct{
		Token:     token,
		Owner:     splitRepoName[0],
		Repo:      splitRepoName[1],
		Path:      path,
		Branch:    branch,
		CommitMsg: commitMsg,
	}

	return uploadStruct, nil

}
