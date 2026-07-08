local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")
-- Ensure this file is actually named 'lua/taka-time/core.lua'
local core = require("taka-time.core")
local storage = require("taka-time.storage")

function M.setup(opts)
  -- 1. Setup Config
  config.setup(opts)

  -- 2. AUTOMATICALLY LOAD SECRET
  if config.options.mongo_uri == "" then
    local saved_uri = storage.load_secret()
    if saved_uri then
      config.options.mongo_uri = saved_uri
    end
  end
  ----------------------------------------------------------------------------------------
  -- 3. Create Commands

  -- TakaInit command
  vim.api.nvim_create_user_command("TakaInit", function()
    local uri = vim.fn.inputsecret("Enter your Mongo URI: ")
    if uri and uri ~= "" then
      storage.save_secret(uri)
      config.options.mongo_uri = uri
      print("TakaTime: URI Saved!")
    else
      print("TakaTime: No URI Entered")
    end
  end, {})

  -------------------------------------------------------------------------------------
  ---taka dashboard
  -------------------------------------------------------------------------------------
  ---taka dashboard
  vim.api.nvim_create_user_command("TakaDash", function()
    -- 1. Ensure the user has actually configured their URI
    local uri = config.options.mongo_uri
    if not uri or uri == "" then
      print("TakaTime: Missing MongoDB URI. Please run :TakaInit first!")
      return
    end

    -- 2. Calculate the size of the floating window (90% of the Neovim screen)
    local width = vim.o.columns
    local height = vim.o.lines
    local win_width = math.floor(width * 0.9)
    local win_height = math.floor(height * 0.9)

    -- Calculate perfectly centered coordinates
    local row = math.floor((height - win_height) / 2)
    local col = math.floor((width - win_width) / 2)

    -- 3. Create a temporary scratch buffer
    local buf = vim.api.nvim_create_buf(false, true)

    -- 4. Open the floating window with a rounded border
    local win = vim.api.nvim_open_win(buf, true, {
      relative = "editor",
      width = win_width,
      height = win_height,
      row = row,
      col = col,
      style = "minimal",
      border = "rounded",
    })

    local bin_path = utils.get_binary_path(utils.BinaryEnum.DASHBOARD)

    -- 5. Construct the command and launch the terminal!
    -- NOTE: Make sure 'taka-dash' is either in your system PATH,
    -- or provide the absolute path to the binary here.
    local cmd = string.format("%s --MongoDBString \"%s\"", bin_path, uri)

    vim.fn.jobstart(cmd, {
      term = true,
      on_exit = function()
        -- When you press 'q' or 'esc' in the dashboard, it exits the Go binary.
        -- This hook automatically closes the Neovim floating window behind it!
        -- if vim.api.nvim_win_is_valid(win) then
        --   vim.api.nvim_win_close(win, true)
        -- end
      end,
    })

    -- 6. Automatically enter Insert mode so Bubble Tea can instantly read your keystrokes (j, k, s, q)
    vim.cmd("startinsert")
  end, {})

  ----------------------------------------------------------------------------------
  -- TakaStatus command
  vim.api.nvim_create_user_command("TakaStatus", function()
    if config.options.mongo_uri and config.options.mongo_uri ~= "" then
      print("TakaTime is configured and running.")
    else
      print("TakaTime is NOT configured. Run :TakaInit")
    end
  end, {})

  -----------------------------------------------------------------------------------
  -- Load the ignore list into RAM on startup
  config.options.ignore_repos = storage.load_ignore_list() or {}

  -- Command to IGNORE the current directory
  vim.api.nvim_create_user_command("TakaIgnore", function()
    local cwd = vim.fn.getcwd()

    -- Check if it's already in the RAM list
    if vim.tbl_contains(config.options.ignore_repos, cwd) then
      print("TakaTime: " .. cwd .. " is already being ignored.")
      return
    end

    -- Add to RAM (Instant)
    table.insert(config.options.ignore_repos, cwd)

    -- Save to Disk
    storage.save_ignore_list(config.options.ignore_repos)

    print("TakaTime: Ignored tracking for " .. cwd)
  end, {})

  ------------------------------------------------------------------------------------------
  -- Command to TRACK (Undo ignore) for the current directory
  vim.api.nvim_create_user_command("TakaTrack", function()
    local cwd = vim.fn.getcwd()
    local updated_list = {}
    local found = false

    -- Rebuild the list, skipping the current directory
    for _, repo in ipairs(config.options.ignore_repos) do
      if repo == cwd then
        found = true
      else
        table.insert(updated_list, repo)
      end
    end

    if not found then
      print("TakaTime: " .. cwd .. " is already being tracked.")
      return
    end

    -- Update RAM (Instant)
    config.options.ignore_repos = updated_list

    -- Save to Disk
    storage.save_ignore_list(config.options.ignore_repos)

    print("TakaTime: Resumed tracking for " .. cwd)
  end, {})

  -- 4. Ensures Binary Exists
  pcall(utils.ensure_binary, utils.BinaryEnum.UPLOAD)
  pcall(utils.ensure_binary, utils.BinaryEnum.DASHBOARD)
  -- 5. START TRACKING (The Fix)
  -- We delegate all logic to core. This sets up CursorMoved, TextChanged, AND BufWritePost
  core.setup_listeners()

  -- Start the background sync timer
  core.start_timer()
end

return M
