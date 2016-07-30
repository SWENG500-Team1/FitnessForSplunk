#
# http request handlers
#

from splunk import auth, search
import splunk.rest 
import splunk.bundle as bundle
import splunk.entity as entity
import splunklib.client as client

import httplib2, urllib, os, time
import base64
import json
import xml.etree.ElementTree as ET

from oauth2client.client import flow_from_clientsecrets
from apiclient.discovery import build

import jsonbyteify
#import splunkmethods

"""
NOTE: To use this script, the following Python library dependencies need 
to be installed in the Splunk Python environment:

pip install --upgrade google-api-python-client

google-api-python-client
uritemplate
six
httplib2
oauth2client
simplejson
pyasn1
pyasn1-modules
rsa

jsonbyteify
"""

# set our path to this particular application directory (which is suppose to be <appname>/bin)
app_dir = os.path.dirname(os.path.abspath(__file__))

# define the web content directory (needs to be <appname>/web directory)
web_dir = app_dir + "/web"

class fitbit_callback(splunk.rest.BaseRestHandler):

    """
    Fitbit OAuth2 Callback Endpoint
    """
    def handle_GET(self):
        
        # Parse authorization code from URL query string
        authCode = None
        errorCode = None
        queryParams = self.request['query']
        for param in queryParams:
            if param == 'code':
                authCode = queryParams['code']
        
        if authCode is None:
            #TODO: Redirect to Error Authorizing Page
            # No Authorization Code, return 400
            self.response.setStatus(400)
            self.response.setHeader('content-type', 'text/html')
            self.response.write("Bad request: No authorization code")
            return
        
        # Pull Client ID and Secret from password store
        #TODO: Don't hardcode username and password
        c = client.connect(host="localhost", port="8089", username="admin", password="root")
        c.namespace.owner = "nobody"
        #c.namespace.app = "fitness_for_splunk"
        c.namespace.app = "TA-FitnessTrackers"
        passwords = c.storage_passwords
                
        # Look for Client Secret
        password = None
        for entry in passwords:
            if entry.realm == "fitbit":
                password = entry
                
        if password is None:
            #TODO: Redirect to Error Authorizing Page
            # No password configured for Fitbit, return 500
            self.response.setStatus(500)
            self.response.setHeader('content-type', 'text/html')
            self.response.write("Server Error: No password")
            return
        
        # Exchange authorization code for OAuth2 token
        http = httplib2.Http()
        clientId = password.username
        clientSecret = password.clear_password
        callback_url = "https://localhost:8089/services/fitness_for_splunk/fitbit_callback"
        auth_uri = "https://www.fitbit.com/oauth2/authorize"
        token_url = "https://api.fitbit.com/oauth2/token"
        
        encoded_id_secret = base64.b64encode( (clientId + ':' + clientSecret) )
        authHeader_value = ('Basic ' + encoded_id_secret)
        headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
        data = {'clientId': clientId, 'grant_type': 'authorization_code', 'redirect_uri': callback_url, 'code': authCode}
        body = urllib.urlencode(data)
        resp, content = http.request(token_url, 'POST', headers=headers, body=body)
        token_json = jsonbyteify.json_loads_byteified(content)
        
        user_id = token_json['user_id']
        access_token = token_json['access_token']
        
        # Get User Profile JSON
        userProfile_url = 'https://api.fitbit.com/1/user/' + user_id + '/profile.json'
        authHeader_value = ('Bearer ' + access_token)
        headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
        body = ""
        resp, content = http.request(userProfile_url, 'GET', headers=headers, body=body)
        profile_json = jsonbyteify.json_loads_byteified(content)
        full_name = profile_json['user']['fullName']
        
        # Store id, name, and token in KV store
        #TODO: Update Existing entries in KV Store
        
        collection_name = "fitbit_tokens"
        if collection_name in c.kvstore:
            c.kvstore.delete(collection_name)

        # Let's create it and then make sure it exists    
        c.kvstore.create(collection_name)
        
        kvstore = c.kvstore['fitbit_tokens']
        kv_jsonstring = json.dumps({"id": user_id, "name": full_name, "token": token_json})
        kvstore.data.insert(kv_jsonstring)
        
        #TODO: Redirect to Success Page
        self.response.setStatus(200)
        self.response.setHeader('content-type', 'text/html')
        #self.response.setHeader('content-type', 'application/json')
        self.response.write(str(kv_jsonstring))
        
    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET
    
    
