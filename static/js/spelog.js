function loadLog(){
	$.get("getspelog?hid="+$('#hid').val(),function(data,status){
	//$.get("http://10.47.22.62:8000/test.txt",function(data,status){
		$('#logloading').css('display','none');
		$('#console').html(generateLog(data))
		window.scroll(0,$('#console')[0].scrollHeight);
	});
}


function generateLog(data){
	var log = data.Log
	if(log.indexOf('[PostBuildScript]') != -1){
		log = log.slice(0,log.indexOf('[PostBuildScript]'))
	}
	//log = log.replace("<script","&lt;script");
	//log = log.replace(/\n/g," <p>");

	if(data.FileName != ''){
		$('#rawlog').removeAttr("disabled")
		$('#rawlog').attr('href',"http://build.dds.com:8081/logs/"+data.FileName+".txt")
	}else{
		$('#rawlog').attr("disabled","disabled")
		$('#rawlog').removeAttr('href')
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
	var htmlstr = '';
	htmlstr += '<pre style="white-space: pre-wrap;white-space: -moz-pre-wrap;white-space: -pre-wrap;white-space: -o-pre-wrap;word-wrap: break-word; ">';
	htmlstr += leftlog;
	if(leftlog){
			htmlstr += 'Finished:'+data.OutCome
	}
	htmlstr += '</pre>';
	return htmlstr;
}

loadLog();
