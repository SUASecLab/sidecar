package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func decide(service string, claims RulesCustomClaims) bool {
	rulesFile, err := os.Open("rules/rules.json")
	if err != nil {
		log.Println("Can not open rules file: ", err)
		return false
	}
	defer rulesFile.Close()

	content, err := io.ReadAll(rulesFile)
	if err != nil {
		log.Println("Can not read rules file: ", err)
		return false
	}

	var rules RulesFile
	json.Unmarshal(content, &rules)

	for _, rule := range rules.Rules {
		if rule.Key == service {
			if rule.Value == "allowed" {
				return true
			} else if rule.Value == "moderator" {
				return claims.Moderator
			} else {
				return false
			}
		}
	}
	return false
}
