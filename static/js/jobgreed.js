/**
 * Created by duanhao on 2016/1/19.
 */

var grids = {};

var _hisPageindex = 1;

function registerGrid(grid){
    if(!grids[grid.attr('id')]){
        grids[grid.attr('id')] = grid;
    }
}

function getGrid(id){
    return grids[id];
}

$.fn.jobgreed = function (projectId) {
    registerGrid(this);
    var _self = this;
    _self.pageNum = 1;
    _self.perPage = 10;
    _self.currentPageData = [];
    _self.pageCount = 0;
	_self.projectId;
	_self.branchName = '';
	_self.branches;
    _self.keyword = '';
    //init with the varibles

//    if (typeof options == 'string') {
//        if(options =='duanhao'){
//            alert(_self.url);
//            return;
//        }
//        _self.url = options;
//    } else {
//        _self.url = options.url;
//        if (options.perPage) {
//            _self.perPage = options.perPage;
//            //do with the varibles
//        }
//    }

	_self.projectId = projectId
	
	_self.setProjectId = function(id){
		_self.projectId  = id;
	}

    _self.loadData = function () {
		$('#loadingjob').css('display',"block")
		_self.css("display","none")
        $.get("listjob?pid="+_self.projectId+"&page="+_self.pageNum+"&size="+_self.perPage+"&branch="+_self.branchName+"&keyword="+_self.keyword,function(data,status){
            if(data.toString().indexOf("DOCTYPE") > 0){
                window.location.href = "login"
            }
            if(data.ProjectId == _self.projectId){
				_self.handleData(data,status);
				checkAllStatus()
				//initTimer();
			}
        });
    }

    _self.handleData = function (data, status) {
        if (true) {  //judge the status

            _self.append('<div id="jobgreedheader"></div><div id="jobgreedbody"></div><div id="jobgreedfooter"></div>')

            _self.getTheHead(data)
            _self.getTheBody(data);
            _self.getTheFooter(data);
			$('#loadingjob').css('display',"none")
			_self.css("display","block")
        }
    }

    _self.getTheHead = function (data) {
		if(_self.currentPageData == data.Jobs){
			return;
		}
		_self.find('#jobgreedheader').html('')
//        if(_self.find('#jobgreedheader').html() != ''){
//            $('#job_count').html('Total: ' + data.Total);
//            return;
//        }
		_self.branches = data.Branches;
        var htmlstr = '';
        htmlstr += '<div style="width: 100%;height: 50px;padding-top: 10px">' +
            '<div style="width: 150px;float: left;margin-left: 10px;height: 30px;line-height: 30px;font-weight: bold;color: #898989" id="job_count">Total: ' + data.Total + ' </div>' +
            '<div style="width: 130px;margin-right: 10px;float: right;background-color: #44A662;">';
        if(_self.projectId){
            htmlstr += '<a class="more_button" style="text-decoration: none;color: white" data-toggle="modal" data-target="#myModal">+ Add New Job</a>'
        }else{
            htmlstr += '<a class="more_button" style="text-decoration: none;color: white" data-toggle="modal" data-target="#warnjob">+ Add New Job</a>'
        }
        if(data.Keyword){
            _self.keyword = data.Keyword
        }
        htmlstr += '</div>' +
            '<div style="width: 150px;float: right;margin-right: 10px">' +
            '<input id="keyword" type="search" value="'+_self.keyword+'" class="form-control" placeholder="Search" style="margin: auto;width:100%;height: 30px;">' +
            '</div>' +
//            '<div class="more_btn_div" style="width: 130px;float: right;margin-right: 10px">' +
//            '<a class="more_button" style="text-decoration: none;">All Type<span class="glyphicon glyphicon-chevron-down dropdown-arrow"></span> </a>' +
//            '<ul class="more-dropdown">' +
//            '<li>' +
//            '<a href="/tlitiwwhtmi/java-gitlab-api/settings" class="ember-view">All Type</a>' +
//            '</li>' +
//            '<li>' +
//            '<a href="/tlitiwwhtmi/java-gitlab-api/settings" class="ember-view">auto build</a>' +
//            '</li>' +
//            '<li>' +
//            '<a href="/tlitiwwhtmi/java-gitlab-api/requests" class="ember-view">manual build</a>' +
//            '</li>' +
//            '</ul>' +
//            '</div>' +
            '<div class="more_btn_div" style="width: 130px;float: right;margin-right: 10px">' +
            '<a class="more_button" style="text-decoration: none;">';
			if(_self.branchName == ''){
				htmlstr += "All Branches"
			}else{
				htmlstr += _self.branchName
			}
			htmlstr += '<span class="glyphicon glyphicon-chevron-down dropdown-arrow"></span></a>' +
            '<ul class="more-dropdown jobbranch">' +
            '<li>' +
            '<a href="javascript:changeBranch()" class="ember-view">All Branches</a>' +
            '</li>';
		
		if(data.Branches){
			for(var i=0;i<data.Branches.length;i++){
				branch = data.Branches[i]
				htmlstr += '<li>' +
	           			   '<a href="javascript:changeBranch(\''+branch.name+'\')" class="ember-view">'+branch.name+'</a>' +
	           			   '</li>';
			}
		}	
		
			
        htmlstr += '</ul>' +
            '</div>' +
            '</div>';
        _self.find('#jobgreedheader').append(htmlstr);
    }


    _self.getTheBody = function (data) {
        if(_self.currentPageData == data.Jobs){
            return ;
        }
        _self.currentPageData = data.Jobs;
        _self.find('#jobgreedbody').html('');
		if(!_self.currentPageData){
            _self.find('#jobgreedbody').append('<div style="font-size: 25px;line-height: 100px;height: 100px;color: rgba(134, 122, 122, 0.95);text-align: center;">There has no job here</div>');
			return;
		}
        var htmlstr = '';
        for (var i = 0; i < _self.currentPageData.length; i++) {
            var job = _self.currentPageData[i];
            htmlstr += '<div class="repoblock animated fadeIn" id="repo'+job.Id+'">';
			htmlstr += _self.generateItemHtml(job);
            htmlstr += '</div>';
			//htmlstr += '<div class="console">'
			//htmlstr += '<div class="log_out_div" id="log'+job.Id+'">'
			//htmlstr += _self.generateLog(job);
			//htmlstr += '</div>'
			//htmlstr += '<div class="full_screen"><img jobid="'+job.Id+'" src="./static/images/fullscreen.png"></div>'
			//htmlstr += '<div style="clear:both">'
			//htmlstr += '</div>'
			
        }
        _self.find('#jobgreedbody').append(htmlstr);
    }

    _self.turunToPage = function(pageIndex){
		//alert(pageIndex)
        _self.pageNum = pageIndex;
        _self.loadData();
    }

    _self.getData = function(id){
        for(var i=0;i<_self.currentPageData.length;i++){
            if(self.currentPageData[i].id == id){
                return _self.currentPageData[i];
            }
        }
    }

    _self.updateSpecificData = function(id){}

    _self.getTheFooter = function (data) {
        var pageCount = Math.ceil(data.Total/_self.perPage);
//        if(data.Total%_self.perPage != 0){
//            pageCount++;
//        }
		if(pageCount ==0){
			pageCount = 1;
		}
//        if(_self.pageCount == pageCount){
//            return;
//        }
        _self.pageCount = pageCount;
        _self.find('#jobgreedfooter').html("");
        var htmlstr = '';
        htmlstr += '<ul class="pagination">' +
            '<li><a href="javascript:getGrid(\'jobs\').turunToPage(1)">&laquo;</a></li>';
        for(var i=1;i<=_self.pageCount;i++) {
            if (i == _self.pageNum) {
                htmlstr += '<li class="active"><a href="javascript:getGrid(\'jobs\').turunToPage('+i+')">' + i + '</a></li>';
            } else {
                htmlstr += '<li><a href="javascript:getGrid(\'jobs\').turunToPage('+i+')">' + i + '</a></li>';
            }

        }
        htmlstr += '<li><a href="javascript:getGrid(\'jobs\').turunToPage('+pageCount+')">&raquo;</a></li>' +
            '</ul>';
        _self.find('#jobgreedfooter').append(htmlstr)
    }
	
	_self.generateItemHtml = function(job){
		var htmlstr = '';
		if(job.Status == "FAILURE"){
				htmlstr += '<div class="lasterror"></div>'
		}
		else if(job.Status == 'SUCCESS'){
			htmlstr += '<div class="lastsuccess"></div>'
		}else{
			htmlstr += '<div class="lastever"></div>'
		}
        htmlstr += '<div class="jobcontent">' +
            '<div class="jobmain">' +
            '<div style="overflow: hidden;height:40%">' +
            '<div class="job_name">' +
            '<h3 style="display: inline-block">'+job.JobName+'</h3>' +
            '<img class="job_edit" jobid="'+job.Id+'" src="./static/images/jobmodify.png">' +
            '</div>' +
            '<label style="margin-left: 20px;font-size: 16px;font-weight: 400;color: #44A662">'+job.BranchName+'</label>' +
            '</div>' +
            '<div>' +
            '<p class="modify-author">' +
            //'<img src="./static/images/jobauth.png" style="width: 20px;height: 20px;margin-right: 5px;margin-left: 15px">' +
            '<span style="color: #898989;margin-right: 15px">'+job.Modifier+'</span>' +
            '</p>' +
            '</div>' +
            '<div>' +
            '<p class="modify-author">' +
            //'<img src="./static/images/jobauth.png" style="width: 20px;height: 20px;margin-right: 5px;margin-left: 15px">' +
            '<span style="color: #898989">'+formattime(job.ModifyTime)+'</span>' +
            '</p>' +
            '</div>' +
            '</div>' +
            '<div class="lastinfo">'+
			'<div>' +
               '<img style="width: 20px;height:18px;" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyMCAyMCI+PHBhdGggZmlsbD0iI0E3QUVBRSIgZD0iTTE2LjcgMi41SDMuM2MtLjMgMC0uNS4yLS41LjV2MTRjMCAuMy4yLjUuNS41aDEzLjRjLjMgMCAuNS0uMi41LS41VjNjMC0uMy0uMi0uNS0uNS0uNXptLS41Ljl2My4xSDMuOFYzLjRoMTIuNHpNMy44IDE2LjZ2LTloMTIuNXY5SDMuOHoiLz48cGF0aCBmaWxsPSIjQTdBRUFFIiBkPSJNOC43IDEzLjRoLS40Yy4yLS4yLjQtLjQuNS0uNi4yLS4yLjMtLjQuNS0uNmwuMy0uNmMuMS0uMi4xLS40LjEtLjYgMC0uMiAwLS40LS4xLS42LS4xLS4yLS4yLS4zLS4zLS40LS4xIDAtLjItLjEtLjQtLjJzLS40LS4xLS42LS4xYy0uMyAwLS42LjEtLjguMi0uMy4xLS41LjMtLjcuNWwuNS42Yy4xLS4xLjItLjIuNC0uMy4xLS4xLjMtLjEuNC0uMS4yIDAgLjQuMS41LjIuMS4xLjIuMy4yLjVzMCAuMy0uMS41LS4yLjQtLjQuNmMtLjIuMi0uNC40LS42LjdsLS44Ljh2LjZIMTB2LS45aC0uOWMtLjEtLjItLjMtLjItLjQtLjJ6TTEwLjYgMTAuN2gyYy0uMi4zLS40LjYtLjUuOC0uMS4zLS4zLjYtLjQuOC0uMS4zLS4yLjYtLjIuOSAwIC4zLS4xLjctLjEgMWgxYzAtLjQgMC0uOC4xLTEuMiAwLS4zLjEtLjcuMi0uOS4xLS4zLjItLjYuNC0uOWwuNi0uOXYtLjVoLTMuMXYuOXoiLz48L3N2Zz4=">' +
               '<span style="color: #898989;margin-right: 10px">Last build:</span><span style="color: #898989">'+formattime(job.LastBuild)+'</span>' +
               '</div>';
		if(job.Status == ""){
			htmlstr += '<div>' +
               '<img style="width: 20px;height:18px;" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9Ii01MDYuOCAzNzcuNSAxNy41IDE5LjUiPjxzdHlsZT4uc3Qwe2ZpbGw6I0E1QUNBRDt9PC9zdHlsZT48cGF0aCBjbGFzcz0ic3QwIiBkPSJNLTQ5OS40IDM5N2MtNC4xIDAtNy40LTMuMy03LjQtNy40czMuMy03LjQgNy40LTcuNCA3LjQgMy4zIDcuNCA3LjQtMy40IDcuNC03LjQgNy40em0wLTEzLjRjLTMuMyAwLTYuMSAyLjctNi4xIDYuMXMyLjcgNi4xIDYuMSA2LjEgNi4xLTIuNyA2LjEtNi4xLTIuOC02LjEtNi4xLTYuMXoiLz48cGF0aCBjbGFzcz0ic3QwIiBkPSJNLTQ5Ny4xIDM5MmMtLjEgMC0uMiAwLS40LS4xbC0yLjYtMS43Yy0uMi0uMS0uMy0uMy0uMy0uNXYtNC4xYzAtLjQuMy0uNi42LS42cy42LjMuNi42djMuOGwyLjMgMS41Yy4zLjIuNC42LjIuOS4xLjEtLjIuMi0uNC4yek0tNDk5LjQgMzgxLjljLTEuMiAwLTIuMi0xLTIuMi0yLjJzMS0yLjIgMi4yLTIuMmMxLjIgMCAyLjIgMSAyLjIgMi4ycy0xIDIuMi0yLjIgMi4yem0wLTMuMWMtLjUgMC0uOS40LS45LjlzLjQuOS45LjkuOS0uNC45LS45LS40LS45LS45LS45ek0tNDkxLjMgMzg3LjZjLS4zIDAtLjctLjItLjgtLjVsLTEuMS0yLjFjLS4yLS41LS4xLTEgLjQtMS4zbDEuMS0uNmMuNS0uMiAxLS4xIDEuMy40bDEuMSAyLjFjLjIuNS4xIDEtLjQgMS4zbC0xLjEuNmMtLjIuMS0uNC4xLS41LjF6bS0uNy0zbC44IDEuNS40LS4yLS44LTEuNS0uNC4yeiIvPjwvc3ZnPg==">' +
               '<span style="color: #898989;margin-right: 10px">Duration:</span><span style="color: #898989">--</span>' +
               '</div>';
		}else{
			htmlstr += '<div>' +
               '<img style="width: 20px;height:18px;" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9Ii01MDYuOCAzNzcuNSAxNy41IDE5LjUiPjxzdHlsZT4uc3Qwe2ZpbGw6I0E1QUNBRDt9PC9zdHlsZT48cGF0aCBjbGFzcz0ic3QwIiBkPSJNLTQ5OS40IDM5N2MtNC4xIDAtNy40LTMuMy03LjQtNy40czMuMy03LjQgNy40LTcuNCA3LjQgMy4zIDcuNCA3LjQtMy40IDcuNC03LjQgNy40em0wLTEzLjRjLTMuMyAwLTYuMSAyLjctNi4xIDYuMXMyLjcgNi4xIDYuMSA2LjEgNi4xLTIuNyA2LjEtNi4xLTIuOC02LjEtNi4xLTYuMXoiLz48cGF0aCBjbGFzcz0ic3QwIiBkPSJNLTQ5Ny4xIDM5MmMtLjEgMC0uMiAwLS40LS4xbC0yLjYtMS43Yy0uMi0uMS0uMy0uMy0uMy0uNXYtNC4xYzAtLjQuMy0uNi42LS42cy42LjMuNi42djMuOGwyLjMgMS41Yy4zLjIuNC42LjIuOS4xLjEtLjIuMi0uNC4yek0tNDk5LjQgMzgxLjljLTEuMiAwLTIuMi0xLTIuMi0yLjJzMS0yLjIgMi4yLTIuMmMxLjIgMCAyLjIgMSAyLjIgMi4ycy0xIDIuMi0yLjIgMi4yem0wLTMuMWMtLjUgMC0uOS40LS45LjlzLjQuOS45LjkuOS0uNC45LS45LS40LS45LS45LS45ek0tNDkxLjMgMzg3LjZjLS4zIDAtLjctLjItLjgtLjVsLTEuMS0yLjFjLS4yLS41LS4xLTEgLjQtMS4zbDEuMS0uNmMuNS0uMiAxLS4xIDEuMy40bDEuMSAyLjFjLjIuNS4xIDEtLjQgMS4zbC0xLjEuNmMtLjIuMS0uNC4xLS41LjF6bS0uNy0zbC44IDEuNS40LS4yLS44LTEuNS0uNC4yeiIvPjwvc3ZnPg==">' +
               '<span style="color: #898989;margin-right: 10px">Duration:</span><span style="color: #898989">'+formatDuration(job.Duration)+'</span>' +
               '</div>';
		}
        if(job.Version){
            htmlstr += '<div>' +
                '<img style="width: 20px;height:18px;" src="/static/images/git.ico">' +
                '<a href="http://git01.dds.com/'+formatProjectPath($('.proactive .proname').html())+'/commit/'+job.Version+'" target="_blank"><span style="color: #898989">'+job.CommitMessage+"("+job.CommitAuthor+")"+'</span></a>' +
                '</div>';
        }

        htmlstr += '</div>' +
            '<div class="joboperate">' +
            '<div style="width: 130px;margin: auto">' +
            '<a class="more_button toglog" style="text-decoration: none;" href="logout?jobid='+job.Id+'" target="_blank" jobid="'+job.Id+'">View log</a>' +
            '</div>' +
            '<div class="more_btn_div" style="width: 130px;margin: auto;margin-top: 10px">' +
            '<a class="more_button" style="text-decoration: none;">More<span class="glyphicon glyphicon-chevron-down dropdown-arrow"></span></a>' +
            '<ul class="more-dropdown">' +
            '<li>' +
            '<a href="javascript:showHistory('+job.Id+',\''+job.JobName+'\')" class="ember-view">Histories</a>' +
            '</li>' +
            //'<li>' +
            //'<a href="#" class="ember-view">Requests</a>' +
            //'</li>' +
            '<li>' +
            '<a href="javascript:deleteJob('+job.Id+')" class="ember-view">Delete</a>' +
            '</li>' +
            '</ul>' +
            '</div>' +
            '</div>' +
            '<div class="jobrun">' +
            '<div class="build_back">';
		if(job.IsRunning){
			htmlstr += '<img class="run" jobid="'+job.Id+'" run="1" src="./static/images/build_stop.png" style="width: 48px;height: 48px;margin: auto;background-color: white;border-radius: 100%;cursor: pointer">';
		}else{
			htmlstr += '<img class="run" jobid="'+job.Id+'" run="0" src="./static/images/run_build.png" style="width: 48px;height: 48px;margin: auto;background-color: white;border-radius: 100%;cursor: pointer">';
		} 
        htmlstr += '</div>' +
            '</div>' +
            '</div>';
        if(job.AutoBuild){
            htmlstr += '<div class="repotagtimes">auto</div>'
        }

		return htmlstr;	
	}
	
	_self.generateLog = function(job){
		var log = job.Log
		if(job.Log.indexOf('[PostBuildScript]') != -1){
			log = log.slice(0,job.Log.indexOf('[PostBuildScript]'))
		}
		log = log.replace(/\n/g,"<p>");
		var htmlstr = '';
		htmlstr += '<div class="log_out"><p>';
		htmlstr += log;
		htmlstr += '</p>';
		if(log){
			if(job.Log.indexOf('[PostBuildScript]') != -1 && job.Status != ''){
				htmlstr += '<p>Finished:'+job.Status+'</p>'
			}else{
				htmlstr += '<p><img src="./static/images/loading.gif" style="height:24px" /></p>'
			}
		}
		htmlstr += '</div>'
		return htmlstr;
	}
	
	_self.updateItem =function(data){
		for(var i=0;i<_self.currentPageData.length;i++){
			if(_self.currentPageData[i].Id == data.Id){
				_self.currentPageData[i] = data;
				break;
			}
		}
		var htmlstr = _self.generateItemHtml(data);
		$("#repo"+data.Id).html(_self.generateItemHtml(data));
		//$("#log"+data.Id).html(_self.generateLog(data));
		//$("#log"+data.Id+" .log_out").scrollTop( $("#log"+data.Id+" .log_out")[0].scrollHeight );
	}

    _self.removeItem = function(id){
        for(var i=0;i<_self.currentPageData.length;i++){
            if(_self.currentPageData[i].Id == id){
                _self.currentPageData.removeArr(i);
                return;
            }
        }
    }

    //_self.loadData();
}

