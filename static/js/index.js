/**
 * Created by duanhao on 2016/1/19.
 */

var firstFlush = true;
var xhr;
var pacBranch = '';
var pacPage = 1;
var pacCid = 0;
var pacLength = 0;

$('#jobs').jobgreed("");

//$('#jobs').on("click",".toglog",function(event){
//    var log_out = $('#log'+$(this).attr("jobid"));
//	var console = $('#log'+$(this).attr("jobid")+' .log_out')
//    if(log_out.css('height') == '450px'){
//		log_out.next().css('display',"none")
//        log_out.css('height',"0px")
//    }else{
//        log_out.css('height','450px')
//		log_out.next().css('display',"block")
//        console.scrollTop( console[0].scrollHeight );
//    }
//});

//$('#jobs').jobgreed('avc')

function logout(){
	$.post("/logoutaction",function(data,status){
		Cookie.unset("activepro");
		location.href="login.html"
	});
}

                                                                 
$('#addProject').on('shown.bs.modal', function () {
	if(firstFlush){
		$.get("/getprojects",function(data,status){
			if(!data){

			}else{
				if(data.toString().indexOf("DOCTYPE") > 0){
					window.location.href = "login"
				}
				$('#pro').typeahead({
					source:data,
					autoSelect: true,
					items:"all",
					displayText:function(item){
						return item.name_with_namespace;
					},
					showHintOnFocus:true});
			}
			$('#realbd').css('display',"block")
			$('#loadingbd').css('display',"none")
			firstFlush = false
		}); 
	}           
                                                
}); 

//$('#addProject').on('hide.bs.modal', function () {
//	$('#loadingbd').css('display',"block")
//	$('#realbd').css('display',"none") 
//});

$('#pro').change(function() {
    var current = $('#pro').typeahead("getActive");
    if (current) {
        if (current.name_with_namespace == $('#pro').val()) {
			$("#proid").val(current.id)
        } else {
            // This means it is only a partial match, you can either add a new item 
            // or take the active if you don't want new items
        }
    } else {
        alert("Please select from the options")
		$('#pro').focus()
    }
}); 

$('#jobs').on('click','.job_edit',function(e){
	var gridJob = getGridJob($(this).attr('jobid'))
	
	$('#tempjobtag').val(gridJob.Id)
	
	$('#editproject').val($('.proactive .proname').html())
	var htmlstr = '';
	var branches = getGrid('jobs').branches;
	if(branches){
		for(var i=0;i<branches.length;i++){
			htmlstr += '<option>'+branches[i].name+'</option>'
		}
	}
	$('#editbranch').html(htmlstr);
	
	$('#editbranch').val(gridJob.BranchName);
	$('#editjobName').val(gridJob.JobName);
	if(gridJob.IsGitlab){
		$('#editblankCheckbox').prop("checked",true)
	}else{
		$('#editblankCheckbox').prop("checked",false)
	}
	if(gridJob.AutoBuild){
		$('#editautobuild').prop("checked",true)
		$('#editjquerycron').css('display','block')
		$('#editjquerycron').html('');
		var initValue = "42 3 * * 5";
		if(gridJob.AutoTime){
			initValue = gridJob.AutoTime.substring(2);
		}
		$('#editjquerycron').cron({
			initial: initValue,
			useGentleSelect: true
		});
	}else{
		$('#editautobuild').prop("checked",false)
		$('#editjquerycron').html('');
		$('#editjquerycron').css('display','none')
	}

	if(gridJob.Deploy == "1"){
		$('.editdtype[value=1]').prop("checked",true)
	}else if(gridJob.Deploy == "2"){
		$('.editdtype[value=2]').prop("checked",true)
	}else{
		$('.editdtype[value=3]').prop("checked",true)
	}

	if(gridJob.Build == '1'){
		$('.editbtype[value=1]').prop("checked",true)
		$('#editshell').css('visibility',"hidden")
		$('#editshell').css('height',"0px")
		$('#editshell').css('display',"none")
		$("#editfortype").css("display","block");
		$("#editforprofile").css("display","block");
		$("#editforpom").css("display","block");
	}else{
		$('.editbtype[value=2]').prop("checked",true)
		$('#editshell').css('display',"block")
		$('#editshell').css('visibility',"visible")
		$('#editshell').css('height',"150px")
		$("#editfortype").css("display","none");
		$("#editforprofile").css("display","none");
		$("#editforpom").css("display","none");
	}
	$('#editpomsrc').val(gridJob.PomPath);
	$('#editprofile').val(gridJob.Profile)
	$('#editemail').val(gridJob.Email);
	$('#editshell').val(gridJob.Shell);


	//check if it is a android project
	if($('.proactive .proname').attr('type') == 'Android'){
		$('.editbtyperadio').css('display','none');
		showEditShell();
	}
	if($('.proactive .proname').attr('type') == 'Java'){
		$('.editbtyperadio').css('display','block');
	}
	
	$('#editJob').modal('show');
});

