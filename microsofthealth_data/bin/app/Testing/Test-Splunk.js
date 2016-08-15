
(function () {
    var splunkjs = require('splunk-sdk');
    var ModularInputs = splunkjs.ModularInputs;
    var Logger = ModularInputs.Logger;
    var Event = ModularInputs.Event;
    var Scheme = ModularInputs.Scheme;
    var Argument = ModularInputs.Argument;
    var utils = ModularInputs.utils;


    var _refreshToken = 'N/A';
    var _username = 'N/A';
    var _fullname = 'N/A';
    var _startDate = 'N/A';
    var _endDate = 'N/A';
    var _name = 'JAY';



    var testdate = "2016/8/5";

    var RunTime = new Date();
    var processRunTime = RunTime.toString();
    var mynote = 'Test';
    var errorFound = false;





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


    exports.streamEvents = function (name, singleInput, eventWriter, done) {
        _username = singleInput.username;
        _fullname = singleInput.fullname;
        _name = singleInput.name;
        _refreshtoken = singleInput.token_json;


        //-Start-----------------------------------------------------------------------------------------------
        //--Splunk Data Test
        //-----------------------------------------------------------------------------------------------------



        var curEvent = new Event({
            stanza: _name,
            data: '{"username":' + '"' + _username + '"' +
            ',"fullname":' + '"' + _fullname + '"' +
            ',"notes":' + '{' + mynote + '}' +
            ',"processRunTime":' + '"' + processRunTime + '"' + '}'
        });
        try {


            console.log(curEvent);
            eventWriter.writeEvent(curEvent);

        }
        catch (e) {
            errorFound = true;
            Logger.error(name, e.message);
            // done(e);
            //--We had an error, die
            return;
        }




        //-END-------------------------------------------------------------------------------------------------





    }
});