function formatDuration(duration){
    var timestr = "";
    var sec = 0;
    var min = 0;
    var hour = 0;
    sec = Math.floor(duration / 1000)
    if(Math.floor(sec/3600) > 0) {
        hour = Math.floor(sec / 3600)
        sec = sec % 3600
    }
    if(Math.floor(sec/60) > 0) {
        min = Math.floor(sec / 60)
        sec = sec % 60
    }
    if(hour != 0) {
        timestr += hour + "h"
    }
    if (min != 0) {
        timestr += min + "m"
    }
    timestr += sec + "s"
    return timestr
}


function formattime(str){
    var strs1 = str.split("T");
    var strs2 = strs1[1].split("+")
    return strs1[0]+" "+strs2[0];
}


function formatVersion(str){
    return str ? str.substring(0,7) : ""
}


function formatProjectPath(str){
    while(str.indexOf(" ")>0){
        str = str.replace(" ","");
    }
    return str;
}


function deleteJob(id){
    $('#deljobid').val(id);
    $('#deljob').modal('show')
}

function sureDelJob(){
    var id = $('#deljobid').val();
    if(id != ""){
        $("#repo"+id).animate({opacity:"0"},function(){
            $("#repo"+id).animate({height:0},function(){
                $("#repo"+id).css('display',"none")
            });
        });
        getGrid('jobs').removeItem(id);
        $.get("jobdel?id="+id,function(data,status){
            if(data.toString().indexOf("DOCTYPE") > 0){
                window.location.href = "login"
            }
            if(data != 'success'){
                alert(data);
            }
        });
    }
    $('#deljob').modal('hide')

}

