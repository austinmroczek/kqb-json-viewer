package bgl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const API_BASE_URL = "https://api-staging.beegame.gg/"
const PROD_BASE_URL = "https://api.beegame.gg/v1/"

type BGLData struct {
	Token       string
	Matches     map[string]int
	Teams       map[string]int
	Players     map[string]int
	HomeID      int
	AwayID      int
	HomeName    string
	AwayName    string
	matchResult MatchResult
	User        User
}

func (b *BGLData) LoadCurrentMatchesLocal() error {
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "/fixtures/match_result.json")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("Could not read json file", err)
		return err
	}
	err = json.Unmarshal(data, &b.matchResult)
	if err != nil {
		log.Println("Could not parse json file", err)
		return err
	}

	b.Matches = make(map[string]int)
	for _, result := range b.matchResult.Results {
		key := result.Away.Name + " @ " + result.Home.Name
		b.Matches[key] = result.ID
	}
	return nil
}

func (b *BGLData) GetMatchNames() []string {
	matches := make([]string, 0)
	for key, _ := range b.Matches {
		matches = append(matches, key)
	}
	return matches
}

func (b *BGLData) LoadPlayersForMatch(match string) {
	matchID := b.Matches[match]
	players := make(map[string]int)
	for _, result := range b.matchResult.Results {
		if result.ID == matchID {
			for _, member := range result.Home.Members {
				players[member.Name] = int(member.ID)
			}
			for _, member := range result.Away.Members {
				players[member.Name] = int(member.ID)
			}
		}
	}
	b.Players = players
}

func (b *BGLData) GetPlayerNames() []string {
	players := make([]string, 0)
	for playerName, _ := range b.Players {
		players = append(players, playerName)
	}
	return players
}

func (b *BGLData) LoadTeamsForMatch(match string) {
	matchID := b.Matches[match]
	teams := make(map[string]int)
	for _, result := range b.matchResult.Results {
		if result.ID == matchID {
			teams[result.Away.Name] = result.Away.ID
			teams[result.Home.Name] = result.Home.ID
			b.HomeID = result.Home.ID
			b.HomeName = result.Home.Name
			b.AwayID = result.Away.ID
			b.AwayName = result.Away.Name
		}
	}
	b.Teams = teams
}

func (b *BGLData) GetTeamNames() []string {
	teams := make([]string, 0)
	for name, _ := range b.Teams {
		teams = append(teams, name)
	}
	return teams
}

func (b *BGLData) GetMe() error {
	url := API_BASE_URL + "me/?format=json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Authorization", "Token "+b.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &b.User)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println("Get Me Status Code Error")
		return errors.New("invalid status code")
	}
	defer res.Body.Close()

	return nil
}

func (b *BGLData) HandleMatchUpdate(result Result) error {
	wd, _ := os.Getwd()
	outPath := filepath.Join(wd, "/tmp/match_update.json")
	output, err := json.MarshalIndent(result, "  ", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b *BGLData) SaveRawOutput(final FinalOutput) error {
	wd, _ := os.Getwd()
	outPath := filepath.Join(wd, "/tmp/match_output.json")
	output, err := json.MarshalIndent(final, "  ", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b *BGLData) LoadCurrentMatches() error {
	url := API_BASE_URL + "matches/?format=json&limit=5&awaiting_results=true"
	method := "GET"

	var teamIDs []int
	for _, team := range b.User.Player.Teams {
		teamIDs = append(teamIDs, team.ID)
	}
	// log.Println(teamIDs)
	url += fmt.Sprintf("&team_id=%d", teamIDs[0])

	// for _, id := range teamIDs {
	// 	url += fmt.Sprintf("&team_id=%d", id)
	// }

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Authorization", "Token "+b.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("invalid status code")
	}

	// fmt.Println(string(body))
	err = json.Unmarshal(body, &b.matchResult)
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(b.matchResult)
	// Initialize empty Match Object
	if len(b.Matches) == 0 {
		b.Matches = make(map[string]int)
	}
	for _, result := range b.matchResult.Results {
		key := result.Away.Name + " @ " + result.Home.Name
		b.Matches[key] = result.ID
	}
	return nil
}
