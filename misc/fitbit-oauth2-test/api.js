var FitbitClient = require('fitbit-client-oauth2');
var fs = require('fs');
//var args = process.argv.slice(2);

// Load JSON from file
var file_json = JSON.parse(fs.readFileSync('./access_token.json', 'utf8'));
var token = file_json.token;
var access_token = token.access_token;
var refresh_token = token.refresh_token;
var scope = token.scope;
var expires_in = token.expires_in;
var expires_at = token.expires_at;

console.log(access_token);
console.log(refresh_token);
console.log(scope);
console.log(expires_in);
console.log(expires_at);

/* Fitbit Stuff */
var clientId = '227MVJ';
var clientSecret = 'df8009bd0ddcb975f9a812e3587e54dd';
var client = new FitbitClient(clientId, clientSecret);
var redirect_uri = 'http://localhost:3000/auth/fitbit/callback';

var client = new FitbitClient(clientId, clientSecret);

//var options = { /* TIME_SERIES_OPTIONS */ };
var options = { date: "2016-06-01" };

function getDailyActivitySummary(token, options)
{
    client.getDailyActivitySummary(token, options)
        .then(function(res) {
            var result = JSON.stringify(res);
            fs.writeFile('./success.json', result, function (err,data) {
                if (err) {
                    return console.log(err);
                }
                //console.log(result);
            });
            console.log('Fitbit data retrieved for', options.date);
            rl.prompt();
        }).catch(function(err) {
            var error = JSON.stringify(err);
            fs.writeFile('./error.json', error, function (err,data) {
                if (err) {
                    return console.log(err);
                }
                //console.log(err);
            });
            console.log('Error retrieving Fitbit data.', error);
            rl.prompt();
        });
}

const readline = require('readline');
const rl = readline.createInterface(process.stdin, process.stdout);

rl.setPrompt('Enter a date (yyyy-MM-dd): ');
rl.prompt();

rl.on('line', function(line) {
    line.trim();
    options.date = line.trim();
    console.log(options);
    getDailyActivitySummary(token, options);
}).on('close', function() {
    console.log('Stopping..');
    process.exit(0);
});