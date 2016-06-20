var express = require('express');
var bodyParser = require('body-parser');
var FitbitClient = require('fitbit-client-oauth2');
var fs = require('fs');

/* Express server */
var app = express();

/* Fitbit Stuff */
var clientId = '';
var clientSecret = '';
var client = new FitbitClient(clientId, clientSecret);
var redirect_uri = 'http://localhost:3000/auth/fitbit/callback';

app.use(bodyParser());

scope = [ 'activity', 'heartrate', 'location', 'nutrition', 'profile', 'settings', 'sleep', 'social', 'weight' ];
console.log(scope);

/* Navigate to this endpoint to start authentication process */
app.get('/auth/fitbit', function(req, res) {
    
    var auth_url = client.getAuthorizationUrl('http://localhost:3000/auth/fitbit/callback', scope);
    res.redirect(auth_url);
});

/* Fitbit API redirects to this callback with authorization code (configured in Fitbit App settings) */
app.get('/auth/fitbit/callback', function(req, res, next) {

    client.getToken(req.query.code, redirect_uri)
        .then(function(token) {

            // ... save your token on session or db ...
            // Write token to file
            var json_string = JSON.stringify(token);
            fs.writeFile('./access_token.json', json_string, function (err,data) {
                if (err) {
                    return console.log(err);
                }
                console.log(json_string);
            });
            
            // then redirect
            //res.redirect(302, '/user');
            console.log(token);
            res.send(token);
        })
        .catch(function(err) {
            // something went wrong.
            res.status(500).send(err);
            console.log(err);
        });
});

app.get('/', function (req, res) {
    
    res.send('Hello World!');
});

app.listen(3000, function() {
    
    console.log('Example app listening on port 3000!');
});