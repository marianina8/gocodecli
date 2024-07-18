Exercise: Add a new history command with a starDate flag

Objective: Create a new command, history, that parses the log file for historical data related to specific URL checks.  The command also takes urls as args.  Create a new flag, –startDate to take in a historical date to use as the start of the url(s) history

./healthcheck history http://www.google.com --startDate 03/25/2024

or

healthcheck.exe history http://www.google.com --startDate 03/25/2024


Requirements:
The tool should accept the history command.
Add a new –startDate flag to additionally parse given a historical start date to print health checks.
It should receive one or more URLs and only print history URL checks after the specified start date