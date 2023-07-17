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
	Conclusion   string `json:"conclusion"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	WorkflowName string `json:"workflow_name"`
	State        string `json:"state"`
}

func main() {
	fmt.Println("Started.")
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			allOfPRStatus, err := github.FetchPrStatus()
			if err != nil {
				log.Fatal(err)
			}
			if err := checkPrStatus(allOfPRStatus); err != nil {
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

func checkPrStatus(result github.AllOfPRStatus) error {
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
					if err := macos.Notify(msg); err != nil {
						fmt.Printf("Notification failed: %s", msg)
					}
				}
				continue
			}

			before, exists := beforeResults[pr.Number].Checks[actionKey]
			var current statusPart
			if !exists {
				current = statusPart{
					Name:         action.Name,
					WorkflowName: action.WorkflowName,
					Conclusion:   action.Conclusion,
					State:        action.State,
					Status:       action.Status,
				}
				if failed {
					msg := fmt.Sprintf("PR %s #%v %s failed", pr.Title, pr.Number, actionKey)
					fmt.Println(msg)
					if err := macos.Notify(msg); err != nil {
						fmt.Printf("Notification failed: %s", msg)
					}
				}
			} else {
				current = statusPart{
					Name:         before.Name,
					WorkflowName: before.WorkflowName,
					Conclusion:   action.Conclusion,
					State:        action.State,
					Status:       action.Status,
				}
				if !isSameState(before, current) && failed {
					fmt.Printf("state changed to failure: %v -> %v\n", before, action)
					msg := fmt.Sprintf("PR %s #%v %s failed", pr.Title, pr.Number, actionKey)
					if err := macos.Notify(msg); err != nil {
						fmt.Printf("Notification failed: %s", msg)
					}
				}
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
