var splunkjs            = require('splunk-sdk');
var ModularInputs       = splunkjs.ModularInputs;
var Logger              = ModularInputs.Logger;
var Event               = ModularInputs.Event;
var Scheme              = ModularInputs.Scheme;
var Argument            = ModularInputs.Argument;
var utils               = ModularInputs.utils;
var fs                  = require("fs");
var path                = require("path");
var AdminUserName       = 'jsprowls';
var AdminPassword       = '38der#IyLL0%n%NI@00n#lc3f';
var service             = new splunkjs.Service({username: AdminUserName, password: AdminPassword});
var unirest             = require("unirest");
var reqauth             = unirest("POST", "https://login.live.com/oauth20_token.srf");
var requser             = unirest("GET", "https://api.microsofthealth.net/v1/me/PROFILE");
var reqdata             = unirest("GET", "https://api.microsofthealth.net/v1/me/Summaries/Hourly");
var _storeClientID      = 'N/A';
var _storeClientKey     = 'N/A';
var _refreshToken       = 'N/A';
var _authToken          = 'N/A';
var _username           = 'N/A';
var _fullname           = 'N/A';
var _lastSyncDate       = 'N/A';
var _startDate          = 'N/A';
var _endDate            = 'N/A';
var _name               = 'JAY';



var testdate            = "2016/8/5";

var RunTime             = new Date();
var processRunTime      = RunTime.toString();
var mynote              = 'Test';
var errorFound          = false;
//var maxAPIpullcount = 40;
//var APIpullcount    = 0;
var APIIterate  = 0;


