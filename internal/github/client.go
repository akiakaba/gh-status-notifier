package github

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var fields = []string{
	"assignees",
	"author",
	"baseRefName",
	"closed",
	"closedAt",
	"headRefName",
	"mergeCommit",
	"mergeStateStatus",
	"mergeable",
	"mergedAt",
	"number",
	"potentialMergeCommit",
	"reviewDecision",
	"reviewRequests",
	"reviews",
	"state",
	"statusCheckRollup",
	"title",
	"updatedAt",
}

func FetchPRStatus() (AllOfPRStatus, error) {
	cmd := exec.Command("gh", "pr", "status", "--json", strings.Join(fields, ","))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var resp AllOfPRStatus
	if err := cmd.Run(); err != nil {
		return resp, errors.Wrap(err, stderr.String())
	}
	data := stdout.Bytes()

	// FOR DEV
	//data, err := os.ReadFile("./.mynotes/fixtures/sample_response.json")
	//if err != nil {
	//	log.Fatal(err)
	//}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
