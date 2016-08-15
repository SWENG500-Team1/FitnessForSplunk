



var loginUrl, clientId, scope, redirectUri, accessToken, authHeader, apiTemplateUrl,  MicrosoftHealth;

loginUrl = "https://login.live.com/oauth20_authorize.srf?client_id={client_id}&scope={scope}&response_type=token&redirect_uri={redirect_uri}";
accessToken = "EwBIA/F0BAAUvSQiG6C8oi/OrqaVrv8+s9GGlDkAAY6Vn5+a6Y8tE5HsTtAR9fYgEnYYl6Fn9hL4EDG6+MJD+eH2ot0kf4K8boKzf66OWBe5tXXb4Oaj2ALFWoiXl6DjS6liZiKHM8rCrNNu1W5MfXIfCQHz0T1e5EDlvhoXQljSLQNJhMJX8tTPOUS00DqJevQ/E/cU9+DkcDTHgCQxFXJYWfEGfbFVCoKNG7WhIGb5aaPmDPO2m1xUtMOK/PVzdOCapDzkCMnWyvi8nif8OF7FsP8s0I7orwtsv4qLjKkNtlkn5HzXX5Qge7x6ysMnsHKWuNe3LSlAItU+uiLsnYKsz6xP7gSbXlvqiaoZTCw0yh1gnr83pSjLVcZhJ68DZgAACGb18aOW/e/tGAL2o+nkFkFSUBGBy8bgQwWgN9D8JrbVb0Wx55BEMUneDSqp907W2MryxxqlknqdTigHVGHXgmlEU751LpSu6ZRReCUngREo5Cl992BXSRyuseCGKazSP+LU+6QOGihjFDQT6S77QcO5kAXT4hkLMTfeWZrpPwTRlxdtxCKCdRY+AlF42VWCrerI6suTEYxL1kuri1IDLD3A3Tr+NIwQGC7+uLkrIYBM4NT3NvRsU9i2THwe4KUQ1fF3Y9ubRXZqruyXwmX9p44jBj0n/pEcnftnOnY+Tlm4BRWQf9zghQ9V7vHgwdk8rCA4SBBHhkuqyCDrCIJOu09snRBiuBxtmw";
authHeader = "Bearer {access_token}";
apiTemplateUrl = "https://api.microsofthealth.net/v1/me/{path}?{parameters}";

MicrosoftHealth = function (options) {
    clientId = options.clientId;
    scope = options.scope;
    redirectUri = options.redirectUri;
};

global.MicrosoftHealth = MicrosoftHealth;

function parseParameters(hash) {
    var split, dictionary, i, param;

    dictionary = {};
    split = hash.split("&");

    for (i = 0; i < split.length; i++) {
        param = split[i].split("=");
        if (param.length === 2) {
            dictionary[param[0]] = param[1];
        }
    }

    return dictionary;
}


function getParameters(parameters) {
    var queryParameters, p;

    queryParameters = "";

    if (parameters) {
        for (p in parameters) {
            if (parameters[p]) {
                queryParameters = queryParameters.concat(encodeURI(p) + "=" + encodeURI(parameters[p]) + "&");
            }
        }
    }

    return queryParameters.substring(0, queryParameters.length - 1);
}


    function query(options) {
        if (!accessToken) {
            throw "User is not authenticated, call login function first";
        }

        var xmlHttpRequest, url,  queryParameters;

        xmlHttpRequest = new global.XMLHttpRequest();

        queryParameters = options.parameters ? getParameters(options.parameters) : "";

        url = apiTemplateUrl.replace("{path}", options.path).replace("{parameters}", queryParameters);
        xmlHttpRequest.open(options.method, url, true);
        xmlHttpRequest.setRequestHeader("Authorization", authHeader.replace("{access_token}", accessToken));

        xmlHttpRequest.onload = function () {
            var request = this;

            if (request.status >= 200 && request.status < 300) {
                //promise.resolve(JSON.parse(request.responseText));
            } else {
              //  promise.reject(request.responseText ? JSON.parse(request.responseText) : {});
            }
        };

        xmlHttpRequest.onerror = function () {
            var request = this;

          //  promise.reject(request.responseText ? JSON.parse(request.responseText) : {});
        };

        xmlHttpRequest.send();

        return "promise";
    }

    MicrosoftHealth.prototype.login = function () {
        var hash, parameters, url;


            url = loginUrl.replace("{client_id}", clientId).replace("{scope}", scope).replace("{redirect_uri}", encodeURIComponent(redirectUri));
            global.location = url;

    };

    MicrosoftHealth.prototype.getProfile = function () {
        return query({
            path: "Profile",
            method: "GET"
        });
    };

/////////////////////////////////////////////////////////////////////////////////////////////////////////

var mshealth = new MicrosoftHealth({
    clientId: "e2a30a63-c396-457a-abbf-1409678a5be4",
    redirectUri: "https://login.live.com/oauth20_desktop.srf",//https://jaysprowls.azurewebsites.net/index.html",
    scope: "mshealth.ReadProfile mshealth.ReadDevices mshealth.ReadActivityHistory mshealth.ReadActivityLocation"

});
mshealth.login();
console.log("Ran Login");

mshealth.getProfile().then(
    function (profile) {
        var greeting = "Hello " + profile.firstName + ". Welcome to the Microsoft Health Cloud API";
        console.log(greeting);

var userinfo = "User: " + profile.firstName + " " +profile.lastName + "\nGender: "+profile.gender + "\nHeight: "+profile.height + "\nWeight: " +profile.weight+"\nPostal Code: "+ profile.postalCode;
console.log(userinfo);
    },
    function (error) {
        onError(error.error,"Profile");
    }
);


console.log("Done");
