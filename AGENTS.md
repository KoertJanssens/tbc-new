# Agent Notes

Notes for AI coding agents working in this repository (`github.com/wowsims/tbc` — TBC-classic WoW sim).

## Building & running sims

- Phase gear presets (P1-P5) live under `ui/<class>/<spec>/gear_sets/p{N}_*.gear.json`.
- Go toolchain is **not** preinstalled on all machines. Install via `brew install go`
  (go.mod pins `go 1.25.0` / toolchain 1.25.4, but a newer brew-installed Go works
  fine and will auto-fetch the pinned toolchain if needed).
- Proto Go bindings (`sim/core/proto/*.pb.go`) are gitignored and must be generated
  before the project will build:
  - `brew install protobuf` (installs `protoc`)
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10` (match the
    `google.golang.org/protobuf` version pinned in `go.mod`), and ensure
    `$(go env GOPATH)/bin` is on `PATH`
  - `protoc -I=./proto --go_out=./sim/core ./proto/*.proto`
- The item database (`assets/database/db.json`) is only compiled in with the
  `with_db` build tag (see `sim/core/database_load.go`). Without `-tags with_db`,
  `ItemsByID` is empty and any real gear set will panic with `No item with id: X`.
  **Always run sim tests/binaries with `-tags with_db`.**
- Example: `go test -tags with_db ./sim/warrior/dps/ -run <TestName> -v`

## Getting DPS numbers out of the test suites

- Existing test suites (e.g. `sim/warrior/dps/dps_test.go` `TestDpsWarrior`) compare
  results against checked-in `.results` fixture files (prototext format) and only
  log DPS via `t.Logf` on a mismatch — they print nothing on a pass.
- To get actual DPS numbers for a specific build/gear combo, write a small
  throwaway `_test.go` file that:
  1. Builds a `*proto.RaidSimRequest` via `core.WithSpec(...)` +
     `core.SinglePlayerRaidProto(...)`, using `core.GetGearSet(dir, name)` and
     `core.GetAplRotation(dir, name)` to load existing gear/APL JSON presets.
  2. Calls `core.RunRaidSim(rsr)` and prints `result.RaidMetrics.Dps.Avg` (and
     `.Stdev`).
  3. Delete the scratch test file afterward — don't leave one-off test files
     checked in.
- Sanity-check new results against the checked-in `.results` fixtures (e.g.
  `sim/warrior/dps/TestDpsWarrior.results`) using the same `RandomSeed: 101`.
  Iteration count differences of ~1-2% are expected/normal.
