//Microsoft HEalth Data API ModularInputs for Splunk


(function() {
    var splunkjs        = require("splunk-sdk");
    var ModularInputs   = splunkjs.ModularInputs;
    var Logger          = ModularInputs.Logger;
    var Event           = ModularInputs.Event;
    var Scheme          = ModularInputs.Scheme;
    var Argument        = ModularInputs.Argument;
    var utils           = ModularInputs.utils;

    var _clientID     = 'e2a30a63-c396-457a-abbf-1409678a5be4'
    var _clientSecret = 'aZdpuFCxAde9QAmBX3U1Sa5'


    exports.getScheme = function() {
        var scheme = new Scheme("Microsoft Health Data");

        scheme.description = "Gets Micrsoft Health Data.";
        scheme.useExternalValidation = true;
        scheme.useSingleInstance = true;

        scheme.args = [
            new Argument({
                name: "fullname",
                dataType: Argument.dataTypeString,
                description: "Please enter your first and lastname.",
                requiredOnCreate: true,
                requiredOnEdit: false
            }),


            new Argument({
                name: "username",
                dataType: Argument.dataTypeString,
                description: "Please enter your Microsoft Health UserName.",
                requiredOnCreate: true,
                requiredOnEdit: false
            }),
            new Argument({
                name: "startdate",
                dataType: Argument.dataTypeString,
                description: "Enter last date serched on or date to start searching on.",
                requiredOnCreate: false,
                requiredOnEdit: false
            }),
            new Argument({
                name: "token_json",
                dataType: Argument.dataTypeString,
                description: "Please enter your Auth Token.",
                requiredOnCreate: true,
                requiredOnEdit: false
            })
        ];

        return scheme;
    };

    exports.validateInput = function(definition, done) {
        var token_json = definition.parameters.token_json;
        var startdate = definition.parameters.date;


          done();
    };

    exports.streamEvents = function(name, singleInput, eventWriter, done) {
        var username = singleInput.username;
        var token_json = JSON.parse(singleInput.token_json)
        var start_date = singleInput.date;

        var token = token_json.token;
        var options = { date: "2016-06-01" };
        options.date = start_date;

        var errorFound = false;
//-----------------------------------------------------------------
var unirest = require("unirest");

var req = unirest("GET", "https://api.microsofthealth.net/v1/me/Summaries/Hourly");

req.query({
  "startTime": "2016-07-15T00:00:00.000Z",
  "endTime": "2016-07-16T00:00:00.000Z"
});

req.headers({
  "authorization": "Bearer EwBIA/F0BAAUvSQiG6C8oi/OrqaVrv8+s9GGlDkAAcfNbsHy71g6GI7iGuP0I3xCpSRpJna1L2hX82hTiEbYhvucN2deDn12LF0CydfZgm5Bl9g6a3DWX2nWLgXNP0NGjOX84ohY5T3mWRNvgQTMJIkifDMuEBGl+gx5MpbUMi3tbbH0mGDfzhvPioh8af5CjPkwrlzp5M48oXGPLomVsZ1BrgAJteQYfGaNxK550U1uza+4TgDyWTyw/E0oxEy82AYzshYND3wuH6TAzSvFJBRVaYBlw+BnvoD6kShiF+4OcUke/o/zsbF4mZoHRyiFySmU0ST55N00aLoF+xnwaIaY/H7DIeIW5G1O/E3CuYIllUKcc5JTjZL9qIFfL2gDZgAACNx6628coiTaGAINqYtN/ztvSX58hXSJvYAcEHZsjyRhFf2RO13PV7XUu6tgyI4MWCWO6jI0W2RueLMytlOO6b2QaTIfShXcNvnuvnirWBhr9I11H2KFV+iLurzZE3w1HKKiieVkbquaYTyYnsBoIgYV0saik34leaRSfkHOGGZjxn14T3S7qIBCxEEyzJBnrBPJjYK/JRl5JhBtR8JXMCUchkxpnE9dYB4wUCXlGFz/QlfkEd+n9Z0QqREUkUwfuxbj5ghmtnNO03DgtSnTrJFETQNTFldp6E0lavB8/GZWOPkzpsk9d8vUHiqWlK4lR/BSZ5S5xNSHxbjlBJID2xCc+dFZ6DXQ7gIu1CIsbU+xFDdyclMRPNpsPAkS1kM6Pi/hfbNN1dRgMRRk0dIDV17vxDC/R6JmrMjIGdlcUopxiWMXCJ87sLoa2I2qV8zVhfp0y+4vwNrwJ09tLJlrrWtAL57gtFkdAcr0VCRFwEBGTfqv02YvbpTsjZia7+dGm73Dc08xuQ3MBOa0v/6oPidJiNcLXLFSLIPXxH7Vl4+oout7MXINRCu4KyGOjl0kiPG3AIx0xNYqYu4dnGLRc2G9g33wiQxPBjomKPipMOZli+B6aQyLyp1d+lH+iB5jSDDO8tQRmKuujcDkybEieyTfmx+FmN/RdlG26OA0geYLAcEkjPvkYA8V3bVICkEBhKNcd7aQwfJl0jX+hGPA5qKSBjYC"
});


req.end(function (res) {
  if (res.error) throw new Error(res.error);

//  console.log(res.body);
  //var result = JSON.stringify(res);
  var count = res.body.summaries.length;
var result = JSON.stringify(res.body.summaries(0));
  ////////////////////////////////////////////////////////////////////////////
  for (var i = 0; i < count && !errorFound; i++) {
    result = JSON.stringify(res.body.summaries(i));
      var curEvent = new Event({
          stanza: name,
        //  data: "number=" + getRandomFloat(min, max)
          data: "{\"username\":" + username + ",\"data\":" + result + "}"
          //\"date\":"+fit_date+",

      });

      try {
          eventWriter.writeEvent(curEvent);
      }
      catch (e) {
          errorFound = true;
          Logger.error(name, e.message);
          done(e);

          // We had an error, die
          return;
      }
  }
  ///////////////////////////////////////////////////////////////////////////


});

      //-----------------------------------------------------------------

        // We're done
        done();
    };

    ModularInputs.execute(exports, module);
})();
