# solution-center-metrics

## Usage

Simply run the executable.

## Maintaining

This tool grabs an existing CSV file from ServiceNow and sorts it, including only the data we want. If the order of columns changes, all that needs to be done is edit the coresponding column number to the switch case in ParseData().

Remember to shift the column number down one to account for the index starting at 0 :)