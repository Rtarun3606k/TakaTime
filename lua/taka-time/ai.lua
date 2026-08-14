local M = {}
local core = require("taka-time.core")

local ok, cmp = pcall(require, "cmp")
if not ok then
  return M
end

-- Supported AI providers
local ai_sources = {
  copilot = true,
  codeium = true,
  supermaven = true,
}

cmp.event:on("confirm_done", function(evt)
  local entry = evt.entry
  if not entry or not entry.source then
    return
  end

  local provider = entry.source.name
  if not ai_sources[provider] then
    return
  end

  -- 1. Extract character count of accepted text
  local char_count = 0
  if entry.completion_item then
    local item = entry.completion_item
    if item.textEdit and item.textEdit.newText then
      char_count = #item.textEdit.newText
    elseif item.insertText then
      char_count = #item.insertText
    elseif item.label then
      char_count = #item.label
    end
  end

  -- 2. Get absolute file path
  local file_path = vim.fn.expand("%:p")

  -- Pass file_path instead of lang
  M.on_ai_accept(provider, char_count, file_path)
end)

function M.on_ai_accept(provider, char_count, file_path)
  -- Push detailed event to core queue with the 'file' key
  core.add_ai_event({
    provider = provider,
    char = char_count,
    file = file_path,
  })

  local config = require("taka-time.config")
  if config.options.debug then
    print(string.format("[TakaTime] AI Accepted: %s | Chars: %d | File: %s", provider, char_count, file_path))
  end
end

return M
