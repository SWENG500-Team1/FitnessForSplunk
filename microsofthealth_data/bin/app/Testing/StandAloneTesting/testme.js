var unirest = require("unirest");
var splunkjs        = require("splunk-sdk");
var ModularInputs   = splunkjs.ModularInputs;
var Logger          = ModularInputs.Logger;
var Event           = ModularInputs.Event;
var Scheme          = ModularInputs.Scheme;
var Argument        = ModularInputs.Argument;
var utils           = ModularInputs.utils;

var _clientID     = 'e2a30a63-c396-457a-abbf-1409678a5be4'
var _clientSecret = 'aZdpuFCxAde9QAmBX3U1Sa5'
var _accessToken  = 'EwBQA/F0BAAUvSQiG6C8oi/OrqaVrv8+s9GGlDkAAU4cE6dpIdjcca8rSG3z8zOXTGXvnK5TwImEmr12eTzLF/1r8beZGH9/xE8UZGR9p8z3bt2EHnJEppEvphI+Q+6+x2o7UQOnl0TjJrT9dpMd8594sryY7BfUlMAojFEA8W57LKczKB7xcNieoVSu5xsRdFMBO8m46LpbQL8tBSzWuAp8Tbdnk9Z9eUlc+6keJBZEBoxD4pySBVLkZzzuveEVJq7JsGauBfle6F5mPJ1ulQfR5F0LgkaBr8+jxRP3KzCXDLkyaN+DIVnx9GNEw8HzLgbVabRY8eZSx0AjNLNYgD7lO3L9ynyzpyzjtMuB/MEDLMReCF4F/De8giEgQMEDZgAACJ1d/837RbuvIAIwqcaMr4W3mrfMs+2fQge6RxU3Brx4IFhjUXZay8Ljfn5IwDki4UzX3l74sn3FSR2vD2Ub+AERVp80wL+ziT/8GAMj3t+pAVYhtxpf1OOgqJYyBeHLL2Zb1clRwQAtqQ5DKjaKKbD1mNvMylV59xfk09FHCGu3zuvF0R+wnFj2Fsh/br1idkJs5SxNtrSBAT4gbnXEUlfiepYnaE6OG86XOP3Zxw1pQnffYY2oCEX2GjJK61M+sngugt2OBq3BP6X+6CLf6zFFETFjY4oJMN2EfIoInyAElZV7EU8SaYTRDz7wGMLxfcbzgzDPhxBMOAcV24QV6cMgSEs69QQaE1FdfrCHJ8+t9fGmi54NsssUofmlPs0o8/b1HesYdCyZHWCuUmPxbd9pfUW+zyPmEAcT/jZGzYDJFIb5UBax/p5VHHgtffm/Ad0IUeVFrIL0WGSrvXJ7lMYZ/uv5WVg+61dTVI/k6B9OpDkM9LtPfYy0hpQlzfDzDSi/wKg3AS2EUcTcAacaTZpgISRFLZxKv2nzK4WPh+MwuouJrzNrWThgPGUT/IEfgdq45WpF8qrTUP/DDrwgZbXOyzZA1MrHjDIqvP2PIgnfYR5k07lpXo3uYffLWvGqIbOW90D2g8qtZi6vOw1qOqEKoWvEK7tU8GcNB5WewEhbeYg2zl+bOkMY21SK60YiK6m+rLrkdsOrFc57FCGLalFa3uvSey56yTRBNAI='
var _refreshtoken = 'MCWEvRoUDxGIvJdlP4R3g*TfVClLRZvEBt2s67m3OiHugUJUDMuqRGzYE23KqLGCo90OSlrTYBuwKJjFgp9OgFP6zYmPzRRusBDmENNe3o4Nv1wmp97z6KDSrrAyr!zL*hmFwVudYTFjA583oaZ1krH5Hn3SKeoEqRNyvsGeIKwTHT5Nbn2ZPAAXbd*NPQTuN8EC!a6tQzdjLIGlF2eWsJyqpeeWmb0BDflBZ34Pg3ahO*DM6rSHrIdpyTj7UKGxAfWpuDduXAg7YFp3*rPE3NazHPLmpUUw4dGNsCwkZtYbtRS8PNEIIiv3GrX8icT2D9GXQvT0JtHNUVfRne6wUXJCaMwV5CWkJOdPRlAS0PSwcbe4*HfG*uihLb*8HoBtBNSrbUWvZnifyWCt*mardbLTDjAhL6ViUCqM7AtOWP1*4'

