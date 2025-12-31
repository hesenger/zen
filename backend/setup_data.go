package main

type App struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
	Command  string `json:"command"`
}

type SetupData struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	GithubToken string `json:"githubToken"`
	Apps        []App  `json:"apps"`
}
