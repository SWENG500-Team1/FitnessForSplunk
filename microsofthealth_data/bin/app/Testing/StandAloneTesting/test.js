
var unirest = require("unirest");

var req = unirest("GET", "https://login.live.com/oauth20_authorize.srf");

req.query({
    "client_id": "e2a30a63-c396-457a-abbf-1409678a5be4",
    "scope": "mshealth.ReadProfile mshealth.ReadDevices mshealth.ReadActivityHistory mshealth.ReadActivityLocation",
    "response_type": "code",
    "redirect_uri": "http://localhost:1337/"
});

req.headers({
    "postman-token": "9a0a48c0-11e5-e8d4-bca9-4cc544a1f304",
    "cache-control": "no-cache",
    "accept": "application/json"
});

req.send("{\r\n    \"question\": \"Favourite programming language?\",\r\n    \"choices\": [\r\n        \"Swift\",\r\n        \"Python\",\r\n        \"Objective-C\",\r\n        \"Ruby\"\r\n    ]\r\n}");

req.end(function (res) {
    if (res.error) throw new Error(res.error);

    console.log(res.body);

});








var req = unirest("POST", "https://login.live.com/oauth20_token.srf");

req.query({
  "client_id": "e2a30a63-c396-457a-abbf-1409678a5be4",
  "redirect_uri": "http://localhost:1337/",
  "client_secret": "aZdpuFCxAde9QAmBX3U1Sa5",
  "code": "Mdf65e226-dafa-60c4-89ad-ee214670f2aa",
  "grant_type": "authorization_code"
});

req.headers({
  "postman-token": "ad003b88-6644-b40b-482a-fc632f1d01ba",
  "cache-control": "no-cache"
});


req.end(function (res) {
  if (res.error) throw new Error(res.error);

  console.log(res.body);
});
