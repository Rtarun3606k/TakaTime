# TakaTime

**TakaTime** is a lightweight, Bring-Your-Own-Database (BYOD) alternative to WakaTime. Built for developers who want deep insights into their coding habits without sending a single byte of telemetry to third-party servers.

This extension seamlessly integrates the TakaTime Go-based engine directly into your editor. Whether you are using **VS Code, VSCodium, Code - OSS**, or AI forks like **Cursor** and **Windsurf**, TakaTime tracks your coding activity, project time, and language usage securely.

<div align="center">

<img width="1080" height="608" alt="output" src="https://github.com/user-attachments/assets/2179e786-ee06-4fea-b845-53759abb60a1" />

</div>

## Key Features

- **100% Privacy & BYOD:** All telemetry is routed exclusively to your personal MongoDB instance. You own your data.
- **Universal Compatibility:** Works flawlessly across the entire Open VSX ecosystem, intelligently grouping your editor stats under a clean `VS Code` umbrella (while natively tracking forks like Cursor).
- **Integrated TUI Dashboard:** View your developer productivity analytics and 365-day heatmaps natively inside the built-in terminal.
- **Smart Heartbeat Detection:** Accurately calculates true coding duration, automatically pausing when you step away from the keyboard.
- **Cross-Platform Parity:** Seamlessly syncs your metrics across VS Code, JetBrains IDEs, and Neovim.

---

## Installation & Setup

Because TakaTime is a Bring-Your-Own-Database engine, you need to provide your own MongoDB URI to start tracking.

### 1. Install the Extension

Install **TakaTime** from the Open VSX Registry or the standard VS Code Marketplace.

### 2. Configure Your Database

1. Open your Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P`).
2. Type and select: **`TakaTime: Configure Database URI`**
3. Paste your MongoDB connection string (e.g., `mongodb+srv://<user>:<password>@cluster...`).
4. Hit Enter. The extension will securely save your URI and immediately begin tracking your heartbeats!

---

## Commands

Access these commands via the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`):

| Command                            | Description                                                                      |
| ---------------------------------- | -------------------------------------------------------------------------------- |
| `TakaTime: Open Dashboard`         | Spawns the interactive Go TUI dashboard directly in your terminal pane.          |
| `TakaTime: Configure Database URI` | Sets or updates your MongoDB connection string.                                  |
| `TakaTime: Check Status`           | Verifies your connection to the database and displays the current session stats. |

---

## GitHub Profile Integration

Show off your coding productivity! TakaTime includes a built-in automated workflow to generate a beautiful, cross-platform **365-Day Contribution Heatmap** directly on your GitHub Profile `README.md`.

Check out the [main TakaTime repository](https://github.com/Rtarun3606k/TakaTime) to set up the GitHub Action and generate your own stats card!

---

## Privacy Policy

TakaTime is designed from the ground up to protect your data.

- **Zero Analytics Servers:** We do not host, intercept, or proxy your data.
- **Direct Connection:** Your editor talks directly to the MongoDB URI you provide.
- **No Code Uploaded:** Only metadata (language, project name, editor, timestamp) is tracked. Your source code is never read or transmitted.

For the full open-source codebase, visit the GitHub repository.

## How to Display Stats on GitHub

To show the graph on your GitHub profile, you need to set up the **Taka-Report** action in your profile repository.

1.  Go to your GitHub Profile repository (`username/username`).
2.  Create a file: `.github/workflows/takatime.yml`.
3.  Copy the workflow configuration from the [Official Repository](https://github.com/Rtarun3606k/TakaTime).

## Privacy Policy

**TakaTime does not send your data to any third-party servers.**

- Your coding activity is sent **only** to the MongoDB URI you provide.
- The extension downloads a helper binary from the official [TakaTime GitHub Releases](https://github.com/Rtarun3606k/TakaTime/releases).
- No telemetry is collected by the extension author.

## Links

- [GitHub Repository](https://github.com/Rtarun3606k/TakaTime)
- [Report an Issue](https://github.com/Rtarun3606k/TakaTime/issues)

---

**Enjoying TakaTime?** ⭐ Star the [repo on GitHub](https://github.com/Rtarun3606k/TakaTime)!
