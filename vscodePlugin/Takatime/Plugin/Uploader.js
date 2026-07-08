const vscode = require("vscode");
const path = require("path");
const os = require("os");
const fs = require("fs");
const { spawn } = require("child_process");
const env = require("./Config");

function getCleanEditorName() {
  const rawName = vscode.env.appName;

  // Group all the standard Microsoft/OSS builds under one clean name
  if (
    rawName === "Visual Studio Code" ||
    rawName === "Visual Studio Code - Insiders" ||
    rawName === "Code - OSS" ||
    rawName === "VSCodium"
  ) {
    return "VS Code";
  }

  // If you want to track Cursor or Windsurf specifically, let them pass through
  if (rawName === "Cursor" || rawName === "Windsurf") {
    return rawName;
  }

  // The Fallback: If it is some random fork we don't know about, return Antigravity
  return "Antigravity";
}

/**
 * Prepares the arguments for the Go binary
 * @param {vscode.TextDocument} document
 * @param {string} mongoUri - We need to pass this explicitly now
 */
function getGoArgs(document, mongoUri) {
  const filePath = document.fileName;

  // Detect Project Name
  let projectName = "Unknown";
  const workspaceFolder = vscode.workspace.getWorkspaceFolder(document.uri);
  if (workspaceFolder) {
    projectName = workspaceFolder.name;
  }

  // 📦 CORRECTED FLAGS (Based on your help output)
  return [
    "-file",
    filePath,
    "-project",
    projectName,
    "-uri",
    mongoUri, // Passing URI as flag (required by your binary)
    "-duration",
    "120", // Sending a default heartbeat duration (optional, fixes "less than 0" error)
    "-editor",
    // "VsCode",
    getCleanEditorName(),
  ];
}

/**
 * Spawns the binary in the background
 * @param {vscode.TextDocument} document
 */
function spawnProcess(document) {
  const config = env.getConfig();
  if (!config || !config.MONGO_URI) return;

  // 👉 1. Locate Binary (UPDATED to match the new 'taka-upload' naming convention)
  const homeDir = os.homedir();
  const isWin = process.platform === "win32";
  const ext = isWin ? ".exe" : "";

  // Notice: 'taka-upload', not 'taka-uploader'!
  const binName = `taka-upload-${env.CURRENT_VERSION}${ext}`;

  const binaryPath = path.join(homeDir, ".takatime", "bin", binName);

  if (!fs.existsSync(binaryPath)) {
    console.warn(
      `TakaTime: Binary not found at ${binaryPath}, skipping upload.`,
    );
    return false;
  }

  // 2. Spawn (Fire & Forget)
  try {
    // We pass 'config.MONGO_URI' to our helper now
    const args = getGoArgs(document, config.MONGO_URI);

    const child = spawn(binaryPath, args, {
      detached: !isWin,
      stdio: "ignore", // Change to "inherit" if you want to see logs in Debug Console
      env: {
        ...process.env,
      },
      windowsHide: true,
    });
    child.unref();
    console.log(`TakaTime: Uploading ${path.basename(document.fileName)}...`);
    return true;
  } catch (err) {
    console.error("TakaTime: Failed to spawn process", err);
    return false;
  }
}

module.exports = { spawnProcess };
