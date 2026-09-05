" TakaTime autoload functions — loaded lazily by Vim on first call

let s:last_heartbeat = 0

function! takatime#heartbeat(file, is_write) abort
  " Rate-limit non-write heartbeats
  if !a:is_write
    let l:now = localtime()
    if (l:now - s:last_heartbeat) < g:takatime_sync_interval
      return
    endif
    let s:last_heartbeat = l:now
  endif

  if empty(a:file) | return | endif
  if !executable(g:takatime_binary)
    " Binary not found — silent fail to not disrupt editing
    return
  endif

  " Build args
  let l:args = [
    \ g:takatime_binary,
    \ '--file', shellescape(a:file),
    \ '--plugin', 'vim',
  \ ]
  if a:is_write
    call add(l:args, '--write')
  endif

  " Fire-and-forget via job (Vim 8+)
  if has('job')
    call job_start(l:args, {'stoponexit': ''})
  else
    " Fallback for older Vim: use system() in background
    call system(join(l:args, ' ') . ' &')
  endif
endfunction

function! takatime#open_dashboard() abort
  if executable(g:takatime_binary)
    call job_start([g:takatime_binary, '--dashboard'])
  else
    echohl WarningMsg
    echom '[TakaTime] Binary not found at: ' . g:takatime_binary
    echohl None
  endif
endfunction

function! takatime#get_status() abort
  return '[TakaTime] Active | Binary: ' . g:takatime_binary
endfunction