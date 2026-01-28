# go-rotate-files
Simple file rotation(log, backups, etc) in specified dir

<h2>Description</h2>

App takes absolute path of dir or dirs, separated by coma, and rotate files in it/them.

Main logic is to keep only N number of most recent files.

App ignores subdirs and files in them.