# ota-status-service

Demo service. It reports the progress of over-the-air update campaigns to
internal dashboards and flags a campaign as halted when too many vehicles
fail.

Everything here is invented for a demo. The VINs are not real and this is not
Lucid code.

## Run

```
go run .
curl localhost:8081/v1/campaigns
```

## Test

```
go test ./...
```

## Pipeline

`.github/workflows/ci.yml` runs `go vet`, the tests and `govulncheck`.

## Open finding

`security/FINDING-002.md` is the same advisory that is open on
`connectivity-gateway`. The scanner reported it here because the module
version matches.

## Conventions

- Never push to `main`. Open a pull request.
- A pull request description lists what changed, the test results and the
  re-scan result.
