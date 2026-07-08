local M = {}
local config = require("taka-time.config")
local uv = vim.uv or vim.loop

---@enum BinaryEnum
M.BinaryEnum = {
	UPLOAD = "taka-upload",
	DASHBOARD = "taka-dashboard",
}

function M.get_binary_path(binary)
	local plugin_root = vim.fn.fnamemodify(debug.getinfo(1).source:sub(2), ":h:h:h")
	local bin_path = plugin_root .. "/" .. binary

	local os_name = uv.os_uname().sysname:lower()
	if string.match(os_name, "windows") ~= nil then
		bin_path = bin_path .. ".exe"
	end

	return bin_path
end

function M.get_version_file_path()
	local plugin_root = vim.fn.fnamemodify(debug.getinfo(1).source:sub(2), ":h:h:h")
	return plugin_root .. "/.version"
end

-- Helper: Read version from disk
function M.get_installed_version(binary)
	local f = io.open(M.get_version_file_path(), "r")
	if not f then
		return nil
	end

	local content = f:read("*a")
	f:close()

	local ok, data = pcall(vim.json.decode, content)
	if ok and type(data) == "table" then
		return data[binary]
	end

	return nil
end

-- Helper: Write version to disk
function M.write_installed_version(binary, version)
	local path = M.get_version_file_path()
	local data = {}

	-- read existing
	local f = io.open(path, "r")
	if f then
		local content = f:read("*a")
		f:close()
		local ok, decoded = pcall(vim.json.decode, content)
		if ok and type(decoded) == "table" then
			data = decoded
		end
	end

	-- update this binary only
	data[binary] = version

	-- write back
	local wf = io.open(path, "w")
	if wf then
		wf:write(vim.json.encode(data))
		wf:close()
	end
end

-- Helper: Detect OS/Arch
local function get_os_info()
	local uname = uv.os_uname()
	local os = uname.sysname:lower()
	local arch = uname.machine:lower()
	if os == "linux" then
		os = "linux"
	elseif os == "darwin" then
		os = "darwin"
	elseif string.match(os, "windows") ~= nil then
		os = "windows"
	else
		return nil, nil
	end
	if arch == "x86_64" then
		arch = "amd64"
	end
	if arch == "aarch64" then
		arch = "arm64"
	end
	return os, arch
end
--------------------------------------------------------------------------------------------
---binary downloads
--------------------------------------------------------------------------------------------

function M.ensure_binary(binary)
	local bin_path = M.get_binary_path(binary)
	local target_ver = config.options.binary_version
local current_ver = M.get_installed_version(binary)
	if vim.fn.filereadable(bin_path) == 1 and current_ver == target_ver then
		local size = vim.fn.getfsize(bin_path)
		if size > 50000 then
			return
		end
	end
	-- 2. Update logic
	if current_ver and current_ver ~= target_ver then
		print(string.format("[Taka] Updating %s %s -> %s...", binary, current_ver, target_ver))
	end

	local os_name, arch = get_os_info()
	if not os_name then
		print("[Taka] Auto-install not supported for this OS.")
		return
	end
	local ext = os_name == "windows" and ".exe" or ""

	local url = string.format(
		"https://github.com/Rtarun3606k/TakaTime/releases/download/%s/%s-%s-%s%s",
		target_ver,
		binary,
		os_name,
		arch,
		ext
	)

	-- 3. Delete old binary and download new one
	if vim.fn.filereadable(bin_path) == 1 then
		os.remove(bin_path)
	end
	print("taka url downloading url " .. url)
	print("[Taka] Downloading " .. binary .. " " .. target_ver .. "...")
	vim.fn.system({ "curl", "-L", "-o", bin_path, url })
	if os_name ~= "windows" then
		vim.fn.system({ "chmod", "+x", bin_path })
	end

	-- 4. Update version file
  M.write_installed_version(binary, target_ver)	-- 1. Check if we are already up to date
	print("[Taka] Successfully installed " .. binary .. " " .. target_ver)
end

-----------------------------------------------------------------------------------------
-----------------------------------------------------------------------------------------

function M.get_os()
	local os_name = uv.os_uname().sysname

	print(os_name, " os name")
	return os_name
end

-- Helper: Ensure path ends with a separator for safe boundary checking
local function ensure_trailing_slash(path)
	-- If it doesn't already end with a slash (or backslash for Windows), add one
	if not path:match("[\\/]$") then
		return path .. "/"
	end
	return path
end

-- Helper: Check if the current directory is inside an ignored repository
function M.is_ignored(current_dir)
	local ignore_list = config.options.ignore_repos or {}

	-- Add a slash to the current directory for safe comparison
	local safe_current = ensure_trailing_slash(current_dir)

	for _, ignored_path in ipairs(ignore_list) do
		-- Add a slash to the ignored path as well
		local safe_ignored = ensure_trailing_slash(ignored_path)

		-- Now it strictly checks folder boundaries, not just partial strings!
		if vim.startswith(safe_current, safe_ignored) then
			return true
		end
	end

	return false
end

return M