/*function deleteJob(id){
    $("#repo"+id).animate({opacity:"0"},function(){
        $("#repo"+id).animate({height:0},function(){
            $("#repo"+id).css('display',"none")
        });
    });
    getGrid('jobs').removeItem(id);
    //$("#repo"+id).remove();
    $.get("jobdel?id="+id,function(data,status){
        if(data != 'success'){
            alert(data);
        }
    });
}*/

function showHistory(jobid,jobName){
    $('#jobHistory').html(jobName+" histories")
    $('#jobhis').modal('show')
    $('#loadinghisbd').css('display',"block")
    $('#hiscontent').css('display',"none")
    $.get("lisyhistories?jobid="+jobid+"&page="+_hisPageindex+"&cid=0",function(data){
        if(data.toString().indexOf("DOCTYPE") > 0){
            window.location.href = "login"
        }
        var htmlstr = handlehisdata(data,_hisPageindex);
        if(htmlstr != ''){
            $('#hiscontent').html('');
            $('#hiscontent').append(htmlstr)
        }
        $('#loadinghisbd').css('display',"none")
        $('#hiscontent').css('display',"block")
    })
}

function handlehisdata(data,index){
    if(data == "error"){
        alert(data)
        return;
    }
    var htmlstr = "";
    if(data){
        for(var i=0;i<data.length;i++){
            htmlstr += '<div class="hisblock">'+
                '<div style="float: left;width: 5%">'+
                '<label style="line-height: 50px;">#'+((index-1)*10+i+1)+'</label>'+
                '</div>'+
                '<div style="float: left;width: 25%">'+
                ' <label style="line-height: 50px;">'+formattime(data[i].StartTime)+'</label>'+
                '</div>'+
                '<div style="float: left;width: 10%">'+
                '<label style="line-height: 50px;">'+data[i].Status+'</label>'+
                '</div>'+
                '<div style="float: left;width: 10%">'+
                '<label style="line-height: 50px;">'+data[i].BuildExecutor+'</label>'+
                '</div>'+
                '<div style="float: left;width: 10%">'+
                '<label style="line-height: 50px;">'+formatDuration(data[i].Duration)+'</label>'+
                '</div>'+
                '<div style="float: left;width: 30%">'+
                '<a href="http://git01.dds.com/'+formatProjectPath($('.proactive .proname').html())+'/commit/'+data[i].Version+'" target="_blank" style="line-height: 50px;"><nobr>'+formateMessage(data[i].Message)+'-'+data[i].CommitAuthor+'</nobr></a>'+
                ' </div>'+
                '<div style="float: left;width: 10%">'+
                '<a href="spelog?id='+data[i].Id+'" target="_blank" style="line-height: 50px;">查看日志</a>'+
                '</div>'+
                '</div>'
        }
    }else{
        htmlstr += '<div style="font-size: 20px;color: rgba(134, 122, 122, 0.95);">The job has no histories to show</div>'
    }

    return htmlstr;
}

Array.prototype.removeArr = function (index) {
    if (isNaN(index) || index>= this.length) { return false; }
    this.splice(index, 1);
}


function formateMessage(msg){
    if(msg.length > 16){
        return msg.substring(0,15)+"..."
    }else{
        return msg;
    }
}

function changeBranch(name){
    if(!name){
        name=""
    }
    var grid = getGrid('jobs')
    if(grid.branchName == name){
        return;
    }
    grid.branchName = name;
    grid.pageNum = 1;
    grid.loadData();
}

$('#joblist').on('keyup',"#keyword", function(event) {
    if (event.keyCode == 13) {
        var grid = getGrid('jobs');
        if(grid && grid.keyword != $('#keyword').val()){
            grid.keyword = $('#keyword').val();
            grid.pageNum = 1;
            grid.loadData();
        }
    }
});



