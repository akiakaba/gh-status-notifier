export interface AllOfPRStatus {
  createdBy: PRStatus[];
  currentBranch: PRStatus;
  needsReview: any[];
}

export interface PRStatus {
  assignees: any[];
  author: Author;
  baseRefName: string;
  closed: boolean;
  closedAt: any;
  headRefName: string;
  mergeCommit: any;
  mergeStateStatus: string;
  mergeable: string;
  mergedAt: any;
  number: number;
  potentialMergeCommit?: PotentialMergeCommit;
  reviewDecision: string;
  reviewRequests: ReviewRequest[];
  reviews: any[];
  state: string;
  statusCheckRollup: StatusCheckRollup[];
  title: string;
  updatedAt: string;
}

export interface Author {
  id: string;
  is_bot: boolean;
  login: string;
  name: string;
}

export interface PotentialMergeCommit {
  oid: string;
}

export interface ReviewRequest {
  __typename: string;
  login: string;
}

export interface StatusCheckRollup {
  __typename: string;
  completedAt?: string;
  conclusion?: string;
  detailsUrl?: string;
  name?: string;
  startedAt: string;
  status?: string;
  workflowName?: string;
  context?: string;
  state?: string;
  targetUrl?: string;
}
