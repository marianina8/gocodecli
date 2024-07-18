Exercise: Enhanced monitoring table command (+Colors and +Spinner)

Objective: Enhance the monitoring table we created in Module 3, by adding color to the up/down status, a spinner that shows up prior to the table loading, and a new column for the last time checked.

./healthcheck monitor http://www.google.com --output="table"

or

healthcheck.exe monitor http://www.google.com --output="table"


Requirements:
The monitor command should display UP status in green and DOWN status in red in table format
The monitor command should shows spinner as the table is loading
The monitor command shows a new column for the table, Last time checked, which is continually updated
