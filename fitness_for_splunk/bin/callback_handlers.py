#
# http request handlers
#

from splunk import auth, search
import splunk.rest 
import splunk.bundle as bundle
import httplib2, urllib, os, time
import base64
#import utils
#import logging as logger
#import telnetlib

# set our path to this particular application directory (which is suppose to be <appname>/bin)
app_dir = os.path.dirname(os.path.abspath(__file__))

# define the web content directory (needs to be <appname>/web directory)
web_dir = app_dir + "/web"

class fitbit_callback(splunk.rest.BaseRestHandler):

    """
    Main endpoint
    """
    def handle_GET(self):
        
        clientId = "227MVJ"
        clientSecret = "df8009bd0ddcb975f9a812e3587e54dd"
        callback_url = "https://localhost:8089/services/fitness_for_splunk/fitbit_callback"
        #encoded_id_secret = base64.b64encode( (clientId + ':' + clientSecret) )
        #authCode = '8cb9534b46ee84c9082f9cd2558416006fc2d11e'
        
        """
        Parse authorization code from URL query string
        """
        # Check for existence of 'code' parameter
        for key in self.request['query']:
            if key == 'code':
                authCode = self.request['query']['code']
            else:
                authCode = None
                
        if authCode is None:
            self.response.setStatus(400)
            self.response.setHeader('content-type', 'text/html')
            self.response.write("Bad request: no authorization code")
        else:
            """
            Exchange authorization code for token
            """
            encoded_id_secret = base64.b64encode( (clientId + ':' + clientSecret) )
            
            url = 'https://api.fitbit.com/oauth2/token'
            authHeader_value = ('Basic ' + encoded_id_secret)
            headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
            data = {'clientId': clientId, 'grant_type': 'authorization_code', 'redirect_uri': callback_url, 'code': authCode}
            body = urllib.urlencode(data)
            
            http = httplib2.Http()
            resp, cont = http.request(url, 'POST', headers=headers, body=body)
            
            """
            Store token in Splunk password store
            """
            # TODO
            
            
            
            self.response.setStatus(200)
            #self.response.setHeader('content-type', 'text/html')
            self.response.setHeader('content-type', 'application/json')
            self.response.write(cont)

    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET


class google_callback(splunk.rest.BaseRestHandler):
    '''
    Status endpoint
    '''
    def handle_GET(self):

        # set output params
        self.response.setStatus(200)
        self.response.setHeader('content-type', 'text/html')
        self.response.write('Google Callback')
    

    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET


class microsoft_callback(splunk.rest.BaseRestHandler):
    '''
    Main endpoint
    '''
    
    def handle_GET(self):
        
        # set output params
        self.response.setStatus(200)
        self.response.setHeader('content-type', 'text/html')
        self.response.write('Microsoft Callback')
    
    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET