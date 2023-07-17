package github

type AllOfPRStatus struct {
	CreatedBy     []PRStatus `json:"createdBy"`
	CurrentBranch PRStatus   `json:"currentBranch"`
	NeedsReview   []any      `json:"needsReview"`
}

type PRStatus struct {
	Assignees            []any                `json:"assignees"`
	Author               Author               `json:"author"`
	BaseRefName          string               `json:"baseRefName"`
	Closed               bool                 `json:"closed"`
	ClosedAt             any                  `json:"closedAt"`
	HeadRefName          string               `json:"headRefName"`
	MergeCommit          any                  `json:"mergeCommit"`
	MergeStateStatus     string               `json:"mergeStateStatus"`
	Mergeable            string               `json:"mergeable"`
	MergedAt             any                  `json:"mergedAt"`
	Number               int                  `json:"number"`
	PotentialMergeCommit PotentialMergeCommit `json:"potentialMergeCommit"`
	ReviewDecision       string               `json:"reviewDecision"`
	ReviewRequests       []ReviewRequest      `json:"reviewRequests"`
	Reviews              []any                `json:"reviews"`
	State                string               `json:"state"`
	StatusCheckRollup    []StatusCheckRollup  `json:"statusCheckRollup"`
	Title                string               `json:"title"`
	UpdatedAt            string               `json:"updatedAt"`
}

type Author struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type PotentialMergeCommit struct {
	OID string `json:"oid"`
}

type ReviewRequest struct {
	Login string `json:"login"`
}

type StatusCheckRollup struct {
	Conclusion   string `json:"conclusion"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	WorkflowName string `json:"workflowName"`
	State        string `json:"state"`
	CompletedAt  string `json:"completedAt"`
	DetailsURL   string `json:"detailsUrl"`
	StartedAt    string `json:"startedAt"`
	Context      string `json:"context"`
	TargetURL    string `json:"targetUrl"`
}
