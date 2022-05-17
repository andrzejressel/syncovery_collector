package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SyncoveryClient struct {
	url *url.URL
}

func CreateSyncoveryClient(urlString string) (*SyncoveryClient, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	client := SyncoveryClient{
		url: u,
	}
	return &client, err
}

type profilesDto struct {
	Rows []struct {
		Number          string    `json:"number"`
		Name            string    `json:"name"`
		RightPath       string    `json:"rightPath"`
		LeftPath        string    `json:"leftPath"`
		MixedProgress   string    `json:"mixedProgress"`
		IsDisabled      string    `json:"isDisabled"`
		IsLocked        string    `json:"isLocked"`
		LastRunString   string    `json:"lastRunString"`
		NextRunString   string    `json:"nextRunString"`
		LastResult      string    `json:"lastResult"`
		CurrentProgress string    `json:"currentProgress"`
		LastRunDateTime time.Time `json:"lastRunDateTime"`
		NextRunDateTime string    `json:"nextRunDateTime"`
		LastRunHadError string    `json:"lastRunHadError"`
		IsRunning       string    `json:"isRunning"`
		IsScheduled     string    `json:"isScheduled"`
	} `json:"Rows"`
}

type Profile struct {
	Number          int
	Name            string
	LastRunDateTime time.Time
	LastRunHadError bool
	IsRunning       bool
	IsScheduled     bool
	IsDisabled      bool
	IsLocked        bool
	NextRunDateTime *time.Time
}

func (client SyncoveryClient) GetProfiles() (*[]Profile, error) {
	path, err := client.url.Parse("profilelist.json")
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(path.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response []profilesDto
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	rows := response[0].Rows

	var profiles []Profile
	for _, row := range rows {
		number, err := strconv.Atoi(row.Number)
		if err != nil {
			return nil, err
		}

		var nextRunDateTime *time.Time = nil
		if row.NextRunDateTime != "" {
			time, err := time.Parse(time.RFC3339, row.NextRunDateTime)
			if err != nil {
				return nil, err
			}
			nextRunDateTime = &time
		}

		profiles = append(profiles, Profile{
			Number:          number,
			Name:            row.Name,
			LastRunDateTime: row.LastRunDateTime,
			LastRunHadError: mustStringToBoolean(row.LastRunHadError),
			IsRunning:       mustStringToBoolean(row.IsRunning),
			IsScheduled:     mustStringToBoolean(row.IsScheduled),
			IsDisabled:      mustStringToBoolean(row.IsDisabled),
			IsLocked:        mustStringToBoolean(row.IsLocked),
			NextRunDateTime: nextRunDateTime,
		})
	}

	return &profiles, nil

}

func mustStringToBoolean(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		panic(err)
	}
	return value
}
