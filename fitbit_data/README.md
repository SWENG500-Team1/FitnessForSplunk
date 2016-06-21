# Setup

1. Set the `SPLUNK_HOME` environment variable to the root directory of your Splunk instance.
* Copy this whole `fitbit_data` folder to `$SPLUNK_HOME/etc/apps`.
* Open a terminal at `$SPLUNK_HOME/etc/apps/fitbit_data/bin/app`.
* Run `npm install`.
* Restart Splunk

# Adding an input

# TODO: Finish README.md

#1. From Splunk Home, click the Settings menu. Under **Data**, click **Data inputs**, and find `Random Numbers`, the input you #just added. **Click Add new on that row**.
#* Click **Add new** and fill in:
#    * `name` (whatever name you want to give this input)
#    * `min` (the minimum value for a random number)
#    * `max` (the maximum value for a random number)
#* Save your input, and navigate back to Splunk Home.
#* Do a search for `sourcetype=random_numbers` and you should see some events.
