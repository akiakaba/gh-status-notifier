import cp from "node:child_process";
import { AllOfPRStatus, StatusCheckRollup } from "./types";

function notify(text: string) {
  cp.execSync(
    `osascript -e 'display notification ${JSON.stringify(
      text
    )} with title "gh status" sound name "Boop"'`
  );
}

type StatusPart = Pick<
  StatusCheckRollup,
  "conclusion" | "state" | "status" | "workflowName" | "name"
>;
type CachedResult = {
  [key: number]: {
    title: string;
    number: number;
    closed: boolean;
    checks: {
      [key: string]: StatusPart & {
        lastCheckedAt: number;
      };
    };
  };
};
const beforeResults: CachedResult = {};

function fetchPrStatus() {
  const resultText = cp
    .execSync(
      "gh pr status --json 'assignees,author,baseRefName,closed,closedAt,headRefName,mergeCommit,mergeStateStatus,mergeable,mergedAt,number,potentialMergeCommit,reviewDecision,reviewRequests,reviews,state,statusCheckRollup,title,updatedAt'",
      { cwd: process.cwd() }
    )
    .toString();
  const result: AllOfPRStatus = JSON.parse(resultText);
  return result;
}

// const r = require("../src/fixtures/sample.json");
// checkPrStatus(r);
// checkPrStatus(r);
checkPrStatus(fetchPrStatus());
const interval = setInterval(() => {
  checkPrStatus(fetchPrStatus());
}, 1000 * 60 * 1);

process.on("SIGINT", () => {
  console.log("received SIGINT");
  console.log("last check result", JSON.stringify(beforeResults));
  clearInterval(interval);
});

function checkPrStatus(result: AllOfPRStatus) {
  console.log("checking PR status");

  for (const prStatus of result.createdBy) {
    if (beforeResults[prStatus.number] == null) {
      beforeResults[prStatus.number] = {
        title: prStatus.title,
        number: prStatus.number,
        closed: prStatus.closed,
        checks: {},
      };
    }

    if (prStatus.closed) {
      delete beforeResults[prStatus.number];
      continue;
    }

    for (const action of prStatus.statusCheckRollup) {
      const name = action.name ?? action.context;

      if (name == null) {
        console.log("cannot get action name", action);
        if (isFailure(action)) {
          notify(`PR ${prStatus.title} #${prStatus.number} failed`);
        }
        continue;
      }

      const lastCheckedAt = Date.now();
      if (beforeResults[prStatus.number].checks[name] == null) {
        beforeResults[prStatus.number].checks[name] = {
          conclusion: action.conclusion,
          state: action.state,
          status: action.status,
          workflowName: action.workflowName,
          name: action.name,
          lastCheckedAt,
        };
        if (isFailure(action)) {
          notify(`PR ${prStatus.title} #${prStatus.number} ${name} failed`);
        }
      } else {
        const before = beforeResults[prStatus.number].checks[name];
        beforeResults[prStatus.number].checks[name] = {
          ...before,
          conclusion: action.conclusion,
          state: action.state,
          status: action.status,
          lastCheckedAt,
        };
        if (!isSameState(before, action) && isFailure(action)) {
          console.log("state changed to failure", before, action);
          notify(`PR ${prStatus.title} #${prStatus.number} ${name} failed`);
        }
      }
    }
  }
}
function isFailure(action: StatusCheckRollup) {
  return action.state === "FAILURE" || action.conclusion === "FAILURE";
}
function isSameState(a: StatusPart, b: StatusPart) {
  return (
    a.conclusion === b.conclusion &&
    a.state === b.state &&
    a.status === b.status
  );
}
