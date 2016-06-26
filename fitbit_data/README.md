# Setup

1. Set the `SPLUNK_HOME` environment variable to the root directory of your Splunk instance.
* Copy this whole `fitbit_data` folder to `$SPLUNK_HOME/etc/apps`.
* Open a terminal at `$SPLUNK_HOME/etc/apps/fitbit_data/bin/app`.
* Run `npm install`.
* Restart Splunk server

# Adding an input
1. From Splunk Home, click the Settings menu. Under **Data**, click **Data inputs**, and find `Fitbit Data`, the input you #just added. **Click Add new on that row**.
* Click **Add new** and fill in:
    * `username` (Username of the account)
    * `date` (Date of data to be retrieved in the format, **YYYY-MM-DD**)
    * `token_json` (JSON token returned from the Fitbit API after authentication)
* Save your input, and navigate back to Splunk Home.
* Do a search for `sourcetype=fitbit_data` and you should see some events.