var savedDate = '2016-07-16T00:00:00.000Z';
var endDate = '2016-07-30T00:00:00.000Z';
var testDate = '2016-07-16T00:00:00.000Z';

 var _newAccessToken = 'aaaa'
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////




//----------------------------------------------------------
var unirest = require("unirest");

var reqauth = unirest("POST", "https://login.live.com/oauth20_token.srf");

reqauth.query({
  "redirect_uri": "http://127.0.0.1:3000"

});

reqauth.headers({
  "content-type": "application/x-www-form-urlencoded"
});

reqauth.form({
  "client_id": _clientID,
  "grant_type": "refresh_token",
  "refresh_token": _refreshtoken,
  "client_secret": _clientSecret
});

  reqauth.end(function (resauth) {
  //if (res.error) throw new Error(res.error);

  //console.log(res.body);

  var _newAccessToken = resauth.body.access_token;

  //console.log(_newAccessToken);


  ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  var req2 = unirest("GET", "https://api.microsofthealth.net/v1/me/Summaries/Hourly");
  var errorFound = false;
  //////////////////////////////////////
  var fs = require('fs');
  var os = require("os");
  var path = 'logs.txt';
  function processInput ( text )
  {
    fs.open(path, 'a', 666, function( e, id ) {
     fs.write( id, text + os.EOL, null, 'utf8', function(){
      fs.close(id, function(){
      // console.log(_newAccessToken);
      });
     });
    });
   }
  /////////////////////////////////////

    processInput("Process Started: "+new Date());
/////////////////////////////////////



function getDateVars(_name,_lastdaterun){
fs.stat('File.txt', function(err, stat) {

    if(err == null) {
        console.log('File exists');
        var fs2 = require('fs');
        var obj = JSON.parse(fs.readFileSync('File.txt', 'utf8'));
      //console.log(obj.dateparm.length);
        for(var j = 0; j < obj.dateparm.length && !errorFound; j++){
          if (obj.dateparm[j].name == _name){
            console.log(obj.dateparm[j]);
            testDate = obj.dateparm[j].lastdaterun;
          }
        }
    } else if(err.code == 'ENOENT') {
        // file does not exist
        var today = new Date(_lastdaterun);

        testDate = today.toISOString();
        var startObj ='{  "dateparm": [{"id": 0, "name": "'+_name+'", "lastdaterun": "'+testDate+'"} ]}'
        fs.writeFile('File.txt', startObj);

    } else {
        console.log('Some other error: ', err.code);
    }



});

};
function customISOstring(){

}
savedDate = testDate;
var today2 = new Date();
getDateVars("Jay","08/01/2016 01:00");
 endDate = today2.toISOString();


console.log('s:'+savedDate);
console.log('e:'+endDate);

/////////////////////////////////////

  req2.query({
    "startTime": savedDate,
    "endTime": endDate
  });

  req2.headers({
  //"authorization": "Bearer "+_accessToken});
    "authorization": "Bearer "+_newAccessToken});
    //  console.log(_newAccessToken);



  /*req.end(function (res) {
    if (res.error) throw new Error(res.error);

  //  console.log(res.body);
    var result = JSON.stringify(res);
  //var o = JSON.parse(result);
    console.log(res.body.summaries.length);
  });*/

  req2.end(function (res2) {
  //  if (res.error) throw new Error(res.error);

  //  console.log(res.body);
    //var result = JSON.stringify(res);
    var count = res2.body.summaries.length;
  var result = JSON.stringify(res2.body.summaries[0]);
  console.log(count);
    ////////////////////////////////////////////////////////////////////////////
    for (var i = 0; i < count && !errorFound; i++) {
      result = JSON.stringify(res2.body.summaries[i]);
        var curEvent = new Event({
            stanza: 'name',
          //  data: "number=" + getRandomFloat(min, max)
            data: "{\"username\":" + 'username' + ",\"data\":" + result + "}"
            //\"date\":"+fit_date+",

        });

        try {
          //  eventWriter.writeEvent(curEvent);
          console.log(curEvent);
            return;//singletest
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

});
