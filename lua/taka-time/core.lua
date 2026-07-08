local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")

-- STATE
local state = {
	last_event_time = os.time(), -- When was the last keystroke?
	pending_duration = 0, -- Accumulated seconds to send
	job_id = 0,
}

-- TIMEOUT: If no activity for 2 mins, don't count that time gap.
local TIMEOUT_SECONDS = 120

-- Internal: The actual upload logic
local function attempt_upload()
	-- 1. Checks
	if state.job_id ~= 0 then
		return
	end -- Busy

	-- If we have very little data (e.g. < 2s), wait for more (debounce)
	if state.pending_duration < (config.options.debounce_seconds or 2) then
		return
	end

	-- 2. Snapshot data
	local time_to_send = state.pending_duration
	state.pending_duration = 0 -- Reset bucket

	-- 3. Prepare Args
	local file_path = vim.fn.expand("%:p")
	local project = vim.fn.fnamemodify(vim.fn.getcwd(), ":t")

	--  THE PRIVACY FILTER
	if utils.is_ignored(vim.fn.getcwd()) then
		state.pending_duration = 0 -- Reset bucket so time doesn't build up
		return
	end

	local cmd = {
		utils.get_binary_path(utils.BinaryEnum.UPLOAD),
		"-uri",
		config.options.mongo_uri,
		"-project",
		project,
		"-file",
		file_path,
		"-duration",
		tostring(time_to_send),
		"-editor",
		"NeoVim",
	}

	if config.options.debug then
		print("[Taka] Syncing " .. time_to_send .. "s...")
	end

	-- 4. Run Job
	state.job_id = vim.fn.jobstart(cmd, {
		on_exit = function(_, code)
			state.job_id = 0
			if code ~= 0 then
				-- Failure: Put time back in bucket
				state.pending_duration = state.pending_duration + time_to_send
				if config.options.debug then
					print("[Taka] Failed. Retrying.")
				end
			elseif config.options.debug then
				print("[Taka] Success.")
			end
		end,
	})
end

-----------------------------------------------------------------------------------
--  THE FIX: Only add time if activity happened recently
local function on_activity()
	local now = os.time()
	local diff = now - state.last_event_time

	-- Only count this time if the gap was small (less than timeout)
	-- If diff > 120s, it means you were away. We ignore that gap.
	if diff < TIMEOUT_SECONDS then
		state.pending_duration = state.pending_duration + diff
	end

	-- Reset the clock for the next event
	state.last_event_time = now
end

-- Setup Autocommands to detect REAL activity
function M.setup_listeners()
	local group = vim.api.nvim_create_augroup("TakaTimeGroup", { clear = true })

	vim.api.nvim_create_autocmd({ "CursorMoved", "CursorMovedI", "TextChanged", "TextChangedI", "InsertEnter" }, {
		group = group,
		callback = on_activity,
	})

	-- On Save, we treat it as activity AND trigger an upload
	vim.api.nvim_create_autocmd("BufWritePost", {
		group = group,
		callback = function()
			on_activity()
			attempt_upload()
		end,
	})
end

-------------------------------------------------------------------------------------
-- Public: Called on Exit
function M.on_exit()
	-- 1. Snapshot the time immediately
	local time_to_send = state.pending_duration

	-- Safety check: If nothing to send, quit
	if time_to_send <= 0 then
		return
	end

	-- Reset the global bucket so we don't double-send (good practice)
	state.pending_duration = 0

	-- 2. Prepare Args (Get context)
	local file_path = vim.fn.expand("%:p")
	local project = vim.fn.fnamemodify(vim.fn.getcwd(), ":t")

	--  THE PRIVACY FILTER
	if utils.is_ignored(vim.fn.getcwd()) then
		state.pending_duration = 0 -- Reset bucket so time doesn't build up
		return
	end

	-- 3. Flush (Synchronous System Call)
	-- We use the LOCAL 'time_to_send' variable here
	vim.fn.system({
		utils.get_binary_path(utils.BinaryEnum.UPLOAD),
		"-uri",
		config.options.mongo_uri,
		"-project",
		project,
		"-file",
		file_path,
		"-duration",
		tostring(time_to_send), -- <--- FIX: Use the snapshot
		"-editor",
		"NeoVim",
	})
end

function M.start_timer()
	local timer = vim.loop.new_timer()
	timer:start(
		1000, -- Wait 1s
		60000, -- Repeat every 60s
		vim.schedule_wrap(function()
			-- Note: We do NOT add time here.
			-- We only check if there is time waiting to be sent.
			attempt_upload()
		end)
	)
end

return M
