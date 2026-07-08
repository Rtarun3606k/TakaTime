# Installation & Development Setup

## Prerequisites

| Component | Version             |
| --------- | ------------------- |
| Go        | >= 1.25.3           |
| MongoDB   | >= 6.0 recommended  |
| Node.js   | >= 18 (recommended) |
| VS Code   | >= 1.80.0           |
| JDK       | 17 or 21            |
| Neovim    | >= 0.9 recommended  |

---

# Core Go Backend

## Install Dependencies

```bash
go mod download
```

## Build

```bash
go build -trimpath -ldflags="-s -w"
```

## Run

```bash
go run .
```

---

# VS Code Extension

## Install Dependencies

```bash
npm install
```

## Lint

```bash
npm run lint
```

## Run Tests

```bash
npm test
```

## Launch Extension Development Host

```bash
code .
```

Press `F5` inside VS Code to launch the Extension Development Host.

---

# Neovim Plugin (LazyVim)

This plugin is configured using the **LazyVim framework**.

## Requirements

- Neovim >= 0.9 recommended

## Installation

Add the plugin to your Lazy.nvim configuration:

```lua
return {
  "Rtarun3606k/TakaTime",
  lazy = false,
  opts = {
    -- debug = false,  -- set to true for verbose logging
  },
}
```

## Initial Setup

After installation, run:

```vim
:TakaInit
```

This will prompt you to enter your MongoDB URI securely.

## Available Commands

```vim
:TakaInit
:TakaDash
:TakaStatus
:TakaIgnore
:TakaTrack
```


---

# JetBrains Plugin

## Requirements

* JDK 17 or JDK 21
* Gradle

## Run IDE Sandbox

```bash
./gradlew runIde
```

## Build Plugin

```bash
./gradlew buildPlugin
```
