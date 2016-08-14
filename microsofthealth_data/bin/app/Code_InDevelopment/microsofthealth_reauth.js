module.exports = function (_clientid, _clientkey, _token) {
   
    var uniRest = require("unirest");
    var reqAuth = uniRest("POST", "https://login.live.com/oauth20_token.srf");
    this.newAccessToken = 'NAT'; 
    var tempError = 'N/A';
    this.ErrorText = '"Error":"'+tempError+'"';
    var AccessToken = 'TK';
    

    reqAuth.query({
        "redirect_uri": "http://127.0.0.1:3000"
       
    });
    //------------------------------------------
    //--Set the auth format
    //------------------------------------------
    reqAuth.headers({
        "content-type": "application/x-www-form-urlencoded"
    });
    //------------------------------------------
    //--Auth Setup to get new token
    //------------------------------------------
    reqAuth.form({
        "client_id": _clientid,
        "grant_type": "refresh_token",
        "refresh_token": _token,
        "client_secret": _clientkey
    });

try {
  reqAuth.end(tempError =function (resAuth) {
      //    = resAuth.body.access_token;
        //------------------------------------------
        //--Data Set the Auth Token
        //------------------------------------------
       
       
        
      
        
    });//Endauth
} catch (error) {
   tempError=error;
}
    
 //this.newAccessToken =test.AccessToken;


this.ErrorText = '"Error":"'+ tempError+'"';

};