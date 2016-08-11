// Copyright 2014 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

(function() {
    var splunkjs        = require("splunk-sdk");
    var ModularInputs   = splunkjs.ModularInputs;
    var Logger          = ModularInputs.Logger;
    var Event           = ModularInputs.Event;
    var Scheme          = ModularInputs.Scheme;
    var Argument        = ModularInputs.Argument;
    var utils           = ModularInputs.utils;

    var splunkjs        = require("splunk-sdk");
    var ModularInputs   = splunkjs.ModularInputs;
    var Logger          = ModularInputs.Logger;
    var Event           = ModularInputs.Event;
    var Scheme          = ModularInputs.Scheme;
    var Argument        = ModularInputs.Argument;
    var utils           = ModularInputs.utils;

    var _clientID     = 'e52600ac-ada7-4d65-93b1-b5013037d4e1'//'e2a30a63-c396-457a-abbf-1409678a5be4'
    var _clientSecret = 'XVXygsDydoJTXQR4KTRQLxd'//'aZdpuFCxAde9QAmBX3U1Sa5'
    var _accessToken  = 'EwBQA/F0BAAUvSQiG6C8oi/OrqaVrv8+s9GGlDkAAU4cE6dpIdjcca8rSG3z8zOXTGXvnK5TwImEmr12eTzLF/1r8beZGH9/xE8UZGR9p8z3bt2EHnJEppEvphI+Q+6+x2o7UQOnl0TjJrT9dpMd8594sryY7BfUlMAojFEA8W57LKczKB7xcNieoVSu5xsRdFMBO8m46LpbQL8tBSzWuAp8Tbdnk9Z9eUlc+6keJBZEBoxD4pySBVLkZzzuveEVJq7JsGauBfle6F5mPJ1ulQfR5F0LgkaBr8+jxRP3KzCXDLkyaN+DIVnx9GNEw8HzLgbVabRY8eZSx0AjNLNYgD7lO3L9ynyzpyzjtMuB/MEDLMReCF4F/De8giEgQMEDZgAACJ1d/837RbuvIAIwqcaMr4W3mrfMs+2fQge6RxU3Brx4IFhjUXZay8Ljfn5IwDki4UzX3l74sn3FSR2vD2Ub+AERVp80wL+ziT/8GAMj3t+pAVYhtxpf1OOgqJYyBeHLL2Zb1clRwQAtqQ5DKjaKKbD1mNvMylV59xfk09FHCGu3zuvF0R+wnFj2Fsh/br1idkJs5SxNtrSBAT4gbnXEUlfiepYnaE6OG86XOP3Zxw1pQnffYY2oCEX2GjJK61M+sngugt2OBq3BP6X+6CLf6zFFETFjY4oJMN2EfIoInyAElZV7EU8SaYTRDz7wGMLxfcbzgzDPhxBMOAcV24QV6cMgSEs69QQaE1FdfrCHJ8+t9fGmi54NsssUofmlPs0o8/b1HesYdCyZHWCuUmPxbd9pfUW+zyPmEAcT/jZGzYDJFIb5UBax/p5VHHgtffm/Ad0IUeVFrIL0WGSrvXJ7lMYZ/uv5WVg+61dTVI/k6B9OpDkM9LtPfYy0hpQlzfDzDSi/wKg3AS2EUcTcAacaTZpgISRFLZxKv2nzK4WPh+MwuouJrzNrWThgPGUT/IEfgdq45WpF8qrTUP/DDrwgZbXOyzZA1MrHjDIqvP2PIgnfYR5k07lpXo3uYffLWvGqIbOW90D2g8qtZi6vOw1qOqEKoWvEK7tU8GcNB5WewEhbeYg2zl+bOkMY21SK60YiK6m+rLrkdsOrFc57FCGLalFa3uvSey56yTRBNAI='
    var _refreshtoken = 'MCWEvRoUDxGIvJdlP4R3g*TfVClLRZvEBt2s67m3OiHugUJUDMuqRGzYE23KqLGCo90OSlrTYBuwKJjFgp9OgFP6zYmPzRRusBDmENNe3o4Nv1wmp97z6KDSrrAyr!zL*hmFwVudYTFjA583oaZ1krH5Hn3SKeoEqRNyvsGeIKwTHT5Nbn2ZPAAXbd*NPQTuN8EC!a6tQzdjLIGlF2eWsJyqpeeWmb0BDflBZ34Pg3ahO*DM6rSHrIdpyTj7UKGxAfWpuDduXAg7YFp3*rPE3NazHPLmpUUw4dGNsCwkZtYbtRS8PNEIIiv3GrX8icT2D9GXQvT0JtHNUVfRne6wUXJCaMwV5CWkJOdPRlAS0PSwcbe4*HfG*uihLb*8HoBtBNSrbUWvZnifyWCt*mardbLTDjAhL6ViUCqM7AtOWP1*4'

    var _newAccessToken = 'aaaa'
    var teststartDate = '';
	var endDate = '';
    var fs = require('fs');

//==========================================================================
//--Input Scheme
//==========================================================================
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
                description: "Enter last date searched on or date to start searching on.",
                requiredOnCreate: false,
                requiredOnEdit: false
            }),
            new Argument({
                name: "token_json",
                dataType: Argument.dataTypeString,
                description: "Please enter your Refresh Auth Token.",
                requiredOnCreate: true,
                requiredOnEdit: false
            })
        ];
        return scheme;
    };