function editJob(){
	var gituser = 0
	if($('#editblankCheckbox').is(':checked')){
		gituser = 1
	}
	var autobuild = 0
	if($('#editautobuild').is(':checked')){
		autobuild = 1
	}
	var autoTime = ""
	if($('#editjquerycron').html() != ''){
		autoTime = "0 "+ $('#editjquerycron').cron("value")
	}
	$.post('editjob',{
		jobid:$('#tempjobtag').val(),
		branch:$('#editbranch').val(),
		name:$('#editjobName').val(),
		gitlab:gituser,
		emails:$('#editemail').val(),
		build:$('.editbtype:checked').val(),
		pompath:$('#editpomsrc').val(),
		profile:$('#editprofile').val(),
		deploytype:$('.editdtype:checked').val(),
		shell:$('#editshell').val(),
		autobuild:autobuild,
		autotime:autoTime
	},function(data,status){if(data == "success"){
			$('#editJob').modal('hide')
			initJobList();
		}else{
			alert(data);
		}
	});
}

function addJob(){
	var gituser = 0
	if($('#blankCheckbox').is(':checked')){
		gituser = 1
	}
	var autobuild = 0
	if($('#autobuild').is(':checked')){
		autobuild = 1
	}
	var autoTime = ""
	if($('#jquerycron').html() != ''){
		autoTime = "0 "+ $('#jquerycron').cron("value")
	}
	$.post('addjob',{
		pid:$('.proactive').attr('proid'),
		branch:$('#branch').val(),
		name:$('#jobName').val(),
		gitlab:gituser,
		emails:$('#email').val(),
		build:$('.btype:checked').val(),
		pompath:$('#pomsrc').val(),
		profile:$('#profile').val(),
		deploytype:$('.dtype:checked').val(),
		shell:$('#shell').val(),
		autobuild:autobuild,
		autotime:autoTime
	},function(data,status){if(data == "success"){
			$('#myModal').modal('hide')
			initJobList();
		}else{
			alert(data);
		}
	});
}

function addproject(){
	$.post('addproject',{
		pid:$("#proid").val(),
		lang:$("#lang").val()
	},function(data,status){
		if(data == "success"){
			$('#addProject').modal('hide');
			Cookie.set("activepro",$("#proid").val(),1);
			initProList();
		}else{
			alert(data);
		}
	});
}   


function editprolist(){
	if($('.prodelete').css('width') != '0px'){
		$('.proname').css('width','100%')
		$('.prodelete').css('width','0%')
	}else{
		$('.proname').css('width','80%')
		$('.prodelete').css('width','20%')
	}
}

