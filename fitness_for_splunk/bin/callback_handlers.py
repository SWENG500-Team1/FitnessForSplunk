#
# http request handlers
#

from splunk import auth, search
import splunk.rest 
import splunk.bundle as bundle
import httplib2, urllib, os, time
import base64
import json
import xml.etree.ElementTree as ET
import splunk.entity as entity

#from test_modules import test_module

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

"""
from oauth2client.client import flow_from_clientsecrets
from apiclient.discovery import build

# set our path to this particular application directory (which is suppose to be <appname>/bin)
app_dir = os.path.dirname(os.path.abspath(__file__))

# define the web content directory (needs to be <appname>/web directory)
web_dir = app_dir + "/web"

class fitbit_callback(splunk.rest.BaseRestHandler):

    """
    Fitbit OAuth2 Callback Endpoint
    """
    def handle_GET(self):
        """
        Parse authorization code from URL query string
        """
        # Check for existence of 'code' parameter
        #for key in self.request['query']:
        #    if key == 'code':
        #        authCode = self.request['query']['code']
        #    else:
        #        authCode = None
                
        #if authCode is None:
        #    self.response.setStatus(400)
        #    self.response.setHeader('content-type', 'text/html')
        #    self.response.write("Bad request: no authorization code")
        #else:
            # TODO: Pull Client ID and Secret from password store
        #    fitbit_secret_filepath = "C:\\Program Files\\Splunk\\etc\\apps\\fitness_for_splunk\\bin\\fitbit_secret.json"
        #    with open(fitbit_secret_filepath, 'r') as file:
        #        fitbit_config = json.load(file)
        #    clientId = fitbit_config['web']['client_id'].encode("utf-8")
        #    clientSecret = fitbit_config['web']['client_secret'].encode("utf-8")
        #    callback_url = fitbit_config['web']['redirect_uris'][0].encode("utf-8")
            
        http = httplib2.Http(disable_ssl_certificate_validation=True)
            
        """
        Exchange authorization code for token
        """
        # Get OAuth2 token
        #    token_url = fitbit_config['web']['token_uri'].encode("utf-8")
        #    encoded_id_secret = base64.b64encode( (clientId + ':' + clientSecret) )
        #    authHeader_value = ('Basic ' + encoded_id_secret)
        #    headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
        #    data = {'clientId': clientId, 'grant_type': 'authorization_code', 'redirect_uri': callback_url, 'code': authCode}
        #    body = urllib.urlencode(data)
            
        #    resp, content = http.request(token_url, 'POST', headers=headers, body=body)
            
        """
        Use token to retrieve Full Name
        """
        #    token_json = json.loads(content)
        #    user_id = token_json['user_id'].encode("utf-8")
        #    access_token = token_json['access_token'].encode("utf-8")
            
            # Get User Profile JSON
        #    userProfile_url = 'https://api.fitbit.com/1/user/' + user_id + '/profile.json'
        #    authHeader_value = ('Bearer ' + access_token)
        #    headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
        #    body = ""
        #    resp, content = http.request(userProfile_url, 'GET', headers=headers, body=body)
            
        #    profile_json = json.loads(content)
        #    full_name = profile_json['user']['fullName'].encode("utf-8")
        #    username = user_id + "////" + full_name.replace(' ', '.')
        #    realm = "fitbit"
            
        """
        Store token in Splunk password store
        """
        # TODO
        
        # Get a sessonKey to access password store
        auth_url = "https://localhost:8089/services/auth/login"
        headers = {'Content-Type': 'application/x-www-form-urlencoded'}
        data = {'username': 'admin', 'password': 'root'}
        body = urllib.urlencode(data)
        resp, content = http.request(auth_url, 'POST', headers=headers, body=body)

        # Load XML
        root = ET.fromstring(content)

        # Parse XML and store sessionKey string
        sessionKey = root.find('sessionKey').text

        # Dummy username and token string to store
        username = "FIEJWO8////Andy.Huynh"
        token = "token"

        """
        # Sample code to retrieve password from Store
        """
        myapp = 'fitness_for_splunk'

        try:
            # list all credentials
            entities = entity.getEntities(['admin', 'passwords'], namespace=myapp, owner='nobody', sessionKey=sessionKey) 
        except Exception, e:
            raise Exception("Could not get %s credentials from splunk. Error: %s" % (myapp, str(e)))

        # return first set of credentials
        #for i, c in entities.items(): 
            #return c['username'], c['clear_password']

        credentials = []

        for i, c in entities.items():
            # Append credentials to array
            credentials.append((c['username'], c['clear_password']))

        # Credentials array is blank when I visit the page
        
        
        
        self.response.setStatus(200)
        #self.response.setHeader('content-type', 'text/html')
        self.response.setHeader('content-type', 'application/json')
        self.response.write(str(', '.join(credentials)))
            
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