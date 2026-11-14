# TakaTime — Classic Vim Plugin

Vim 8+ / Vim 9 VimScript plugin for TakaTime coding activity tracking.

## Installation

### Using vim-plug
\```vim
Plug 'Rtarun3606k/TakaTime', { 'rtp': 'vimPlugin/TakaTime' }
\```

### Using Vundle
\```vim
Plugin 'Rtarun3606k/TakaTime'
" After installing, set runtimepath:
set runtimepath+=~/.vim/bundle/TakaTime/vimPlugin/TakaTime
\```

### Manual
Copy the `plugin/` and `autoload/` folders into `~/.vim/`.

## Configuration

Add to your `.vimrc`:
\```vim
" Path to the TakaTime Go binary (default: ~/.takatime/takatime-bin)
let g:takatime_binary = expand('~/.takatime/takatime-bin')

" Heartbeat interval in seconds for non-write events (default: 120)
let g:takatime_sync_interval = 120
\```

## Commands

| Command | Description |
|---|---|
| `:TakaTimeDashboard` | Open TUI dashboard |
| `:TakaTimeStatus` | Show plugin status |

## How it works

On `BufWritePost`, the plugin calls the Go binary immediately.
On `CursorMoved`, it sends a heartbeat at most once per `sync_interval` seconds.
All calls are non-blocking via Vim's `job_start()` API (Vim 8+).