function initProList(){
	$('#loadinglist').css('display',"block")
	$('#prolist').css('display',"none")
	$.get("ownedproject",function(data,status){

		if(!data){
			$('#prolist').html("there is no project");
		}else {
			if(data.toString().indexOf("DOCTYPE") > 0){
				window.location.href = "login"
			}
			$('#prolist').html("");
			//$('#prolist').append('<li style="margin-top: -10px;margin-bottom: -10px;"><div style="height:100%;width:100%;background: none;border:none;text-align: right;padding-right: 10px;"><span data-toggle="modal" data-target="#addProject"><img style="width:24px" src="./static/images/addpro.png" /></span><span onclick="editprolist()"><img style="width:24px" src="./static/images/removepro.png" /></span></div></li>')
			//$('#prolist').append('<li><a class="btn btn-default" data-toggle="modal" data-target="#addProject" href="#" role="button" style="width: 100%;margin: auto; border-radius: 12px; border-style: dashed;background-color: white">+ Add Project</a></li>');
			//$('#prolist').append('<li><a class="btn btn-default" href="javascript:editprolist()" role="button" style="width: 100%;margin: auto; border-radius: 12px; border-style: dashed;background-color: white">Project</a></li>');
			if(data){
				if(data != 'error'){
					var activePro = Cookie.get("activepro")
					var htmlstr = '';
					var isProInList = false;

					for(var i=0;i<data.length;i++){
						if(data[i].ProjectId == activePro){
							isProInList = true;
							break;
						}
					}

					if(!isProInList){
						for(var i=0;i<data.length;i++){
							if(i ==0){
								htmlstr += '<li class="animated fadeInRight"><div class="proitem proactive" proid="'+data[i].ProjectId+'"><div class="proname" type="'+data[i].Language+'">'+data[i].ProjectName+'</div><div class="prodelete" style="background-color:red;color:white;height:49px"><img style="width:24px" src="./static/images/removepro.png"></div></div></li>'
							}else{
								htmlstr += '<li class="animated fadeInRight"><div class="proitem" proid="'+data[i].ProjectId+'"><div class="proname" type="'+data[i].Language+'">'+data[i].ProjectName+'</div><div class="prodelete" style="background-color:red;color:white;height:49px"><img style="width:24px" src="./static/images/removepro.png"></div></div></li>'
							}
						}
					}else{
						for(var i=0;i<data.length;i++){
							if(data[i].ProjectId == activePro){
								htmlstr += '<li class="animated fadeInRight"><div class="proitem proactive" proid="'+data[i].ProjectId+'"><div class="proname" type="'+data[i].Language+'">'+data[i].ProjectName+'</div><div class="prodelete" style="background-color:red;color:white;height:49px"><img style="width:24px" src="./static/images/removepro.png"></div></div></li>'
							}else{
								htmlstr += '<li class="animated fadeInRight"><div class="proitem" proid="'+data[i].ProjectId+'"><div class="proname" type="'+data[i].Language+'">'+data[i].ProjectName+'</div><div class="prodelete" style="background-color:red;color:white;height:49px"><img style="width:24px" src="./static/images/removepro.png"></div></div></li>'
							}
						}
					}

					//$('#prolist').prepend(htmlstr)
					$('#prolist').append(htmlstr)
				}
			}
		}


		$('#prolist').css('display',"block")
		$('#loadinglist').css('display',"none")
		flushPacList();
		initJobList();
	});
}

function searchPac(){
	pacCid = 0;
	pacPage = 1;
	pacLength = 0;
	initPacList();
}

function flushPacList(){
	pacCid = 0;
	pacPage = 1;
	pacBranch = '';
	pacLength = 0;
	initPacList();
}

function changePacBranch(name){
	if(!name){
		name=""
	}
	if(pacBranch == name){
		return;
	}
	pacCid = 0;
	pacPage = 1;
	pacLength = 0;
	pacBranch = name;
	initPacList();
}

function nextpage(){
	$('#pacloadbtn').remove();
	pacPage++;
	initPacList();
}

