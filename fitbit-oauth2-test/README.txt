DEPENDENCIES:
npm install fitbit-client-oauth2
npm install express
npm install body-parser

USAGE:
token.js:
-Starts a server with an endpoint for authentication with Fitbit API.
-Fitbit API redirects to callback (specific to Andy H's Fitbit App credentials)
-Access Token is written in token.json file.

api.js:
-Procedurally retrieves Fitbit Data for specified date.
-Writes retrieved data to success.json file.