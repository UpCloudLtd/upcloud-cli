# upctl maintenance agent

Status: **design draft, not approved for implementation**. Open decisions are marked **[PROPOSED]** or **[OPEN]**. Decided items are marked **[DECIDED]**.

## What it is, in plain words

A scheduled and event-driven GitHub Actions workflow run by a coding agent (GitHub Agentic Workflows). It does four things:

1. Finds work: bugs, feature requests, new API operations, missing fields, and SDK v8 to v9 migration steps.
2. Writes a finding: one GitHub issue per finding, with what it found, the goal, and an overall assessment.
3. Waits for a maintainer. The maintainer talks to the agent in the issue comments. Nothing is implemented until the maintainer approves.
4. After approval, makes one small change, opens one pull request, and comments the PR link on the issue.

The agent never merges, never pushes to `main`, and never approves its own work. A human reviews every PR.

## Goals

- Move `upctl` from SDK v8 to v9, one command operation at a time.
- Work through incoming issues.
- Add commands for new API resources.
- Find parity gaps between the CLI and the OpenAPI spec, at the field and operation level.
- Find bugs and propose fixes.
- Keep every step visible in the GitHub UI. No local tooling is needed to review or approve.

## Non-goals

- No automatic merging, releasing, or tagging.
- No live API calls in the first version (see [Access](#access-to-the-api)).
- No change to the public CLI behaviour without a maintainer's approval.
- No generation of whole command groups in one PR.

## How a finding moves through the system

```
trigger -> discover -> finding issue -> maintainer conversation -> approval -> PR -> human review
```

### 1. Triggers

| Trigger | When | What it does |
|---|---|---|
| Issues | A new issue, or a label is added | Triage. Turns bug and feature reports into findings. |
| Schedule | Once a day | Discovery: bugs, parity gaps, migration steps. Creates at most 5 new findings per day. |
| SDK version check | Part of the daily run | Compares the `upcloud-go-api/v9` version in `go.mod` with the latest v9 tag. If it is newer, checks which new operations exist and creates findings for them. A release in the SDK repo cannot trigger this repo's workflow, so the check is a poll. |
| Manual | A maintainer starts it | Same as the schedule, on demand. |

**[DECIDED] OpenAPI changes are detected through the SDK, not the spec.** A change to the public spec does not start work by itself. The SDK must first ship the new operation. This avoids a finding that cannot be implemented yet. The daily run may compare the spec with the last SDK release and note the gap in the daily summary, but it does not create a finding for it.

### 2. Finding issue

Each finding is one GitHub issue with a fixed structure:

- **Title**: `[agent] <kind>: <short description>`, e.g. `[agent] parity: server create is missing field "x"`.
- **Labels**: `agent:finding`, a kind (`bug`, `feature`, `parity`, `sdk-v9`, `new-operation`), and `agent:awaiting-maintainer`.
- **Body**:
  - **What was found**: the evidence, with links to the code and the spec line.
  - **Goal**: what would change.
  - **Overall assessment**: impact, risk, effort (small, medium, large), and whether it is one PR.
  - **Proposed change**: a short plan, no code.
  - **How to approve**: the exact comment text below.

### 3. Conversation

- The maintainer replies in the issue. The agent reads replies and answers questions.
- The agent may update the finding body with the decisions made in the conversation.
- The agent does not implement anything during the conversation.
- **[PROPOSED]** A finding with no reply for 14 days is closed by the agent with a comment. It can be reopened by a maintainer.

### 4. Approval

- The approval is a comment that is exactly:

  ```
  @agent approved
  ```

  The comment must contain only that text, after trimming whitespace. Other comments, such as "looks good @agent approved", do not count.
- Only a member of the CODEOWNERS team may approve (see [Security](#security)).
- To disregard a finding, a maintainer closes the issue. The agent does not reopen closed findings.
- Editing an approval comment after it was read, or removing it, cancels the approval. The agent re-checks the comment before it starts work.

### 5. Implementation and PR

- The agent works on one operation per PR. Examples: `server create`, or a single bug fix. It does not change a whole command group in one PR.
- The agent pushes to a branch named `agent/<issue-number>-<slug>`. It never pushes to `main`.
- It opens a PR that links the issue and states what changed and how it was checked.
- It comments the PR link on the issue.
- The PR goes through the normal review and CI. The branch protection rules on `main` are unchanged.

## Priorities (what the agent works on first)

**[DECIDED]** The agent picks findings in this order. Each finding gets a priority label (`agent:p0` to `agent:p5`). The main goal is that `upctl` keeps parity with the SDK v9 and the OpenAPI spec first, then serves customer requests, then runs the long-term SDK v9 migration and general bug fixing.

| Priority | Kind | Why it comes here |
|---|---|---|
| P0 | Security issues and crashes that affect users, confirmed by a reproduction | Highest risk. Handled first, even out of order. |
| P1 | New API operations that the SDK v9 already exposes but the CLI does not | Parity with the SDK is the main goal. |
| P2 | Visible parity gaps between the CLI and the OpenAPI spec or SDK v9 (for example a missing field the user can see, or a missing operation on an existing resource) | Silent wrong output is harder to notice than a crash. |
| P3 | Feature requests from customers (issues) that need no API change | Customer demand, after parity. Needs product judgement, so it always waits for a maintainer. |
| P4 | SDK v8 to v9 migration, one operation per PR | Long-standing. Mechanical, so it fills the gaps. |
| P5 | Confirmed general bugs with no security or crash impact | Fixed when the queue has room. |

Ordering inside a priority: most-used resources first (server, storage, network), then the rest. Sort by the issue reaction count and the age of the report.


## Access to the API

**[DECIDED]** Spec-only mode for the first version. The agent reads the OpenAPI spec and the SDK source. It does not call the live API.

- Behaviour checks would use a mock server in the workflow. The mock is not in scope for the first version.
- A live read-only token would come later, as a separate decision, with its own review.

## Secrets

**What this means:** GitHub stores secrets (tokens, keys) for a repo. Workflows read them at run time. A secret is exposed to the workflow, so a leaked or misused one can do damage.

**Decision:** the spec-only plan needs no API secrets. The agent's access to the model is managed by GitHub Agentic Workflows.

**What the agent needs to write:** an issue, a comment, and a branch. This is done with a fine-grained personal access token scoped to this repo only, with `contents: write` (for `agent/*` branches only, enforced by branch protection on `main` and the workflow rules), `issues: write`, and `pull-requests: write`. It has no admin, no `workflow`, and no `actions` scope. The token is the only secret. It is stored in a protected GitHub environment with required reviewers. Fine-grained tokens expire, so the owner must rotate them.

## Security

Security is a design constraint, not a feature. The threat model:

| Threat | Mitigation |
|---|---|
| Prompt injection from issue text, comments, or spec text | Untrusted text is filtered by the agent's integrity layer before the agent reads it. The agent's instructions say issue text is data, never instructions. |
| A user tricks the agent into running code | The agent does not run code from issues or from PR branches. Only the repo's own build and test run. |
| A non-maintainer approves work | Approval checks the commenter against CODEOWNERS (see below). |
| A stale or edited approval | Approval is checked again right before work starts. |
| The agent pushes to `main` or merges | Branch protection on `main` requires review. The agent's token cannot push to `main`. |
| The agent changes its own permissions, workflows, or CI | The agent may not edit `.github/`, `agent/`, or release files. A PR that touches them is blocked by a CODEOWNERS rule. |
| Secret leakage in logs or PRs | The agent cannot read secrets. Safe outputs are checked for secrets before they are written. |
| Runaway cost or spam | At most 5 new findings per day. One agent run per finding at a time. |
| Loops (the agent reacting to its own comments) | Events from the agent's own account are ignored. |
| Fork PRs or outside contributors | The workflow does not run for fork PRs with secrets. |

### Approver rule

- **[DECIDED]** Only members of the CODEOWNERS team may approve. The workflow checks team membership with the GitHub API before it starts.
- The check runs in the trigger job and again before work starts.

## Output format for the report

Every finding issue and the daily summary use the same shape, so a maintainer can scan them quickly:

```
Found: <one line>
Goal: <one line>
Assessment: impact <low|medium|high>, risk <low|medium|high>, effort <small|medium|large>
Evidence: <links>
Proposed change: <short plan>
Approve with a comment that is exactly: @agent approved
```

## Daily cap

**[DECIDED]** At most 5 new findings per day. The cap counts only new issues. Comments and PRs for existing findings do not count.

## Verification

**[DECIDED]** Verification is unit tests and linters only, using the repo's existing commands. Ideas for stronger checks are in [Future improvements](#future-improvements).

## Future improvements

Not in the first version. Revisit if the unit-test-only checks let regressions through.

- Mock-server runs of the changed command, comparing the request and the output with the expected shape.
- A live read-only token for parity checks, as a separate security review.
- Agent review of other people's PRs.
- A metric on approved versus disregarded findings, to tune the daily cap and priorities.

## Open questions

- **[OPEN]** Who owns the fine-grained token, and who rotates it when it expires?
- **[OPEN]** Whether the agent should also review other people's PRs (not in scope now).


## Files in this folder

- `README.md`: this design document.
- `workflow.md`: the agent definition (draft, not compiled yet).
