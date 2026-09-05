" plugin/taka-time.vim
"
" ROLE: This is the Neovim entry point shim.
" It is NOT the classic Vim plugin (see vimPlugin/TakaTime/ for that).
" This file is loaded by Neovim's plugin loader and delegates all work
" to the Lua module at lua/taka-time/init.lua via require('taka-time').
"
" For classic Vim 8+/9 users: use vimPlugin/TakaTime/ instead.

" TakaTime — Classic Vim plugin (Vim 8+ / Vim 9)
" Tracks coding activity and syncs with the TakaTime Go binary.
" https://github.com/Rtarun3606k/TakaTime

if exists('g:loaded_takatime') | finish | endif
let g:loaded_takatime = 1

" User-configurable options with defaults
if !exists('g:takatime_sync_interval')
  let g:takatime_sync_interval = 120
endif

if !exists('g:takatime_binary')
  " Default binary location — override in vimrc if needed
  let g:takatime_binary = expand('~/.takatime/takatime-bin')
endif

augroup TakaTime
  autocmd!
  " Sync heartbeat on every file save
  autocmd BufWritePost * call takatime#heartbeat(expand('<afile>:p'), 1)
  " Track cursor activity (non-blocking, rate-limited internally)
  autocmd CursorMoved,CursorMovedI * call takatime#heartbeat(expand('%:p'), 0)
  " Track file open
  autocmd BufEnter,BufReadPost * call takatime#heartbeat(expand('<afile>:p'), 0)
augroup END

command! TakaTimeDashboard call takatime#open_dashboard()
command! TakaTimeStatus    echo takatime#get_status()