from bottle import route, run, request # Python server library
import sys
import httplib2, urllib
import base64

# Hello World route example
@route('/hello')
def hello():
    return "Hello World!"

# Fitbit callback route
@route('/auth/fitbit/callback')
def fitbit_callback():
    
    # Edit these variables to suit you
    clientID = '227MVJ'
    clientSecret = 'df8009bd0ddcb975f9a812e3587e54dd'
    encoded = base64.b64encode( (clientID + ':' + clientSecret) )
    callback_url = 'https://localhost:8089/services/fitness_for_splunk/fitbit_callback'
    authCode = '' # Need to fill in auth cod
    
    # Request for a token
    url = 'https://api.fitbit.com/oauth2/token'
    authHeader_value = ('Basic ' + encoded)
    headers = {'Authorization': authHeader_value, 'Content-Type': 'application/x-www-form-urlencoded'}
    data = {'clientId': clientID, 'grant_type': 'authorization_code', 'redirect_uri': callback_url, 'code': authCode}
    body = urllib.urlencode(data)
    
    http = httplib2.Http()
    resp, cont = http.request(url, 'POST', headers=headers, body=body)
    
    # Print response content (token) to screen
    return cont

run(host='localhost', port=3000, debug=True)