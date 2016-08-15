var splunkjs = require('splunk-sdk');
var ModularInputs = splunkjs.ModularInputs;
var Logger = ModularInputs.Logger;
var Event = ModularInputs.Event;
var Scheme = ModularInputs.Scheme;
var Argument = ModularInputs.Argument;
var utils = ModularInputs.utils;
var fs = require("fs");
var path = require("path");
var AdminUserName = 'jsprowls';
var AdminPassword = '38der#IyLL0%n%NI@00n#lc3f';
var service = new splunkjs.Service({ username: AdminUserName, password: AdminPassword });

var _storeClientID = 'N/A';
var _storeClientKey = 'N/A';
var _refreshToken = 'N/A';




var testdate = "2016/8/5";

var RunTime = new Date();
var processRunTime = RunTime.toString();
var mynote = 'Test';
var errorFound = false;
//var maxAPIpullcount = 40;
//var APIpullcount    = 0;
var APIIterate = 0;


service.login(function (err, success) {
    if (err) {
        throw err;
    }

    console.log("Login was successful: " + success);
    //--LoggedIn


   

    // Create a Service instance and log in 
    var service2 = new splunkjs.Service({
        username: AdminUserName,
        password: AdminPassword,
        scheme: "https",
        host: "localhost",
        port: "8089",
        version: "5.0"
    });

    // Print installed apps to the console to verify login



    var endpoint = new splunkjs.Service.Endpoint(service2, "servicesNS/nobody/fitness_for_splunk/storage/collections/data/microsoft_tokens");
    endpoint.get("results", { offset: 1 }, function (res) {
        console.log(res);

    });





    //1-Start-----------------------------------------------------------------------------------------------
    //1--Gather CLientID and Client Secret Key from storage passwords
    //1-----------------------------------------------------------------------------------------------------
    service.storagePasswords().fetch(
        function (err, storagePasswords) {
            //   console.log(storagePasswords.list());
            if (err)
            { /* handle error */ }
            else {
                // Storage password was created successfully
                if (storagePasswords.list().length > 0) {
                    for (var i = 0; i < storagePasswords.list().length; i++) {
                        //  console.log( storagePasswords.list()[i]);
                        if (storagePasswords.list()[i]._properties.realm = 'microsoft') {
                            console.log("realm:" + storagePasswords.list()[i]._properties.realm);


                            _storeClientID = storagePasswords.list()[i]._properties.username
                            console.log("_storeClientID:" + _storeClientID);
                            _storeClientKey = storagePasswords.list()[i]._properties.clear_password
                            console.log("_storeClientKey:" + _storeClientKey)


                            //2-END-------------------------------------------------------------------------------------------------
                        }
                    }



                }

            }
        });
    //1-END-------------------------------------------------------------------------------------------------



});//EndLogin