package main

type AuthenticationResponse struct {
	Allowed bool `json:"allowed"`
}

type IssuanceResult struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

type RulesFile struct {
	Rules []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"rules"`
}

type RulesCustomClaims struct {
	Moderator bool     `json:"moderator"`
	Tags      []string `json:"tags"`
	UUID      string   `json:"uuid"`
}

type UserInfo struct {
	Exists  bool   `json:"exists"`
	IsAdmin bool   `json:"isAdmin"`
	UUID    string `json:"uuid"`
}

type User struct {
	Uuid         string   `json:"uuid"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	VisitCardUrl string   `json:"visitCardUrl"`
	Tags         []string `json:"tags"`
}
type ValidationResult struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}
