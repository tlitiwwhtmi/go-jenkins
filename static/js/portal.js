/**
 * Created by duanhao on 2016/4/15.
 */

$(function(){

    function initPortalData(){
        $.get("portaldata",function(data){
            if(data != 'error'){
                $('#jobcount').html(data.JobCount)
                $('#avgdura').html(formatDuration(data.AveTime))
                $('#failure').html((data.Faiure * 100).toFixed(2) + '%')
                $('#paccount').html(data.PacCount)
                $('#usercount').html(data.UserCount)
            }
        });
    }

    function initBuildDay(){
        $.get('buildsday',function(data){
            if(data != 'error'){
                var xaxisArr = new Array();
                var yaxisArr = new Array();
                for(var i=0;i<data.length;i++){
                    var obj = data[i];
                    xaxisArr.push(obj.StartDate)
                    yaxisArr.push(obj.Total)
                }
                $('#placeholder3xx3').highcharts({
                    chart:{
                        type:"spline"
                    },
                    title: {
                        text: '',
                        x: -20 //center
                    },
                    xAxis: {
                        categories: xaxisArr
                    },
                    yAxis: {
                        title: {
                            text: 'times'
                        },
                        plotLines: [{
                            value: 0,
                            width: 1,
                            color: '#808080'
                        }]
                    },
                    tooltip: {
                        valueSuffix: ''
                    },
                    legend: {
                        enabled:false
                    },
                    series: [{
                        name: 'builds',
                        data: yaxisArr
                    }],
                    credits:{
                        enabled:false
                    },
                    exporting:{
                        enabled:false
                    }
                });

            }
        })
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

    initPortalData();
    initBuildDay();
})