//==========================================================================
//--Validate the Input
//==========================================================================
    exports.validateInput = function(definition, done) {
        var token_json = definition.parameters.token_json;
        var startdate = definition.parameters.date;
          done();
    };
//==========================================================================
//--Code Modules for Application
//==========================================================================	
     //------------------------------------------
     //--PreCalculate End Date
     //------------------------------------------
     var d = new Date();
     var d2 = new Date();
     d2.setMinutes(d.getMinutes() - d.getMinutes());
     d2.setSeconds(d.getSeconds() - d.getSeconds());
     d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
       endDate = d2.toISOString();	
	 //------------------------------------------
     //--PreCalculate Test Start Date
     //------------------------------------------
     var s = new Date();
     var s2 = new Date();
     s2.setHours(s.getHours() - 1);
	 s2.setMinutes(s.getMinutes() - s.getMinutes());
     s2.setSeconds(s.getSeconds() - s.getSeconds());
     s2.setMilliseconds(s.getMilliseconds() - s.getMilliseconds());
       teststartDate = s2.toISOString();	  
	   
	   
	   
	
//==========================================================================
//--Process To: Refresh Token / Get Data / Send to Splunk
//==========================================================================
    exports.streamEvents = function(name, singleInput, eventWriter, done) {
//------------------------------------------------------------------------------------------------------------------------------
	var username    = singleInput.username;
        var fullname    = singleInput.fullname;
         _refreshtoken  = singleInput.token_json;
        var start_date  = singleInput.date;
        var _StartTime   =teststartDate//"2016-07-30T00:00:00.000Z"
        var _EndTime    = endDate//"2016-07-30T15:00:00.000Z"

        var errorFound = false;

	var unirest = require("unirest");
	var reqauth = unirest("POST", "https://login.live.com/oauth20_token.srf");
	var reqdata = unirest("GET", "https://api.microsofthealth.net/v1/me/Summaries/Hourly");

reqauth.query({
  "redirect_uri": "http://127.0.0.1:3000"
});
//------------------------------------------
//--Set the auth format
//------------------------------------------
reqauth.headers({
  "content-type": "application/x-www-form-urlencoded"
});
//------------------------------------------
//--Auth Setup to get new token
//------------------------------------------
reqauth.form({
  "client_id": _clientID,
  "grant_type": "refresh_token",
  "refresh_token": _refreshtoken,
  "client_secret": _clientSecret
});

//------------------------------------------
//--Data Set the time
//------------------------------------------
reqdata.query({
  "startTime": _StartTime,
  "endTime": _EndTime
});

reqauth.end(function (resauth) {
var _newAccessToken = resauth.body.access_token;
//------------------------------------------
//--Data Set the Auth Token
//------------------------------------------
reqdata.headers({
  "authorization": "Bearer "+_newAccessToken//token_json//EwBIA/F0BAAUvSQiG6C8oi/OrqaVrv8+s9GGlDkAAcfNbsHy71g6GI7iGuP0I3xCpSRpJna1L2hX82hTiEbYhvucN2deDn12LF0CydfZgm5Bl9g6a3DWX2nWLgXNP0NGjOX84ohY5T3mWRNvgQTMJIkifDMuEBGl+gx5MpbUMi3tbbH0mGDfzhvPioh8af5CjPkwrlzp5M48oXGPLomVsZ1BrgAJteQYfGaNxK550U1uza+4TgDyWTyw/E0oxEy82AYzshYND3wuH6TAzSvFJBRVaYBlw+BnvoD6kShiF+4OcUke/o/zsbF4mZoHRyiFySmU0ST55N00aLoF+xnwaIaY/H7DIeIW5G1O/E3CuYIllUKcc5JTjZL9qIFfL2gDZgAACNx6628coiTaGAINqYtN/ztvSX58hXSJvYAcEHZsjyRhFf2RO13PV7XUu6tgyI4MWCWO6jI0W2RueLMytlOO6b2QaTIfShXcNvnuvnirWBhr9I11H2KFV+iLurzZE3w1HKKiieVkbquaYTyYnsBoIgYV0saik34leaRSfkHOGGZjxn14T3S7qIBCxEEyzJBnrBPJjYK/JRl5JhBtR8JXMCUchkxpnE9dYB4wUCXlGFz/QlfkEd+n9Z0QqREUkUwfuxbj5ghmtnNO03DgtSnTrJFETQNTFldp6E0lavB8/GZWOPkzpsk9d8vUHiqWlK4lR/BSZ5S5xNSHxbjlBJID2xCc+dFZ6DXQ7gIu1CIsbU+xFDdyclMRPNpsPAkS1kM6Pi/hfbNN1dRgMRRk0dIDV17vxDC/R6JmrMjIGdlcUopxiWMXCJ87sLoa2I2qV8zVhfp0y+4vwNrwJ09tLJlrrWtAL57gtFkdAcr0VCRFwEBGTfqv02YvbpTsjZia7+dGm73Dc08xuQ3MBOa0v/6oPidJiNcLXLFSLIPXxH7Vl4+oout7MXINRCu4KyGOjl0kiPG3AIx0xNYqYu4dnGLRc2G9g33wiQxPBjomKPipMOZli+B6aQyLyp1d+lH+iB5jSDDO8tQRmKuujcDkybEieyTfmx+FmN/RdlG26OA0geYLAcEkjPvkYA8V3bVICkEBhKNcd7aQwfJl0jX+hGPA5qKSBjYC"
});
//------------------------------------------
//--Get Data
//------------------------------------------
////////////////////////////////////////////////////////////////////////////
reqdata.end(function (resData) {
if (resData.error) throw new Error(resData.error);
//  console.log(res.body);
//  var result = JSON.stringify(res);
var count = resData.body.summaries.length;
var result = JSON.stringify(resData.body.summaries[0]);
//--Loop through the hourly events if there are more than one
  for (var i = 0; i < count && !errorFound; i++) {
    result = JSON.stringify(resData.body.summaries[i]);
      var curEvent = new Event({
          stanza: name,
          data: '{"username":' + '"'+username +'"'+ ',"fullname":' + '"'+fullname +'"'+ ',"data":' + result + '}'
      });
      try {
          eventWriter.writeEvent(curEvent);
      }
      catch (e) {
          errorFound = true;
          Logger.error(name, e.message);
          done(e);
          //--We had an error, die
          return;
      }
  }
});
});


//------------------------------------------------------------------------------------------------------------------------------
        //--We're done
        done();
    };
//==========================================================================	

    ModularInputs.execute(exports, module);
})();