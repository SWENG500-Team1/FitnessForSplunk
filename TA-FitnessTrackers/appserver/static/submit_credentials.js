require([
     "jquery",
     "splunkjs/mvc/simplexml/ready!"
     ], function(
         $
     ) {
         $("#btn-submit").on("click", function (){
             alert("Do Stuff");
         });
     });
