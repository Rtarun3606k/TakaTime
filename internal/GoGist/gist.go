package gogist

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

func UpdateGist(token string, gistID string, content string) error {

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 2. Define the Update
	// We only need to provide the fields we want to change.
	// Notice we use github.String() to safely create pointers.
	gist := &github.Gist{
		Description: github.String("📊 TakaTime Stats (Auto-Updated)"),
		Files: map[github.GistFilename]github.GistFile{
			"taka-stats.txt": {Content: github.String(content)},
		},
	}

	_, _, err := client.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		log.Println("Some error occured during editing gist", err)
		return err
	}

	return nil
}

func UpdateReadMe(githubToken string, repoName string, content string) error {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	path := "README.md"
	startMarker := "<!--takatime-start-->"
	endMarker := "<!--takatime-end-->"
	splitRepoName := strings.Split(repoName, "/")

	getReadMeFileContent, _, _, err := ghClient.Repositories.GetContents(ctx, splitRepoName[0], splitRepoName[1], path, nil)

	if err != nil {
		log.Println("Could not get filecontents ", err)
		return err
	}

	currentText, err := getReadMeFileContent.GetContent()
	if err != nil {
		log.Println("Could not get filecontents ", err)
		return err
	}

	// newBlock := fmt.Sprintf("%s\n```text\n%s\n```\n%s", startMarker, content, endMarker)
	newBlock := fmt.Sprintf("%s\n\n%s\n\n%s", startMarker, content, endMarker)

	// Regex to replace existing block (dot matches newline with (?s))
	re := regexp.MustCompile(fmt.Sprintf(`(?s)%s.*?%s`, regexp.QuoteMeta(startMarker), regexp.QuoteMeta(endMarker)))

	var newReadme string
	if re.MatchString(currentText) {
		// Replace existing block
		newReadme = re.ReplaceAllString(currentText, newBlock)
	} else {
		// Append to end if markers missing (Safeguard)
		fmt.Println("⚠️ Markers not found, appending to end of README.")
		newReadme = currentText + "\n" + newBlock
	}

	if newReadme == currentText {
		log.Println("No changes dectected")
		return nil
	}

	msg := "chore : update ReadMe takatime"

	opts := &github.RepositoryContentFileOptions{
		Message: &msg,
		Content: []byte(newReadme),
		SHA:     getReadMeFileContent.SHA,
	}

	_, _, err = ghClient.Repositories.UpdateFile(ctx, splitRepoName[0], splitRepoName[1], path, opts)
	if err != nil {
		log.Println("Error occured during updaing file contents ", err)
		return err
	}

	return nil
}