function initPacList(){
	$('#pacloading').css('display',"block")
	$('#paccontent').css('display',"none")
	var projectId = $('.proactive').attr('proid');
	if(!projectId){
		var pachtml = '<div style="font-size: 20px;margin-top: 100px;color: rgba(134, 122, 122, 0.95);">there is no package</div>'
		$('#paclist').html('')
		$('#paclist').html(pachtml)
		$('#pacloading').css('display',"none")
		$('#paccontent').css('display',"block")
		return;
	}
	$.get('listpac?pid='+projectId+'&cid='+pacCid+'&page='+pacPage+'&branch='+pacBranch+'&q='+$('#pacQ').val(),function(data){
		if(data.toString().indexOf("DOCTYPE") > 0){
			window.location.href = "login"
		}
		if(data != "error"){
			var branches = data.Branches
			var htmlstr = '';
			if(pacBranch != ''){
				htmlstr += '<a class="more_button" style="text-decoration: none;">'+pacBranch+'<span class="glyphicon glyphicon-chevron-down dropdown-arrow"></span></a>'
			}else{
				htmlstr += '<a class="more_button" style="text-decoration: none;">All Branches<span class="glyphicon glyphicon-chevron-down dropdown-arrow"></span></a>'
			}
			htmlstr += '<ul class="more-dropdown jobbranch">'
			htmlstr += '<li><a href="javascript:changePacBranch()" class="ember-view">All Branches</a></li>'
			for(var i=0;i<branches.length;i++){
				htmlstr += '<li><a href="javascript:changePacBranch(\''+branches[i].name+'\')" class="ember-view">'+branches[i].name+'</a></li>'
			}
			htmlstr += '</ul>'
			$('#pacbranch').html('');
			$('#pacbranch').html(htmlstr);

			var paclist = data.GenPacs;
			var pachtml = '';
			var tempCid = 0;
			if(paclist){
				for(var i=0;i<paclist.length;i++){
					var genPac = paclist[i];
					pachtml += '<div class="pacblock">';
					pachtml += '<div style="float: left;width: 5%"><label>#'+(pacLength+i+1)+'</label></div>'
					pachtml += '<div style="float: left;width: 5%"><label>'+genPac.History.Branch+'</label></div>'
					pachtml += '<div style="float: left;width: 20%"><label>'+formattime(genPac.History.StartTime)+'</label></div>'
					pachtml += '<div style="float: left;width: 25%"><a href="http://git01.dds.com/'+formatProjectPath($('.proactive .proname').html())+'/commit/'+genPac.History.Version+'" target="_blank"><nobr>'+formateMessage(genPac.History.Message)+'-'+genPac.History.CommitAuthor+'</nobr></a></div>'
					pachtml += '<div style="float: left;width: 10%"><label>'+genPac.History.BuildExecutor+'</label></div>'
					pachtml += '<div style="float: left;width: 35%">'
					for(var j=0;j<genPac.Pacs.length;j++){
						var pac = genPac.Pacs[j]
						pachtml += '<label>'+pac.Name+'</label>'
						if(tempCid == 0){
							tempCid = pac.Id
						}else{
							if(tempCid > pac.Id){
								tempCid = pac.Id
							}
						}
					}
					pachtml += '</div>'
					pachtml += '<div style="clear: both"></div>'
					pachtml += '</div>'
				}
				pacLength += paclist.length

				if(data.Total == 10){
					pachtml += '<button id="pacloadbtn" onclick="javascript:nextpage()" type="button" class="btn btn-default btn-lg" style="width: 95%;margin-top: 20px">More</button>'
				}else{
					if($('#pacloadbtn').length > 0){
						$('#pacloadbtn').remove();
					}
				}
			}else{
				if(pacPage == 1){
					pachtml += '<div style="font-size: 20px;margin-top: 100px;color: rgba(134, 122, 122, 0.95);">The Project did not make any package till now</div>'
				}else{
					if($('#pacloadbtn').length > 0){
						$('#pacloadbtn').remove();
					}
				}
			}
			if(pacCid == 0){
				$('#paclist').html('')
				$('#paclist').html(pachtml)
			}else{
				$('#paclist').append(pachtml)
			}
			pacCid = tempCid
			$('#pacloading').css('display',"none")
			$('#paccontent').css('display',"block")
		}
	});
}

function initJobList(){
	var projectId = $('.proactive').attr('proid')
	jobGrid = getGrid('jobs')
	jobGrid.pageNum = 1
	jobGrid.branchName = "";
	jobGrid.setProjectId(projectId);
	jobGrid.loadData();
}

initProList();

$('#prolist').on("click","li .prodelete",function(e){
	$.get("deleteproject?pid="+$(this).parent().attr('proid'),function(data,status){
		if(data.toString().indexOf("DOCTYPE") > 0){
			window.location.href = "login"
		}
		if(data != "success"){
			alert(data);
		}
	});
	if($(this).parent().attr('class') == "proitem proactive"){
		$(this).parent().remove();
		if($('#prolist li .proitem').length != 0){
			$($('#prolist li .proitem')[0]).attr("class","proitem proactive");
		}
		flushPacList();
		initJobList();
	}else{
		$(this).parent().remove();
	}
});

