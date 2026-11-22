---
name: writing-unit-tests
description: Use when adding or changing a Go unit test (including the failing test of a TDD red step) for a use case, worker or other package under internal/ (gomock, mockgen, table-driven tests, *_test.go without a build tag), when a use case has no or thin unit tests, when mocks are missing or hand-rolled, or when deciding what a unit test should assert.
---

# Writing unit tests

## Overview

A unit test runs one use case against **generated mocks of its ports** and checks what the **spec** promises: the result, the error, what was written, and what must not happen. It never touches Postgres, Redis or the network. Coverage is a by-product: a line that is executed but whose behaviour no assertion checks is not tested.

Other layers: persisted state, transactions and SQL go to `writing-integration-tests`. Flows across workers, chain and webhooks go to `writing-e2e-scenarios`.

## With TDD

`superpowers:test-driven-development` decides *when*: test first, watch it fail, then the minimal code. This skill decides *how* in this repo. In the red step, steps 2–6 below still apply (spec first, generated mocks, real Principal). The failing run against missing or buggy code **is** the prove-fail for that test, so you don't need a mutation for it. For a bug that has a `t.Skip("known bug …")` test, remove the skip: that is your red. Use mutations (Prove-fail) for tests you add around code that already works.

## Recipe

1. **Find what already exists.** `grep -rn "<UseCaseName>\|New<UseCaseName>" internal/usecases/<domain>/*_test.go`. Tests often live in a shared file (`rbac_test.go`, `management_helpers_test.go`) or in `<usecase>_unit_test.go`, not `<usecase>_test.go`. Extend the existing file, or add `<usecase>_test.go` for the gaps. If you replace an old hand-rolled (testify) test, delete it in the same change. Existing tests can be weak too: check them with prove-fail like your own.
2. **Expectations come from the spec, not the code.** Open the feature spec (`docs/02-requirements/specs/F-XXX-*.md`; `docs/02-requirements/features.md` maps features). Before you read the use case body, write down the behaviours to prove, with their BR ids. Then read the code only to learn how to call it. Some behaviours belong to another layer: permission gates in `internal/handler/intersectors/rpc_actions.go` when the use case takes no Principal, response fields set by the handler, refusal reasons decided by `internal/access`. Move them to "not tested here, and why". Don't drop them silently. The domain `CLAUDE.md` describes today's code. It is never the expectation.
   **Required rows.** Whatever the spec says, the behaviour list also has a row for each shape below that the use case has. Mark each row covered (test name) or N/A with the reason. Rules: `.claude/rules/rules.md`.

   | The use case… | Required behaviour | Rule |
   |---|---|---|
   | reads a row by id or key | **Miss vs failure:** `repository.ErrNotFound` gives the domain not-found; any other repo error is returned as an error, never as not-found, "absent", or a reason to create. One test each. | R-021 |
   | is run by a River worker (`workers/` calls its `Execute`) | **Re-run after commit:** run `Execute` a second time on the state the first run committed (the worker's next-job insert failed, River retries the step). It writes, credits, broadcasts and enqueues nothing new, and returns the same next step. | R-023 |
   | takes a Principal and a tenant-scoped id | **Other tenant:** a caller of another company gets the same error as for a missing row, and no data. | R-019, R-020 |
   | guards an invariant on a set (last admin, cap, "not yet credited") | **Concurrency:** belongs to `writing-integration-tests` § Patterns. List the row here and point to that test. | R-022 |

3. **When the code disagrees with the spec**, decide by the spec:
   - **The spec states it and the code violates it.** File a bug with `bd create --type=bug` (file:line, observed vs expected). Write the test as the spec says, starting with `t.Skip("known bug demo1: <one line>")`. Remove the skip once and see the test fail for that reason before you put it back; otherwise the skipped test may assert nothing. Never write the test to the code's behaviour. If an **existing** test pins the buggy behaviour, leave it unchanged and name it in the bug: the fix PR updates it.
   - **The spec is silent.** Either ask, or pin today's behaviour with a `// Q-XXX: spec silent; pins current behaviour` comment above the test. Add the question to `docs/05-progress/open-questions.md` with the next free Q number, in both its quick-scan table and its details section.
   - **The behaviour looks wrong but you can't prove it** (for example, the cache is cleared before a commit that may fail). That is a finding for your report, not an expectation. Don't assert it.
4. **Mocks are generated, never written by hand** (R-016). See the next section.
5. **Build the use case with its constructor**, mocks for every port. The caller is a real Principal: `principal.ForUser(id, &companyID, entity.RoleAdmin)`, `principal.ForAPIKey(...)`, `principal.System("job")`. Use the real `internal/authz` Policy; never mock it.
6. **Write the tests** (see Assertions). One behaviour per `t.Run`. Name each one after the outcome ("rolls back and publishes nothing when the commit fails").
7. **Prove every test can fail.** See Prove-fail.
8. **Run:** `go test -count=1 -run '^TestX' ./internal/usecases/<domain>/`, then the whole package, then `go vet ./internal/usecases/<domain>/` and `gofmt -l internal/usecases/<domain>/`.

## Mocks

- A use case's ports are either canonical repository interfaces in `internal/repository/<domain>` (R-018) or the domain's own `ports/`. The mocks for both come from `//go:generate` lines in `internal/usecases/<domain>/generate.go`. Look there first, then `grep -n "func NewMock" internal/usecases/<domain>/mocks*_test.go`.
- **A mock is missing?** Add one `//go:generate` line to `generate.go` in the form the file already uses, which is usually bare `mockgen -destination=mocks_<name>_test.go -package=<pkg> <import path> <Interface>`. Regenerate only your package with `go generate ./internal/usecases/<domain>/`. `make mocks` regenerates the whole repository.
- **`-package` must match the test files that use the mock** (`<pkg>` or `<pkg>_test`).
- **New tests never use testify `mock.Mock`**, even when the neighbouring tests do.

## Time and logging

- **Clock.** If the use case takes a `Clock` port, use its generated `MockClock` and return a fixed time **far from today** (for example `time.Date(2031, 3, 1, 12, 0, 0, 0, time.UTC)`), so timestamps can be checked exactly. If the fixed time is near the real date, a regression back to `time.Now()` passes or fails depending on the hour. If it calls `time.Now()` directly, that violates R-014. Bound the value (`before := time.Now()` … `after`), and file a task to inject a clock. Don't change production code in a test-only change.
- **Logs.** Never assert on log output (R-017). If the package has no `TestMain`, add one in `main_test.go` to discard logs: `log.Init(slog.New(slog.NewTextHandler(io.Discard, nil)))`, then `os.Exit(m.Run())`. Also check `go vet -tags=integration` on the package, because integration files may live in the same package. Never use `func init()`.

## Assertions

| Prove | How |
|---|---|
| The returned error | `require.ErrorIs(t, err, domain.ErrX)` on the sentinel. Don't use `Contains(err.Error(), "failed to …")` unless the message itself is the contract (shown to users). |
| What was written | `DoAndReturn` that snapshots the argument (`saved = *k`), then assert the fields the spec names. |
| What must not happen | Leave out the `EXPECT`. gomock fails the test on any unexpected call (for example, no role lookup after a denial, or no event after a failed commit). |
| Call order | `gomock.InOrder(...)`, only when order is a requirement (for example, publish after commit). |
| Writes join the transaction | Give the fake transaction a marked context (`context.WithValue`), and expect the repository call with exactly that context. |

**Use distinctive test data.** Pick values that are not defaults: `IsSandbox: true`, `RateLimit: 300` rather than `false` or `60`. Otherwise code that hard-codes the constant passes. Don't mirror the implementation field by field when the spec only promises a few of the fields. Don't assert on values you passed straight into a mock and read straight back.

## Prove-fail

For every behaviour, change one line of production code so the behaviour breaks. The named test must go red. Then restore the line.

1. Back up the file: `cp internal/usecases/<domain>/<file>.go <scratchpad>/<file>.<your-unique-tag>.orig`. The scratchpad can be shared between sessions, so give the backup a unique name.
2. Make one change, with the Edit tool or a small script in the scratchpad that applies one replacement, runs the test, restores the file and runs `cmp`. A script is faster for many mutations. Examples: `<=` to `<`, drop a call, swap `txCtx` for `ctx`, publish before the commit.
3. Run the test, and see **that test** fail with an assertion or gomock message. A compile error is not a red.
4. Restore the file from the backup and run `cmp` against it.

A mutation that no test catches means an assertion is weak or missing. Fix the test, not the mutation. List the mutations and the test each one broke in your PR.

## Report / PR test plan

Include:

- the behaviours, with their BR ids;
- the required rows (Recipe step 2), each with its test name or N/A and why;
- the mutations table;
- the bugs and questions you filed, with their ids;
- what you deliberately didn't test, and why;
- the exact commands you ran.

## Common mistakes

| Mistake | Fix |
|---|---|
| Reading the code first, then asserting what it does | Spec first (step 2). A test that matches the code also protects its bugs. |
| Encoding a suspicious behaviour "because that's what it does today" | Bug, question or finding (step 3), never an expectation. |
| Finding a bug in passing and only mentioning it | `bd create --type=bug` with file:line. Put the id in the test's skip. |
| Hand-writing a mock, or using testify `mock.Mock` | Add a `//go:generate` line (Mocks). |
| `make mocks` for one package | `go generate ./internal/usecases/<domain>/` |
| Asserting `err.Error()` text | `ErrorIs` on the sentinel. |
| Concluding "no tests exist" from file names | grep for the use case name (step 1). |
| Trusting coverage numbers quoted in a doc or spec | Measure: `go test -coverprofile=… && go tool cover -func=…` |
| A compile error counted as a red | Fix the mutation so it builds, then rerun. |
| Test data that equals the defaults (`false`, `0`, `60`) | Distinctive values (Assertions). |
| A fixed clock set to today | A date years away (Time and logging). |
| A permission test for a use case that takes no Principal | The gate is in `rpc_actions.go`. Note it under "not tested here" (step 2). |

## Environment notes

In an isolated worktree, `git` and compound shell commands (pipes, `&&`, `$(...)`) may be refused. Use one plain command per call with literal paths (`go -C <worktree> test …` works; shell variables inside `rtk proxy` may be refused), quote globs (zsh), and check that production files are restored with `cmp`, not `git diff`. The `rtk` hook shortens `go test` output to "N passed" and hides skips and failure detail. Use `rtk proxy go test -v …` when you need to see them.
