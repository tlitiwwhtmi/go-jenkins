var int;
var currentData;
var xhr;
var firstLoad = true;

function loadLog(){
	if(firstLoad){
		$('#logloading').css('display','');
	}
	if(xhr){
		xhr.abort();
	}
	xhr = $.get("getlog?jobid="+$('#jobid').val()+"&logname="+$('#logname').val(),function(data,status){
		if(firstLoad){
			$('#logloading').css('display','none');
			firstLoad = false;
		}

		$('#logname').val(data.LogName)
		if(data.LogName != ""){
			$('#rawlog').removeAttr("disabled")
			$('#rawlog').attr('href',"http://build.dds.com:8081/logs/"+data.LogName+".txt")
		}else{
			$('#rawlog').attr("disabled","disabled")
			$('#rawlog').removeAttr('href')
		}
		if(data.NeedRead){
			if(JSON.stringify(data) != JSON.stringify(currentData)){
				currentData = data;
				$('#console').html(generateLog(currentData))
				window.scroll(0,$('#console')[0].scrollHeight);
			}
		}
		setTimeout(function(){loadLog();},3000)
	});
}


function generateLog(data){
	var log = data.Log
	if(log.indexOf('[PostBuildScript]') != -1){
		log = log.slice(0,log.indexOf('[PostBuildScript]'))
	}
	var logarr = log.split("\n")
	var leftlog = ""
	if(logarr.length > 200){
		var leftlogarr = logarr.slice(logarr.length-200,logarr.length)
		for(var i=0;i<leftlogarr.length;i++){
			leftlog += leftlogarr[i] + "\n"
		}
	}else{
		for(var i=0;i<logarr.length;i++){
			leftlog += logarr[i] + "\n"
		}
	}
	//log = log.replace("<script","&lt;script");
	//log = log.replace(/\n/g," <p>");
	var htmlstr = '';
	htmlstr += '<pre style="white-space: pre-wrap;white-space: -moz-pre-wrap;white-space: -pre-wrap;white-space: -o-pre-wrap;word-wrap: break-word; ">';
	htmlstr += leftlog;
	if(leftlog){
		if(!data.IsRunning){
			htmlstr += 'Finished:'+data.OutCome
		}else{
			htmlstr += '<p><img src="./static/images/loading.gif" style="height:24px" /></p>'
		}
	}
	htmlstr += '</pre>';
	return htmlstr;
}

loadLog();
//initTimer();

function initTimer(){
	if(int){
		clearInterval(int)
	}
	int = setInterval('loadLog()',2000)
}