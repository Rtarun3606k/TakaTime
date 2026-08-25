package utils

import (
	"reflect"
	"testing"

	"github.com/Rtarun3606k/TakaTime/internal/types"
)

func TestFormmatUpload(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		repoURL   string
		path      string
		branch    string
		commitMsg string
		want      types.UploadStruct
		wantErr   bool
	}{
		{
			name:      "Happy Path - All Valid Inputs",
			token:     "ghp_12345",
			repoURL:   "Rtarun3606k/TakaTime",
			path:      "stats/weekly.png",
			branch:    "main",
			commitMsg: "Update weekly stats",
			want: types.UploadStruct{
				Token:     "ghp_12345",
				Owner:     "Rtarun3606k",
				Repo:      "TakaTime",
				Path:      "stats/weekly.png",
				Branch:    "main",
				CommitMsg: "Update weekly stats",
			},
			wantErr: false,
		},
		{
			name:      "Empty Commit Message Fallback",
			token:     "ghp_abcde",
			repoURL:   "user/repo",
			path:      "image.png",
			branch:    "dev",
			commitMsg: "", // Triggers the len(commitMsg) == 0 branch
			want: types.UploadStruct{
				Token:     "ghp_abcde",
				Owner:     "user",
				Repo:      "repo",
				Path:      "image.png",
				Branch:    "dev",
				CommitMsg: "Adding toadys stats",
			},
			wantErr: false,
		},
		{
			name:      "Invalid Repo URL Format (No Slash)",
			token:     "ghp_invalid",
			repoURL:   "just-the-repo-name", // Triggers the splitRepoName < 2 branch
			path:      "file.txt",
			branch:    "main",
			commitMsg: "init",
			want:      types.UploadStruct{},
			wantErr:   true,
		},
		{
			name:      "Multiple Slashes (Takes first two parts)",
			token:     "token",
			repoURL:   "owner/repo/extra",
			path:      "file",
			branch:    "main",
			commitMsg: "msg",
			want: types.UploadStruct{
				Token:     "token",
				Owner:     "owner",
				Repo:      "repo",
				Path:      "file",
				Branch:    "main",
				CommitMsg: "msg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormmatUpload(tt.token, tt.repoURL, tt.path, tt.branch, tt.commitMsg)

			if (err != nil) != tt.wantErr {
				t.Errorf("FormmatUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormmatUpload() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