$('#prolist').on("click","li .proname",function(e){
	$('#prolist li .proname').each(function(){
		if($(this).parent().attr("class") == "proitem proactive"){
			$(this).parent().attr("class","proitem");
			return;
		}
	});
	$(e.target).parent().attr("class","proitem proactive");
	Cookie.set("activepro",$(e.target).parent().attr("proid"),1);
	flushPacList();
	initJobList();
});

$('#myModal').on('shown.bs.modal', function () {
	$('#project').val($('.proactive .proname').html())
	var htmlstr = '';
	var branches = getGrid('jobs').branches;
	for(var i=0;i<branches.length;i++){
		htmlstr += '<option>'+branches[i].name+'</option>'
	}

	$('#branch').html(htmlstr);
	$("#fortype").css("display","block");
	$("#forprofile").css("display","block");
	$("#forpom").css("display","block");

	//check if it is a android project
	if($('.proactive .proname').attr('type') == 'Android'){
		$('.btyperadio').css('display','none');
		showShell();
	}
	if($('.proactive .proname').attr('type') == 'Java'){
		$('.btyperadio').css('display','block');
		hideShell();
	}

});

function runjob(id){
	$.get('jobrun?id='+id,function(data,status){
		//alert(data);
		if(data.indexOf("DOCTYPE") > 0){
			location.href = "login"
		}
	});
}

$('#jobs').on('click','.build_back img',function(e){
	var tag = $(this).attr("run");
	if(tag == 1){
		$('#temp').val($(this).attr('jobid'))
		$('#stopjob').modal('show');
	}else{
		runjob($(this).attr('jobid'))
		//alert($(this).attr('jobid'))
		$(this).attr('src',"./static/images/build_stop.png")
		$(this).attr("run","1")
		$(this).parent().css("background-size","0% 100%")
	}
});

function stopJob(){
	$.get('stopjob?id='+$('#temp').val(),function(data,status){
		if(data.toString().indexOf("DOCTYPE") > 0){
			window.location.href = "login"
		}
		//do what you want
	});
	$('#stopjob').modal('hide');
	$("img[jobid='"+$('#temp').val()+"'].run").attr("src","./static/images/run_build.png")
	$("img[jobid='"+$('#temp').val()+"'].run").attr("run","0")
	$("img[jobid='"+$('#temp').val()+"'].run").parent().css("background-size","100% 100%")
	//alert($('#temp').val());
}

function checkJobsStatus(ids){
//	$.ajax({
//		type:"post",
//		url:"jobstatus",
//		data:{
//			"ids":ids.toString()
//		},
//		async:false,
//		success:function(data,status){
//			checkData(data);
//		}
//	});
	if(xhr){
		xhr.abort();
	}
	xhr = $.post('jobstatus',{
		'ids':ids.toString()
		//'data':JSON.stringify(getGrid('jobs').currentPageData)
	},function(data,status){
		checkData(data);
		if(getJobIds(getGrid('jobs').currentPageData).toString() == getJobIds(data).toString()){
			setTimeout(function(){checkAllStatus();},5000);
		}
	});
}
function checkData(data){
	for(var i=0;i<data.length;i++){
		var gridJob = getGridJob(data[i].Id)
		if(gridJob){
			if(gridJob.HistoryId != data[i].HistoryId || gridJob.Status != data[i].Status || gridJob.Log != data[i].Log){
				getGrid('jobs').updateItem(data[i])
			}
			if(data[i].IsRunning){
				$('#repo'+data[i].Id+' .build_back').css("background-size",(data[i].Progress * 100).toFixed(2) + '%' + " 100%")
			}else{
				$('#repo'+data[i].Id+' .build_back').css("background-size","100% 100%")
			}
		}
	}
}

function getGridJob(id){
	var current = getGrid('jobs').currentPageData
	if(current){
		for(var i=0;i<current.length;i++){
			if(current[i].Id == id){
				return current[i]
			}
		}
	}
}

