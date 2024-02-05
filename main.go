package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			username := strings.TrimSpace(r.Form.Get("username"))
			if username != "" {
				userID, err := getGitHubUserID(username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(w, "GitHub ID for %s is %d", username, userID)
				return
			}
		}
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func getGitHubUserID(username string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned non-200 status code: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return 0, err
	}

	return user.ID, nil
}
