// Copyright (C) Microsoft Corporation.  All rights reserved.
// This is a Javascript library that can be used to communicate with the Microsoft Health Cloud API
// directly using CORS. This library allows you to get an access token from MSA and then query the API

/// options supports the following parameters:
///     clientId - your client id as registered in https://account.live.com/developers/applications/index
///     scope - a space separated list of permissions to request to the user, the supported values are
///             mshealth.ReadProfile mshealth.ReadDevices mshealth.ReadActivityHistory mshealth.ReadActivityLocation
///             for an updated list of scopes please visit http://developer.microsofthealth.com
///     redirectUri - the redirect uri registered in MSA

(function (global) {
    "use strict";

    var loginUrl, clientId, scope, redirectUri, accessToken, authHeader, apiTemplateUrl, Promise, MicrosoftHealth;

    loginUrl = "https://login.live.com/oauth20_authorize.srf?client_id={client_id}&scope={scope}&response_type=token&redirect_uri={redirect_uri}";
    accessToken = "";
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

        var xmlHttpRequest, url, promise, queryParameters;
        promise = new Promise();
        xmlHttpRequest = new global.XMLHttpRequest();

        queryParameters = options.parameters ? getParameters(options.parameters) : "";

        url = apiTemplateUrl.replace("{path}", options.path).replace("{parameters}", queryParameters);
        xmlHttpRequest.open(options.method, url, true);
        xmlHttpRequest.setRequestHeader("Authorization", authHeader.replace("{access_token}", accessToken));

        xmlHttpRequest.onload = function () {
            var request = this;

            if (request.status >= 200 && request.status < 300) {
                promise.resolve(JSON.parse(request.responseText));
            } else {
                promise.reject(request.responseText ? JSON.parse(request.responseText) : {});
            }
        };

        xmlHttpRequest.onerror = function () {
            var request = this;

            promise.reject(request.responseText ? JSON.parse(request.responseText) : {});
        };

        xmlHttpRequest.send();

        return promise;
    }

    MicrosoftHealth.prototype.login = function () {
        var hash, parameters, url;

        hash = global.location.hash;

        if (hash) {
            parameters = parseParameters(hash.substring(1));
            accessToken = parameters["access_token"];
        } else {
            url = loginUrl.replace("{client_id}", clientId).replace("{scope}", scope).replace("{redirect_uri}", encodeURIComponent(redirectUri));
            global.location = url;
        }
    };

    MicrosoftHealth.prototype.getProfile = function () {
        return query({
            path: "Profile",
            method: "GET"
        });
    };

    MicrosoftHealth.prototype.getSummaries = function (options) {
        if (!options || !options.period) {
            throw "A period is required to call the summaries API";
        }

        return query({
            path: "Summaries/" + options.period,
            method: "GET",
            parameters: {
                startTime: options.startTime,
                endTime: options.endTime,
                deviceIds: options.deviceIds,
                maxPageSize: options.maxPageSize
            }
        });
    };

    MicrosoftHealth.prototype.getActivities = function (options) {
        return query({
            path: "Activities",
            method: "GET",
            parameters: {
                activityIds: options.activityIds,
                activityTypes: options.activityTypes,
                activityIncludes: options.activityIncludes,
                splitDistanceType: options.splitDistanceType,
                startTime: options.startTime,
                endTime: options.endTime,
                deviceIds: options.deviceIds,
                maxPageSize: options.maxPageSize
            }
        });
    };

    MicrosoftHealth.prototype.getDevices = function (deviceId) {
        var path = deviceId ? "Devices/" + encodeURIComponent(deviceId) : "Devices";

        return query({
            path: path,
            method: "GET"
        });
    };

    Promise = function () { };

    Promise.prototype.then = function (onResolved, onRejected) {
        this.onResolved = onResolved;
        this.onRejected = onRejected;
    };

    Promise.prototype.resolve = function (value) {
        this.onResolved(value);
    };

    Promise.prototype.reject = function (value) {
        this.onRejected(value);
    };

})(this);