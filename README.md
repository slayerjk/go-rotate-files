# go-rotate-files
Simple file rotation(log, backups, etc) in specified dir

<h2>Description</h2>

App takes absolute path of dir or dirs(separated by coma), and rotate files in it/them.

Main logic is to keep only N number of most recent files.

App ignores subdirs and files in them.

<h2>Flags</h2>

```
("log-dir", logsPathDefault, "set custom log dir")
("keep-logs", 7, "set number of logs to keep after rotation")
("d", "NONE", "REQUIRED, abs path of dir or dirs, separeted by coma")
("r", -1, "REQUIRED, most recent files to keep")
```