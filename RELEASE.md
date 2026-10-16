# Release checklist

## Source/version
- [ ] `VERSION`, CLI version and release tag agree.
- [ ] Go 1.26 baseline and changelog are current.

## Calendar/time gate
- [ ] trigger offset must match the configured IANA timezone at the trigger instant.
- [ ] elapsed hours and calendar days retain distinct semantics.
- [ ] DST-gap calendar results fail visibly.
- [ ] weekend names are valid/unique and do not cover all seven days.
- [ ] holiday dates are valid/unique.
- [ ] lifecycle cannot affect deadline arithmetic.
- [ ] no jurisdiction/holiday rules are hidden in code.

## Engineering gate
- [ ] `gofmt` clean.
- [ ] `go vet ./...` passes.
- [ ] `go test -race -count=1 ./...` passes.
- [ ] native/version and Linux ARM64 builds pass.
- [ ] Bafflewick and calendar-boundary acceptance pass.

## Evidence
- [ ] all PRD criteria map to executable evidence.
- [ ] generated EDD is pinned to exact source commit.
- [ ] generated comparison has no unexplained mismatch.
