// Set the client credentials and the OAuth2 server
var credentials = {
  clientID: 'e2a30a63-c396-457a-abbf-1409678a5be4',
  clientSecret: 'aZdpuFCxAde9QAmBX3U1Sa5',
  site: 'https://login.live.com/oauth20_token.srf'
};

// Initialize the OAuth2 Library 
var oauth2 = require('simple-oauth2')(credentials);

// Authorization oauth2 URI
var authorization_uri = oauth2.authCode.authorizeURL({
  redirect_uri: ' https://login.live.com/oauth20_authorize.srf',
  scope: 'mshealth.ReadProfile mshealth.ReadDevices mshealth.ReadActivityHistory mshealth.ReadActivityLocation',
  state: ''
});

// Redirect example using Express (see http://expressjs.com/api.html#res.redirect)
res.redirect(authorization_uri);

// Get the access token object (the authorization code is given from the previous step).
var token;
var tokenConfig = {
  code: '<code>',
  redirect_uri: 'http://localhost'
};

// Callbacks
// Save the access token
oauth2.authCode.getToken(tokenConfig, function saveToken(error, result) {
  if (error) { console.log('Access Token Error', error.message); }

  token = oauth2.accessToken.create(result);
});

// Promises
// Save the access token
oauth2.authCode.getToken(tokenConfig)
.then(function saveToken(result) {
  token = oauth2.accessToken.create(result);
})
.catch(function logError(error) {
  console.log('Access Token Error', error.message);
});