class google_callback(splunk.rest.BaseRestHandler):

    """
    Google OAuth2 Callback Endpoint
    """
    def handle_GET(self):
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
            # TODO: Pull Client ID and Secret from password store?
            google_secret_filepath = "C:\\Program Files\\Splunk\\etc\\apps\\fitness_for_splunk\\bin\\google_secret.json"
            
            # Create flow object
            flow = flow_from_clientsecrets(google_secret_filepath, scope='https://www.googleapis.com/auth/fitness.activity.read https://www.googleapis.com/auth/fitness.body.read https://www.googleapis.com/auth/userinfo.profile', redirect_uri='https://localhost:8089/services/fitness_for_splunk/google_callback')
            
            # Set access_type to 'offline' to retrieve Refresh token
            flow.params['access_type'] = 'offline'
            
            #auth_uri = flow.step1_get_authorize_url()
            #redirect user to auth_uri
            
            """
            Exchange authorization code for token
            """
            credentials = flow.step2_exchange(authCode)
            http_auth = credentials.authorize(httplib2.Http())
            
            """
            Use token to retrieve Full Name
            """
            # Build service object
            service = build('plus', 'v1', http=http_auth)
            profile = service.people().get(userId='me').execute()
            
            user_id = profile['id']
            name = profile['name']
            
            # Concat UserId and FullName for storage in password store
            username = user_id + "////" + name['givenName'] + "." + name['familyName']
            realm = "google"
           
            """
            Store token in Splunk password store
            """
            # TODO
            
            """
            Write Response
            """
            self.response.setStatus(200)
            self.response.setHeader('content-type', 'text/html')
            #self.response.setHeader('content-type', 'application/json')
            self.response.write(username)
            
    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET


class microsoft_callback(splunk.rest.BaseRestHandler):
    """
    Microsoft OAuth2 Callback Endpoint
    """
    def handle_GET(self):
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
            # TODO: Pull Client ID and Secret from password store?
            microsoft_secret_filepath = "C:\\Program Files\\Splunk\\etc\\apps\\fitness_for_splunk\\bin\\microsoft_secret.json"
            
            # Create flow object
            flow = flow_from_clientsecrets(microsoft_secret_filepath, scope='mshealth.ReadProfile mshealth.ReadActivityHistory mshealth.ReadDevices mshealth.ReadActivityLocation offline_access', redirect_uri='https://localhost:8089/services/fitness_for_splunk/microsoft_callback')
            
            # Set access_type to 'offline' to retrieve Refresh token
            flow.params['access_type'] = 'offline'
            
            #auth_uri = flow.step1_get_authorize_url()
            #redirect user to auth_uri
            
            """
            Exchange authorization code for token
            """
            credentials = flow.step2_exchange(authCode)
            cred_json = json.loads(credentials.to_json())
            token = cred_json['token_response']
            access_token = token['access_token'].encode("utf-8")
            user_id = token['user_id'].encode("utf-8")
            
            """
            Use token to retrieve Full Name
            """
            # Get User Profile JSON
            userProfile_url = 'https://api.microsofthealth.net/v1/me/Profile'
            authHeader_value = ('bearer ' + access_token)
            headers = {'Authorization': authHeader_value}
            body = ""
            http = httplib2.Http()
            resp, content = http.request(userProfile_url, 'GET', headers=headers, body=body)
            
            profile_json = json.loads(content)
            username = user_id + "////" + profile_json['firstName']
            realm = "microsoft"
            
            """
            Store token in Splunk password store
            """
            # TODO
            
            """
            Write Response
            """
            self.response.setStatus(200)
            #self.response.setHeader('content-type', 'text/html')
            self.response.setHeader('content-type', 'application/json')
            self.response.write(username)
            
    # listen to all verbs
    handle_POST = handle_DELETE = handle_PUT = handle_VIEW = handle_GET