service.login(function(err, success) {
    if (err) {
        throw err;
    }

    console.log("Login was successful: " + success);
    //--LoggedIn


    //==========================================================================
    //--Input Scheme
    //==========================================================================
    exports.getScheme = function () {
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
    exports.validateInput = function (definition, done) {
        var token_json = definition.parameters.token_json;
        var startdate = definition.parameters.date;
        done();
    };
    //------------------------------------------
    //--PreCalculate End Date
    //------------------------------------------
    function GetBeginOfHourDate() {
        var d = new Date();
        var d2 = new Date();
        d2.setMinutes(d.getMinutes() - d.getMinutes());
        d2.setSeconds(d.getSeconds() - d.getSeconds());
        d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
        var tempdate = d2.toISOString();
        return tempdate;
    };
    //------------------------------------------
    //--Convert Origional Start Date
    //------------------------------------------
    function GetISODate(date) {
        var d = new Date();
        var d2 = new Date(date);
        d2.setMinutes(d.getMinutes() - d.getMinutes());
        d2.setSeconds(d.getSeconds() - d.getSeconds());
        d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
        var tempdate = d2.toISOString();
        return tempdate;
    };
    //------------------------------------------
    //--Calculate Date Difference by hour
    //------------------------------------------
    function GetDateDifference(date1,date2) {
        var d = new Date();
        var d2 = new Date(date1);
        d2.setMinutes(d.getMinutes() - d.getMinutes());
        d2.setSeconds(d.getSeconds() - d.getSeconds());
        d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
        var s = new Date();
        var s2 = new Date(date2);
        s2.setMinutes(s.getMinutes() - s.getMinutes());
        s2.setSeconds(s.getSeconds() - s.getSeconds());
        s2.setMilliseconds(s.getMilliseconds() - s.getMilliseconds());
        var diff = Math.abs(s2.getTime() - d2.getTime()) / 3600000;
        return diff;
    };
    //------------------------------------------
    //--Convert Origional Start Date
    //------------------------------------------
    function GetDateAddedHours(date,hours) {
        var d = new Date();
        var d2 = new Date(date);
        d2.setMinutes(d.getMinutes() - d.getMinutes());
        d2.setSeconds(d.getSeconds() - d.getSeconds());
        d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
        d2.setHours(d.getHours()+40);
        var tempdate = d2.toISOString();
        return tempdate;
    };
  
  
//exports.streamEvents = function (name, singleInput, eventWriter, done) {
        //_username = singleInput.username;
        //_fullname = singleInput.fullname;
        //_name = singleInput.name;
        //_refreshtoken = singleInput.token_json;
        //var checkpointDir       = this._inputDefinition.metadata["checkpoint_dir"];
       // var checkpointFilePath  = path.join(checkpointDir, "MSHEalthCPs_" + _name + ".txt");

// Create a Service instance and log in 
var service2 = new splunkjs.Service({
  username:AdminUserName,
  password:AdminPassword,
  scheme:"https",
  host:"localhost",
  port:"8089",
  version:"5.0"
});

// Print installed apps to the console to verify login



 var endpoint = new splunkjs.Service.Endpoint(service2, "servicesNS/nobody/fitness_for_splunk/storage/collections/data/microsoft_tokens");
 endpoint.get("results", {offset: 1}, function(res) {
      console.log(res);

});

 
/*service2.Entity().fetch(function(err, apps) {
  if (err) {
    console.log("Error retrieving apps: ", err);
    return;
  }

  console.log("Applications:");

  var appsList = apps.list();
  for(var i = 0; i < appsList.length; i++) {
    var app = appsList[i];
    console.log("  App " + i + ": " + app.name);
  }
});*/






//1-Start-----------------------------------------------------------------------------------------------
//1--Gather CLientID and Client Secret Key from storage passwords
//1-----------------------------------------------------------------------------------------------------
service.storagePasswords().fetch(
    function(err, storagePasswords) {
     //   console.log(storagePasswords.list());
        if (err) 
            { /* handle error */ }
        else {
        // Storage password was created successfully
        if (storagePasswords.list().length >0 ){
            for(var i = 0; i <  storagePasswords.list().length; i++) {
          //  console.log( storagePasswords.list()[i]);
           if (storagePasswords.list()[i]._properties.realm = 'microsoft') 
           {
            console.log( "realm:"+storagePasswords.list()[i]._properties.realm);

            
            _storeClientID = storagePasswords.list()[i]._properties.username
            console.log( "_storeClientID:"+_storeClientID);
           _storeClientKey = storagePasswords.list()[i]._properties.clear_password
            console.log("_storeClientKey:"+_storeClientKey)
            //2-Start-----------------------------------------------------------------------------------------------
            //2--Gather Token From KV Store
            //2-----------------------------------------------------------------------------------------------------
                _refreshToken = 'MCWEvRoUDxGIvJdlP4R3g*TfVClLRZvEBt2s67m3OiHugUJUDMuqRGzYE23KqLGCo90OSlrTYBuwKJjFgp9OgFP6zYmPzRRusBDmENNe3o4Nv1wmp97z6KDSrrAyr!zL*hmFwVudYTFjA583oaZ1krH5Hn3SKeoEqRNyvsGeIKwTHT5Nbn2ZPAAXbd*NPQTuN8EC!a6tQzdjLIGlF2eWsJyqpeeWmb0BDflBZ34Pg3ahO*DM6rSHrIdpyTj7UKGxAfWpuDduXAg7YFp3*rPE3NazHPLmpUUw4dGNsCwkZtYbtRS8PNEIIiv3GrX8icT2D9GXQvT0JtHNUVfRne6wUXJCaMwV5CWkJOdPRlAS0PSwcbe4*HfG*uihLb*8HoBtBNSrbUWvZnifyWCt*mardbLTDjAhL6ViUCqM7AtOWP1*4'
                console.log("_refreshToken:"+_refreshToken);
                    //3-Start-----------------------------------------------------------------------------------------------
                    //3--Refresh Token for new authToken
                    //3-----------------------------------------------------------------------------------------------------
                    //3--Auth Setup to get new token
                    reqauth.query({"redirect_uri": "http://127.0.0.1:3000"});
                    reqauth.headers({"content-type": "application/x-www-form-urlencoded"});
                    reqauth.timeout(60000);
                    reqauth.form({
                        "client_id": _storeClientID,
                        "grant_type": "refresh_token",
                        "refresh_token": _refreshToken,
                        "client_secret": _storeClientKey
                    });
                     reqauth.end(function (resauth) {
                         _authToken = resauth.body.access_token;
                         console.log("_authToken:"+_authToken);
                        //4-Start-----------------------------------------------------------------------------------------------
                        //4--Get Profile Data / Checkpoint
                        //4-----------------------------------------------------------------------------------------------------
                        requser.headers({"authorization": "Bearer " + _authToken});
                        requser.timeout(30000);
                        requser.end(function (resUser) {
                        if (resUser.error) {
                                working = false;
                                throw new Error(resUser.error);
                            }
                            var resultprof = JSON.stringify(resUser.body);
                            TempSyncDate = resUser.body.lastUpdateTime;
                            LastAccountSync=TempSyncDate;
                            var curEventprof = new Event({
                                    stanza: _name,
                                    data: '{"username":' + '"' + _username + '"' +
                                    ',"fullname":' + '"' + _fullname + '"' +
                                    ',"notes":' + '{' + mynote + '}' +
                                    ',"processRunTime":' + '"' + processRunTime + '"' +
                                    ',"dataType":' + '"' + 'User-Profile' + '"' +
                                    ',"data":' + resultprof + '}'
                                });
                                try {
                                    console.log(curEventprof);
                                    _lastSyncDate = resUser.body.lastUpdateTime;
                                    console.log(_lastSyncDate);
                                    // eventWriter.writeEvent(curEventprof);
                                    
                                    //4--Calculate Date
                                    _endDate = GetBeginOfHourDate();
                                    console.log("_endDate:"+_endDate);
                                    _startDate = GetISODate(testdate);
                                    console.log("_startDate:"+_startDate);
                                    console.log("------");
                                    APIpullcount = GetDateDifference(testdate,RunTime)
                                    console.log(APIpullcount);
                                    var tempEndDate = ''
                                    var tempStartDate = ''
                                    var runOnce = false;
                                   

                                   
                                    //5-Start-----------------------------------------------------------------------------------------------
                                    //5--Get Summary Data
                                    //5-----------------------------------------------------------------------------------------------------
                                    var needtorun = true;
                                    ////for (var x = 0; x < APIpullcount && !errorFound; x++) {
                                       /* while (needtorun==true) {//Testing more than 48

                                            if (runOnce == true) {
                                                if (APIpullcount - APIIterate - maxAPIpullcount >= 0) {
                                                    //Updates Start and End
                                                    console.log("AAA");
                                                    _startDate = GetDateAddedHours(_startDate, maxAPIpullcount);
                                                    _endDate = GetDateAddedHours(_startDate, maxAPIpullcount);
                                                    APIIterate = APIIterate + GetDateDifference(_startDate, _endDate)

                                                    console.log("APIIterate:"+APIIterate);
                                                }
                                                else {
                                                    console.log("ZZZ");
                                                    _startDate = GetDateAddedHours(_startDate, (APIpullcount - APIIterate));
                                                    _endDate = GetBeginOfHourDate();
                                                    APIIterate = APIIterate + GetDateDifference(_startDate, _endDate)
                                                    needtorun = false;
                                                    console.log("APIIterate:"+APIIterate);
                                                }
                                            }
                                            if (runOnce == false) {
                                                if (APIpullcount - APIIterate - maxAPIpullcount >= 0) {
                                                    console.log("BBB");
                                                    //Update End
                                                    _endDate = GetDateAddedHours(_startDate, maxAPIpullcount);
                                                    APIIterate = APIIterate + GetDateDifference(_startDate, _endDate)
                                                    console.log("APIIterate:"+APIIterate);
                                                }
                                                else {
                                                    //Leave Data alone
                                                    APIIterate = APIIterate + GetDateDifference(_startDate, _endDate)
                                                    console.log("YYY");
                                                    needtorun = false;
                                                    console.log("APIIterate:"+APIIterate);
                                                }
                                            }
                                    console.log("_endDate2:"+_endDate);
                                   console.log("_startDate2:"+_startDate);*/
                                    
                                        

                                        reqdata.query({ "startTime": _startDate, "endTime": _endDate });
                                        reqdata.headers({ "authorization": "Bearer " + _authToken });
                                        reqdata.end(function (resData) {//!!Step 4 Get Summaries Process!!
                                            if (resData.error) {
                                                throw new Error(resData.error);
                                            }
                                            //  console.log(res.body);
                                            //  var result = JSON.stringify(res);
                                            var count = resData.body.summaries.length;
                                            var result = JSON.stringify(resData.body.summaries[0]);
                                            //--Loop through the hourly events if there are more than one
                                            for (var i = 0; i < count && !errorFound; i++) {
                                                result = JSON.stringify(resData.body.summaries[i]);

                                                var curEvent = new Event({
                                                    stanza: _name,
                                                    data: '{"username":' + '"' + _username + '"' +
                                                    ',"fullname":' + '"' + _fullname + '"' +
                                                    ',"notes":' + '{' + mynote + '}' +
                                                    ',"processRunTime":' + '"' + processRunTime + '"' +
                                                    ',"dataType":' + '"' + 'User-Summary' + '"' +
                                                    ',"data":' + result + '}'
                                                });
                                                try {
                                                    APIIterate++;
                                                    console.log(APIIterate);
                                                  ///  console.log(curEvent);
                                                    //eventWriter.writeEvent(curEvent);
                                                    /* 
                                                     var startObj = '  "lastcheckpoint": { "name": "' + _username + '"' +
                                                         ',"lastprocessrun": "' + processRunTime + '"' +
                                                         ',"startdate": "' + _StartTime + '"' +
                                                         ',"enddate": "' + _EndTime + '"' +
                                                         ',"lastaccountsync": "' + LastAccountSync + '"' +
                                                         '}';
     
                                                     fs.truncateSync(checkpointFilePath, 0);
     
                                                     fs.appendFileSync(checkpointFilePath, startObj);//_EndTime);
                                                     */
                                                }
                                                catch (e) {
                                                    errorFound = true;
                                                    Logger.error(name, e.message);
                                                    // done(e);
                                                    //--We had an error, die
                                                    return;
                                                }
                                            }
                                        });
                                        //runOnce = true;//Testing more than 48
                                        //
                                       //}//Testing more than 48
                                    ////}
                                    //5-END-------------------------------------------------------------------------------------------------
                                    }
                                catch (e) {
                                    errorFound = true;
                                    Logger.error(name, e.message);
                                    //done(e);
                                    //--We had an error, die
                                    return;
                                }
                            
                        });

                        //4-END-------------------------------------------------------------------------------------------------
                    });

                    //3-END-------------------------------------------------------------------------------------------------

            //2-END-------------------------------------------------------------------------------------------------
            }
        }
        
            
            
        }
        
        }
    });
    //1-END-------------------------------------------------------------------------------------------------


 
});//EndLogin