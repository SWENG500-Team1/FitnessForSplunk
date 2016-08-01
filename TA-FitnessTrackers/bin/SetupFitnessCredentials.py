import splunk.rest
import splunk.auth
import splunk.admin as admin
import splunklib.client as client

appName = "TA-FitnessTrackers"
realmGoogle = "Google"
realmFitBit = "FitBit"
stanza = "apiCredentials"
googleClientID = "google_client_id"
googleClientSecret = "google_client_secret"
fitBitClientID = "fitbit_client_id"
fitBitClientSecret = "fitbit_client_secret"

"""Handles REST requests to create app credentials."""
class FitnessTrackerSetup(splunk.rest.BaseRestHandler):
    def setup(self):
        if self.requestedAction == admin.ACTION_CREATE:
            for arg in ['google_client_id', 'google_client_secret',
                        'fitbit_client_id', 'fitbit_client_secret']:
                self.supportedArgs.addOptArg(arg)

    def handle_GET(self):
        c = client.connect(app=appName,
                           username="Splunk",
                           password=splunk.getSessionKey())

        #Clear existing passwords
        for password in c.storage_passwords:
            c.storage_passwords.delete(password.name)

        #Write new values to the password store.
        c.storage_passwords.create(username=self.callerArgs.data[googleClientID],
                                   password=self.callerArgs.data[googleClientSecret],
                                   realm=realmGoogle)
        c.storage_passwords.create(username=self.callerArgs.data[fitBitClientID],
                                   password=self.callerArgs.data[fitBitClientSecret],
                                   realm=realmFitBit)

    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET
