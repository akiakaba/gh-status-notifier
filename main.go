package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sisisin/gh-status-notifier/internal/github"
	"github.com/sisisin/gh-status-notifier/internal/macos"
)

const interval = 1 * time.Minute

// FIXME: キーが PR number なので複数のリポジトリに対応しない
var beforeResults = make(map[int]cachedResult, 0)

type cachedResult struct {
	Title  string                `json:"title"`
	Number int                   `json:"number"`
	Closed bool                  `json:"closed"`
	Checks map[string]statusPart `json:"checks"`
}

type statusPart struct {
	Conclusion string `json:"conclusion"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	State      string `json:"state"`
}

func main() {
	fmt.Println("Started.")
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			allOfPRStatus, err := github.FetchPRStatus()
			if err != nil {
				log.Fatal(err)
			}
			if err := checkPRStatus(allOfPRStatus); err != nil {
				log.Fatal(err)
			}
			time.Sleep(interval)
		}
	}()
	sig := <-trap
	fmt.Printf("received %s\n", sig.String())
	if j, err := json.Marshal(beforeResults); err == nil {
		fmt.Printf("last check result: %s\n", j)
	} else {
		log.Fatal(err)
	}
}

func checkPRStatus(result github.AllOfPRStatus) error {
	fmt.Println("checking PR status")

	for _, pr := range result.CreatedBy {
		if _, exists := beforeResults[pr.Number]; !exists {
			empty := make(map[string]statusPart, 0)
			beforeResults[pr.Number] = cachedResult{
				Title:  pr.Title,
				Number: pr.Number,
				Closed: pr.Closed,
				Checks: empty,
			}
		}
		if pr.Closed {
			delete(beforeResults, pr.Number)
			continue
		}

		for _, action := range pr.StatusCheckRollup {
			actionKey := genKey(action)
			failed := isFailure(action)
			if actionKey == "" {
				fmt.Printf("cannot get action name/context: %v\n", action)
				if failed {
					msg := fmt.Sprintf("PR %s #%v failed", pr.Title, pr.Number)
					notify(msg)
				}
				continue
			}

			before, _ := beforeResults[pr.Number].Checks[actionKey]
			current := statusPart{
				Name:       action.Name,
				Conclusion: action.Conclusion,
				State:      action.State,
				Status:     action.Status,
			}
			if !isSameState(before, current) && failed {
				msg := fmt.Sprintf("PR %s #%v %s failed", pr.Title, pr.Number, actionKey)
				fmt.Println(msg)
				notify(msg)
			}
			beforeResults[pr.Number].Checks[actionKey] = current
		}
	}

	return nil
}

func genKey(action github.StatusCheckRollup) string {
	if action.Name == "" {
		return action.Context
	}
	return action.Name
}

func isFailure(action github.StatusCheckRollup) bool {
	return action.State == "FAILURE" || action.Conclusion == "FAILURE"
}

func isSameState(s, t statusPart) bool {
	return s.Conclusion == t.Conclusion &&
		s.State == t.State &&
		s.Status == t.Status
}

func notify(msg string) {
	if err := macos.Notify(msg); err != nil {
		fmt.Printf("Notification failed: %s\n", msg)
	}
}
