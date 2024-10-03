package models

type GithubUser struct {
	Login string `json:"login"`
}

type PatreonUser struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}
