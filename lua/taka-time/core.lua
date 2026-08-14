local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")
local uv = vim.uv or vim.loop

-- STATE
local state = {
	last_event_time = os.time(),
	pending_duration = 0,
	ai_accepted_count = 0, -- Total count integer
	ai_events = {},        -- Array of { provider = "", char = 0, lang = "" }
	job_id = 0,
	timer = nil,
}

local TIMEOUT_SECONDS = 120

-- Public API to push AI events into queue
function M.add_ai_event(event)
	state.ai_accepted_count = state.ai_accepted_count + 1
	table.insert(state.ai_events, event)
end

-- Internal Upload Logic
local function attempt_upload()
	if state.job_id ~= 0 then
		return
	end

	-- Check if we have anything to flush (time OR AI activity)
	if state.pending_duration < (config.options.debounce_seconds or 2) and state.ai_accepted_count == 0 then
		return
	end

	-- Snapshot state
	local time_to_send = state.pending_duration
	local ai_count_to_send = state.ai_accepted_count
	local ai_events_to_send = state.ai_events

	-- Reset in-memory queue
	state.pending_duration = 0
	state.ai_accepted_count = 0
	state.ai_events = {}

	local file_path = vim.fn.expand("%:p")
	local project = vim.fn.fnamemodify(vim.fn.getcwd(), ":t")

	if utils.is_ignored(vim.fn.getcwd()) then
		return
	end

	-- Encode metadata array to JSON string
	local ai_json_metadata = vim.json.encode(ai_events_to_send)

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
		"-ai-accepted",
		tostring(ai_count_to_send),
		"-ai-metadata",
		ai_json_metadata, -- Pass JSON string to Go CLI
		"-editor",
		"NeoVim",
	}

	if config.options.debug then
		print(string.format("[Taka] Syncing %ds | AI Count: %d | AI Metadata: %s", time_to_send, ai_count_to_send, ai_json_metadata))
	end

	state.job_id = vim.fn.jobstart(cmd, {
		on_exit = function(_, code)
			state.job_id = 0
			if code ~= 0 then
				-- Fault Tolerance: Restore queue if network/process fails
				state.pending_duration = state.pending_duration + time_to_send
				state.ai_accepted_count = state.ai_accepted_count + ai_count_to_send
				for _, evt in ipairs(ai_events_to_send) do
					table.insert(state.ai_events, evt)
				end
			end
		end,
	})

	if state.job_id <= 0 then
		state.job_id = 0
		state.pending_duration = state.pending_duration + time_to_send
		state.ai_accepted_count = state.ai_accepted_count + ai_count_to_send
		for _, evt in ipairs(ai_events_to_send) do
			table.insert(state.ai_events, evt)
		end
	end
end

-----------------------------------------------------------------------------------
local function on_activity()
	local now = os.time()
	local diff = now - state.last_event_time

	if diff < TIMEOUT_SECONDS then
		state.pending_duration = state.pending_duration + diff
	end

	state.last_event_time = now
end

function M.setup_listeners()
	local group = vim.api.nvim_create_augroup("TakaTimeGroup", { clear = true })

	vim.api.nvim_create_autocmd({ "CursorMoved", "CursorMovedI", "TextChanged", "TextChangedI", "InsertEnter" }, {
		group = group,
		callback = on_activity,
	})

	vim.api.nvim_create_autocmd("BufWritePost", {
		group = group,
		callback = function()
			on_activity()
			attempt_upload()
		end,
	})

	vim.api.nvim_create_autocmd("VimLeavePre", {
		group = group,
		callback = M.on_exit,
	})
end

function M.clear_timer()
	if state.timer then
		state.timer:stop()
		state.timer:close()
		state.timer = nil
	end
end

function M.on_exit()
	M.clear_timer()

	local time_to_send = state.pending_duration
	local ai_count_to_send = state.ai_accepted_count
	local ai_events_to_send = state.ai_events

	if time_to_send <= 0 and ai_count_to_send <= 0 then
		return
	end

	state.pending_duration = 0
	state.ai_accepted_count = 0
	state.ai_events = {}

	local file_path = vim.fn.expand("%:p")
	local project = vim.fn.fnamemodify(vim.fn.getcwd(), ":t")

	if utils.is_ignored(vim.fn.getcwd()) then
		return
	end

	vim.fn.system({
		utils.get_binary_path(utils.BinaryEnum.UPLOAD),
		"-uri",
		config.options.mongo_uri,
		"-project",
		project,
		"-file",
		file_path,
		"-duration",
		tostring(time_to_send),
		"-ai-accepted",
		tostring(ai_count_to_send),
		"-ai-metadata",
		vim.json.encode(ai_events_to_send),
		"-editor",
		"NeoVim",
	})
end

function M.start_timer()
	M.clear_timer()

	state.timer = uv.new_timer()
	state.timer:start(
		1000,
		60000,
		vim.schedule_wrap(function()
			attempt_upload()
		end)
	)
end

return M
