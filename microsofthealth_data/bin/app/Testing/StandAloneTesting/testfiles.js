//////////////////////////////////////
var fs = require('fs');
var savedDate = '';
var endDate = '';
var errorFound = false;
var dateRecordPulled = false;
var lastID = 0;

//------------------------------------------
//--PreCalculate End Date
//------------------------------------------
var d = new Date();
var d2 = new Date();
d2.setMinutes(d.getMinutes() - d.getMinutes());
d2.setSeconds(d.getSeconds() - d.getSeconds());
d2.setMilliseconds(d.getMilliseconds() - d.getMilliseconds());
  endDate = d2.toISOString();

function getDateVars(_name,_lastdaterun,cb){
fs.stat('File.txt', function(err, stat) {

//==========================================================================
//--Determine Last Input Date
//==========================================================================
//------------------------------------------
//--Log File Exist
//------------------------------------------
    if(err == null) {
        console.log('File exists');
        var fs2 = require('fs');
        var obj = JSON.parse(fs.readFileSync('File.txt', 'utf8'));
      //console.log(obj.dateparm.length);
        for(var j = 0; j < obj.dateparm.length && !errorFound; j++){
          lastID = obj.dateparm[j].id;
          console.log('last id:'+  lastID);

          if (obj.dateparm[j].name == _name){
            console.log(obj.dateparm[j]);
            savedDate = obj.dateparm[j].lastdaterun;
            var replaceOld = JSON.stringify(obj.dateparm[j]);
            var replaceNew = '{"id": 0, "name": "'+_name+'", "lastdaterun": "'+endDate+'"} ';
            var replaceStr = JSON.stringify(obj);
            var replaceEnd = replaceStr.replace(replaceOld,replaceNew);
            fs.writeFile('File.txt', replaceEnd);
              //console.log('old:'+replaceStr);
              //console.log('new:'+replaceEnd);
              dateRecordPulled = true;
            cb();
            break;//Exit Loop - Record Found
          }

        }
        //------------------------------------------
        //--Log File Exist but user date record does not
        //------------------------------------------
        if (dateRecordPulled==false)
        {
          var today = new Date(_lastdaterun);
          var tempDate = today.toISOString();
          savedDate  = tempDate;
          var newID = lastID + 1;
            var allData = JSON.stringify(obj);
            var newAllData = allData.replace(']}',',{"id":'+newID+', "name": "'+_name+'", "lastdaterun": "'+endDate+'"} ]}');
            fs.writeFile('File.txt', newAllData);
            cb();
        }
      }
//------------------------------------------
//--Log File does not Exist / Create New File
//------------------------------------------
  else if(err.code == 'ENOENT') {
        // file does not exist
        var today = new Date(_lastdaterun);
        var tempDate = today.toISOString();
        savedDate  = tempDate;
        var startObj ='{  "dateparm": [{"id": 1, "name": "'+_name+'", "lastdaterun": "'+endDate+'"} ]}';
        fs.writeFile('File.txt', startObj);
        cb();
    } else {
        console.log('Some other error: ', err.code);
    }

});

};


getDateVars("Jay","07/29/2016 09:00",function(){

  console.log('s:'+savedDate);
  console.log('e:'+endDate);
});
