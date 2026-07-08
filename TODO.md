# TODO

## Performance

### Replace `go-enry` with generated GitHub Linguist lookup tables

**Status:** Backlog

#### Motivation

The current language detection relies on `go-enry`, which provides excellent coverage but comes with a noticeable binary size increase and additional startup overhead.

#### Goals

- [ ] Replace `go-enry` with compile-time generated lookup tables from GitHub Linguist.
- [ ] Generate extension → language mappings during build using `go:generate`.
- [ ] Handle ambiguous extensions with a small override/heuristic layer.
- [ ] Benchmark detection speed before and after migration.
- [ ] Verify language detection accuracy against the current implementation.

#### Expected Benefits

- Reduce upload binary size by approximately **8–9 MB**.
- Improve CLI startup time.
- Remove the external `go-enry` dependency.
- Retain support for ambiguous extensions using lightweight overrides (e.g. `.rs`, `.m`, `.h`).
- Maintain GitHub Linguist compatibility without runtime parsing.

#### Notes

- This is a performance optimization, **not a current priority**.
- Revisit after the language detection pipeline is stable across all supported editor plugins.
- Measure startup time and binary size before beginning the migration to establish a baseline.
