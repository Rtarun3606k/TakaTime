## [2.2.8] - 2026-09-04

### Added

- changed version



## [2.2.7] - 2026-08-31

### Added

- changed version

## [2.2.6] - 2026-05-21

### Added

- logo for plugin

## [2.2.5] - 2026-05-21

### Fixed

- Fixed JetBrains document listener lifecycle handling to prevent invalid listener removal exceptions when closing editors or split tabs.

## [2.2.4] - 2026-05-20

### Added

- Integrated TUI Dashboard: TakaTime now launches its analytics dashboard directly inside the IDE's built-in terminal, eliminating external windows and context switching.

- Native Status Bar Menu: Added a native popup menu accessible from the status bar. Users can quickly configure the MongoDB URI or launch the dashboard directly from the IDE.

- Dynamic IDE Detection: The telemetry engine now detects and logs the exact JetBrains IDE in use, including IntelliJ IDEA, PyCharm, WebStorm, and others.

### Changed

- Off-Thread Downloading: Moved version checking and binary downloading fully off the UI thread to prevent IDE freezes and improve responsiveness.

- Dynamic Heartbeat Duration: Replaced hardcoded telemetry intervals with a dynamic heartbeat calculator that accurately tracks active coding time while pausing during inactivity.

### Fixed

- Cross-Platform Version Sync: Implemented semantic version-aware synchronization to prevent VS Code and JetBrains extensions from repeatedly downgrading each other.

- Accurate Project Detection: Fixed an issue where temporary files such as `Main.kt` were incorrectly detected as project names. TakaTime now correctly resolves the workspace root directory.