function getJobIds(data){
	var ids = new Array()
	if(data){
		for(var i=0;i<data.length;i++){
			ids.push(data[i].Id)
		}
	}
	return ids;
}

function checkAllStatus(){
	var page = getGrid('jobs').currentPageData
	
	checkJobsStatus(getJobIds(page))
}

function initTimer(){
	if(int){
		clearInterval(int)
	}
	int = setInterval('checkAllStatus()',2000)
}

function fullScreen(){
	alert(1231);
}

$('#jobs').on("click",".full_screen img",function(event){
	consoleDiv = $('#log'+$(this).attr("jobid"));
	consoleDiv.toggleClass("full_screen_div");
	
	if(consoleDiv.hasClass('full_screen_div')){
		$(document.body).css("overflow","hidden")
		consoleDiv.css('width',"100%");
		consoleDiv.css('height',$(window).height());
		$(this).parent().css("position","fixed")
		$(this).parent().css("right","0")
		$(this).parent().css("top","0")
		$(this).parent().css("z-index","9999999")
		
		
	}else{
		consoleDiv.css('width',"90%");
		consoleDiv.css('height',"450px");
		$(this).parent().css("position","static")
		$(document.body).css("overflow","scroll")
	}
	
});

function showShell(){
	$('#shell').css("display","block")
	$('#shell').css('visibility',"visible")
	$('#shell').css('height',"150px")
	$("#fortype").css("display","none");
	$("#forprofile").css("display","none");
	$("#forpom").css("display","none");
}

function hideShell(){
	$('#shell').css('visibility',"hidden")
	$('#shell').css('height',"0px")
	$('#shell').css("display","none");
	$("#fortype").css("display","block");
	$("#forprofile").css("display","block");
	$("#forpom").css("display","block");
}

function showEditShell(){
	$('#editshell').css("display","block")
	$('#editshell').css('visibility',"visible")
	$('#editshell').css('height',"150px")
	$("#editfortype").css("display","none");
	$("#editforprofile").css("display","none");
	$("#editforpom").css("display","none");
}

function hideEditShell(){
	$('#editshell').css('visibility',"hidden")
	$('#editshell').css('height',"0px")
	$('#editshell').css("display","none")
	$("#editfortype").css("display","block");
	$("#editforprofile").css("display","block");
	$("#editforpom").css("display","block");
}


$('.btype').change(function(){
	var val = $('.btype:checked').val();
	if(val == 2){
		showShell();
	}else{
		hideShell();
	}
});

$('.editbtype').change(function(){
	var val = $('.editbtype:checked').val();
	if(val == 2){
		showEditShell();
	}else{
		hideEditShell();
	}
});


$('#searchpro').on('keyup',function(){
	var searchQ = $(this).val();
	searchQ = searchQ.toLowerCase();
	var prolist = $('#prolist .proname')
	for(var i=0;i<prolist.length;i++){
		var proname = $(prolist[i]).html();
		proname = formatProjectPath(proname).toLowerCase()
		if(proname.indexOf(searchQ)<0){
			$(prolist[i]).parent().parent().css('display','none')
		}else{
			$(prolist[i]).parent().parent().css('display','block')
		}
	}
})




$('#autobuild').change(function(){
	var checked = $('#autobuild').prop("checked");
	if(checked){
		$('#jquerycron').css("display","block");
		var initValue = "42 3 * * 5";
		if($('#jquerycron').html()){
			initValue = $('#jquerycron').cron("value");
		}
		$('#jquerycron').html('');
		$('#jquerycron').cron({
			initial: initValue,
			useGentleSelect: true
		});
	}else{
		$('#jquerycron').css("display","none");
	}
});

$('#editautobuild').change(function(){
	var checked = $('#editautobuild').prop("checked");
	if(checked){
		$('#editjquerycron').css("display","block");
		var initValue = "42 3 * * 5";
		if($('#editjquerycron').html()){
			initValue = $('#editjquerycron').cron("value");
		}
		$('#editjquerycron').html('');
		$('#editjquerycron').cron({
			initial: initValue,
			useGentleSelect: true
		});
	}else{
		$('#editjquerycron').css("display","none");
	}
});








                                                       