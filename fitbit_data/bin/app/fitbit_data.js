(function() {
    
    // Splunk SDK variables
    var splunkjs        = require("splunk-sdk");
    var ModularInputs   = splunkjs.ModularInputs;
    var Logger          = ModularInputs.Logger;
    var Event           = ModularInputs.Event;
    var Scheme          = ModularInputs.Scheme;
    var Argument        = ModularInputs.Argument;
    var utils           = ModularInputs.utils;

    // Fitbit variables
    var FitbitClient    = require("fitbit-client-oauth2");
    var clientId        = '227MVJ';
    var clientSecret    = 'df8009bd0ddcb975f9a812e3587e54dd';
    var client          = new FitbitClient(clientId, clientSecret);
        
    // Miscellaneous variables
    var fs              = require("fs");
    
    // getScheme method returns the introspection scheme
    exports.getScheme = function() {

        var scheme = new Scheme("Fitbit Data");

        scheme.description = "Streams events containing Fitbit data from Fitbit API.";
        scheme.useExternalValidation = true;    // if true, must define validateInput method
        scheme.useSingleInstance = true;        // if true, all instances of mod input passed to
                                                //   a single script instance; if false, user 
                                                //   can set the interval parameter under "more settings"

        // Arguments
        scheme.args = [
            new Argument({
                name: "username",
                dataType: Argument.dataTypeString,
                description: "Fitbit account name.",
                requiredOnCreate: true,
                requiredOnEdit: false
            }),
            /*new Argument({
                name: "access_token",
                dataType: Argument.dataTypeString,
                description: "Fitbit API access token for the user account.",
                requiredOnCreate: false,
                requiredOnEdit: false
            }),
            new Argument({
                name: "refresh_token",
                dataType: Argument.dataTypeString,
                description: "Fitbit API refresh token for the user account.",
                requiredOnCreate: false,
                requiredOnEdit: false
            }),*/
            new Argument({
                name: "date",
                dataType: Argument.dataTypeString,
                description: "Date of Fitbit data to retrieve.",
                requiredOnCreate: false,
                requiredOnEdit: false
            }),
            new Argument({
                name: "token_json",
                dataType: Argument.dataTypeString,
                description: "JSON token to access Fitbit API for the user account.",
                requiredOnCreate: false,
                requiredOnEdit: false
            })
        ];

        return scheme;
    };

    // validateInput method validates the script's configuration (optional)
    exports.validateInput = function(definition, done) {

        // Local variables
        //var username = definition.parameters.username;
        //var access_token = definition.parameters.access_token;
        //var refresh_token = definition.parameters.refresh_token;
        var date = definition.parameters.date;
        var token_json = definition.parameters.token_json;

        // Validation checking
        // TODO: Perform validation on date
        /*if (min >= max) {
            done(new Error("min must be less than max; found min=" + min + ", max=" + max));
        }
        else if (count < 0) {
            done(new Error("count must be a positive value; found count=" + count));
        }
        else {
            done();
        }*/
        done();
    };

    // Stream data as text or as XML, using checkpoints
    exports.streamEvents = function(name, singleInput, eventWriter, done) {        
        // TODO: Refactor data retrieval logic into a function
        // Modular input logic
        /*var getDailyActivitySummary = function (token, options) {
            
            client.getDailyActivitySummary(token, options)
                .then(function(res) {
                    var result = JSON.stringify(res);
                    fs.writeFile('./success.json', result, function (err,data) {
                        if (err) {
                            return console.log(err);
                        }
                        //console.log(result);
                    });
                    //console.log('Fitbit data retrieved for', options.date);
                }).catch(function(err) {
                    
                    
                    var error = JSON.stringify(err);
                    fs.writeFile('./error.json', error, function (err,data) {
                        if (err) {
                            return console.log(err);
                        }
                        //console.log(err);
                    });
                    //console.log('Error retrieving Fitbit data.', error);
                });
        };*/
        
        // Local variables
        var username = singleInput.username;
        //var access_token = singleInput.access_token;
        //var refresh_token = singleInput.refresh_token;
        var fit_date = singleInput.date;
        var token_json = JSON.parse(singleInput.token_json);    // Deserialize JSON token
        
        // Set token and options
        var token = token_json.token;
        var options = { date: "2016-06-01" };
        options.date = fit_date;
        
        // ----------- need to refactor ------------
        
        // Retrieve Fitbit data for specified date
        client.getDailyActivitySummary(token, options)
            .then(function(res) {
                // Serialize JSON result
                var result = JSON.stringify(res);
                
                // Package data in an Event object
                var curEvent = new Event({
                    stanza: name,
                    data: "{\"username\":" + username + ",\"date\":"+fit_date+",\"data\":" + result + "}"
                });
                
                // Stream event
                try {
                    eventWriter.writeEvent(curEvent);
                }
                catch (e) {
                    //errorFound = true;
                    Logger.error(name, e.message);
                    done(e);
                    
                    // We had an error, die
                    return;
                }
            }).catch(function(err) {
                // Serialize JSON error
                var error = JSON.stringify(err);
                Logger.error(name, error);
                done(e);
                // We had an error, die
                return;
            })
            
        // -----------------------------------------
        
        // We're done
        done();
    };

    ModularInputs.execute(exports, module);
})();