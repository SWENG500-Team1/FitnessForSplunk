"""TA-GoogleFitness"""
import sys

from splunklib.modularinput import Scheme, Argument, Event, Script, EventWriter
import requests

class GoogleFitness(Script):
    """
    Returns a Modular Input Scheme.
    """
    def get_scheme(self):
        #Returns Scheme
        scheme = Scheme("Google Fitness")
        scheme.description = "Gets fitness data from Google Fitness API"
        scheme.use_external_validation = True
        scheme.use_single_instance = True

        client_session = Argument("client_session")
        client_session.title = "Client Session"
        client_session.data_type = Argument.data_type_string
        scheme.add_argument(client_session)

        return scheme

    """
    Recieves settings from splunk, uses them to get data from Google and outputs
    them to standard out.
    """
    def stream_events(self, inputs, ew):
        for input_name, input_item in inputs.inputs.iteritems():
            ew.log(EventWriter.INFO, "Making request for Data Sources")
            request = requests.get('https://www.googleapis.com/fitness/v1/users/me/dataSources/',
                                   Authorization='Bearer:'
                                   + input_item["client_session"])
            event = Event()
            event.stanza = input_name
            event.data = request.json()
            ew.write_event(event)

if __name__ == "__main__":
    sys.exit(GoogleFitness().run(sys.argv))
