# Setup

1. Set the `SPLUNK_HOME` environment variable to the root directory of your Splunk instance.
* Copy this whole `microsoft_data` folder to `$SPLUNK_HOME/etc/apps`.
* Open a terminal at `$SPLUNK_HOME/etc/apps/microsoft_data/bin/app`.
* Run `npm install`.
* Restart Splunk server

# Adding an input
1. From Splunk Home, click the Settings menu. Under **Data**, click **Data inputs**, and find `Fitbit Data`, the input you #just added. **Click Add new on that row**.
* Click **Add new** and fill in:
    * `fullname` (first and last of the account)
    * `username` (Username of the account)
    * `startdate` (Date of data to be retrieved in the format using 24 hours time, **YYYY-MM-DD-HH**  Ex: 2016-07-13-18)
    * `token_json` (JSON token returned from the Microsoft Health API after authentication)
* Save your input, and navigate back to Splunk Home.
* Do a search for `sourcetype=microsoft_data` and you should see some events.



----------------------------------------------------------------------------------
#Testing / Clear Data from Splunk (Windows)
Reset all of the data to start collecting fresh data
1. Open the Command line and Navigate to `$SPLUNK_HOME/bin/` (c:\Program Files\Splunk\bin)
2. Enter the following commands
  * splunk stop
  * splunk clean eventdata -index yourindex
  * splunk start


#Error Information in Splunk
Navigate to the following directory to view the log files: `C:\Program Files\Splunk\var\log\splunk`
Contains information relating to input failures: splunkd.log
