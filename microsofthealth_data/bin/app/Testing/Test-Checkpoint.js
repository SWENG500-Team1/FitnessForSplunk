//


(function () {
var splunkjs            = require('splunk-sdk');
var ModularInputs       = splunkjs.ModularInputs;
var Logger              = ModularInputs.Logger;
var Event               = ModularInputs.Event;
var Scheme              = ModularInputs.Scheme;
var Argument            = ModularInputs.Argument;
var utils               = ModularInputs.utils;
var fs                  = require("fs");
var path                = require("path");
var unirest = require("unirest");
var reqauth = unirest("POST", "https://login.live.com/oauth20_token.srf");
var requser = unirest("GET", "https://api.microsofthealth.net/v1/me/PROFILE");
var _storeClientID      = 'dsdfasdf34wsghyt43szx*4xdfga32#4';
var _storeClientKey     = 'a2sdf4sare3425xacxv43';
var _refreshToken       = 'N/A';
var _authToken          = 'N/A';
var _username           = 'N/A';
var _fullname           = 'N/A';
var _lastSyncDate       = 'N/A';
var _lastCheckpointDate = 'N/A';
var _startDate          = 'N/A';
var _endDate            = 'N/A';
var _name               = 'TEST';
var _INPUTDate          = 'N/A';
var checkpointFileContents = "";

var testdate            = "2016/8/5";
var Step                = "0";

var RunTime             = new Date();
var processRunTime      = RunTime.toString();
var mynote              = '"Notes":"Test"';
var errorFound          = false;
var checkpointexist     = false;
//var maxAPIpullcount = 40;
//var APIpullcount    = 0;
var APIIterate  = 0;

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
        var d = new Date(date);
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
  
exports.streamEvents = function (name, singleInput, eventWriter, done) {//Streaming
        _username = singleInput.username;
        _fullname = singleInput.fullname;
        var name = singleInput.name;
        _INPUTDate = singleInput.startdate;
//--------------------------------------------------------
        var checkpointDir = this._inputDefinition.metadata["checkpoint_dir"];
        var checkpointFilePath = path.join(checkpointDir, " MSHEalthCPs_" + name + ".txt");
        // Set the temporary contents of the checkpoint file to an empty string
        var checkpointFileContents = "";
        try {
            checkpointFileContents = utils.readFile("", checkpointFilePath);
        }
        catch (e) {
            fs.appendFileSync(checkpointFilePath, "");
        }
       
        
        
//1-Start-----------------------------------------------------------------------------------------------
//1--Gather CLientID and Client Secret Key from storage passwords
//1-----------------------------------------------------------------------------------------------------

            
        
        
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
                         //------------------------------------------
                        //--Get / Setup Checkpoint
                        //------------------------------------------
                   
                      

                        requser.headers({"authorization": "Bearer " + _authToken});
                        requser.timeout(30000);
                        requser.end(function (resUser) {
                        if (resUser.error) {
                                throw new Error(resUser.error);
                            }
                            var resultprof = JSON.stringify(resUser.body);
                            _lastSyncDate = GetBeginOfHourDate(resUser.body.lastUpdateTime);
                            console.log(_lastSyncDate);
                            var testDateCheck = new Date(_lastCheckpointDate);
                            var testDateSync = new Date(_lastSyncDate);
                            

//-=------------------------------------------------------------
                            if (checkpointFileContents.length <= 0) {// Check point never logged
                                _endDate = GetBeginOfHourDate(_lastSyncDate);
                                _lastCheckpointDate = _endDate;//last sync date!!!
                                _startDate = GetBeginOfHourDate(_INPUTDate)
                                Step = '1';
                            }
                            else {//Log Exist
                                if (typeof checkpointFileContents.lastCheckpointDate != 'undefined') {
                                    //It is a Date
                                _endDate = GetBeginOfHourDate(_lastSyncDate);
                                _lastCheckpointDate = endDate;//last sync date!!!
                                _startDate = checkpointFileContents.lastCheckpointDate;
                                Step = '2';
                                }
                                else {//Log was an Error / Allow for reset
                                    _endDate = GetBeginOfHourDate(_lastSyncDate);
                                    _lastCheckpointDate = GetBeginOfHourDate(_endDate);//last sync date!!!
                                    _startDate = _endDate;
                                    Step = '3';
                                }

                            }
                            mynote = checkpointFileContents;
                            APIpullcount = GetDateDifference(_startDate,_endDate);
//-=------------------------------------------------------------



                            
                         //   mynote = '"_lastSyncDate":"'+testDateSync+'","_lastCheckpointDate":"'+testDateCheck+'","input":"'+_INPUTDate+'",'+checkpointFileContents;
                            var curEventprof = new Event({
                                    stanza: name,
                                    data: '{"username":' + '"' + _username + '"' +
                                    ',"fullname":' + '"' + _fullname + '"' +
                                    ',"notes":' + '{' + mynote + '}' +
                                    ',"processRunTime":' + '"' + processRunTime + '"' +
                                    ',"dataType":' + '"' + 'User-Profile' + '"' +
                                    ',"recordCount":' + '"' + APIpullcount + '"' +
                                    ',"data":' + resultprof + '}'
                                });
                                try {
                                    //console.log(curEventprof);
                                    eventWriter.writeEvent(curEventprof);                                                               
                                    
                                    
                                    var startObj = '  "lastcheckpoint": { "name": "' + _username + '"' +
                                        ',"lastprocessrun": "' + processRunTime + '"' +
                                        ',"startdate": "' + _startDate + '"' +
                                        ',"enddate": "' + _endDate + '"' +
                                        ',"lastCheckpointDate": "' + _lastCheckpointDate + '"' +
                                        ',"lastSyncDate": "' + _lastSyncDate + '"' +
                                        ',"Step": "' + GetBeginOfHourDate(_INPUTDate) + '"' +
                                        '}';

                                    fs.truncateSync(checkpointFilePath, 0);             //Clear Check point
                                    fs.appendFileSync(checkpointFilePath, startObj);    //Write Check point
                                  
                                    }
                                catch (e) {
                                    errorFound = true;
                                    Logger.error(name, e.message);
                                    done(e);
                                    //--We had an error, die
                                    return;
                                }
                            
                        });

                        //4-END-------------------------------------------------------------------------------------------------
                    });

                    //3-END-------------------------------------------------------------------------------------------------

            //2-END-------------------------------------------------------------------------------------------------
            
        
        
      
      
    //1-END-------------------------------------------------------------------------------------------------
    
    };//Streaming

    ModularInputs.execute(exports, module);
})();
