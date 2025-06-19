var dashboard = {};

//获取一周前的日期
function GetStartTime() {
　　 //let date =  new Date(new Date().getTime()- 6 * 24 * 60* 60 * 1000); //获取一周前的日期
	let date =  new Date();	//当前时间
    let year = date.getFullYear(); 　//返回的是年份
    let month = date.getMonth() + 1; //返回的是月份 记得+1 才是当前月 （从0开始 0代表1月） 
	if (month < 10) {
		month = "0" + month;
	}
　  let dates = date.getDate()　　　　//返回的是日 date上面定义了 这里命名为 dates
	if (dates < 10) {
		dates = "0" + dates;
	}
      //let day = date.getDay();　　　　　//返回的是周几
    // console.log(year.toString().padStart(4,'0')+month+day);
    //console.log(year.toString()+'-' + month.toString()+'-' +dates.toString()); //输出时将年月日转换为 字符串 形式，不然三个数字会进行相加，导致值错误
	var strdate = year.toString() + "-" + month.toString() + "-" + dates.toString()+"T00:00:00";
	return strdate;
}

function GetLastWeek() {
　　 let date =  new Date(new Date().getTime()- 6 * 24 * 60* 60 * 1000); //获取一周前的日期
    let year = date.getFullYear(); 　//返回的是年份
    let month = date.getMonth() + 1; //返回的是月份 记得+1 才是当前月 （从0开始 0代表1月） 
	if (month < 10) {
		month = "0" + month;
	}
　  let dates = date.getDate()　　　　//返回的是日 date上面定义了 这里命名为 dates
	if (dates < 10) {
		dates = "0" + dates;
	}
      //let day = date.getDay();　　　　　//返回的是周几
    // console.log(year.toString().padStart(4,'0')+month+day);
    //console.log(year.toString()+'-' + month.toString()+'-' +dates.toString()); //输出时将年月日转换为 字符串 形式，不然三个数字会进行相加，导致值错误
	var strdate = year.toString() + "-" + month.toString() + "-" + dates.toString()+"T00:00:00";
	return strdate;
}


Date.prototype.format = function(format) {
    /*
     * eg:format="YYYY-MM-dd hh:mm:ss";
     * new Date().format("yyyy.MM.dd");

     */
    var o = {
        "M+" :this.getMonth() + 1, // month
        "d+" :this.getDate(), // day
        "h+" :this.getHours(), // hour
        "m+" :this.getMinutes(), // minute
        "s+" :this.getSeconds(), // second
        "q+" :Math.floor((this.getMonth() + 3) / 3), // quarter
        "S" :this.getMilliseconds()
    // millisecond
    }
    if (/(y+)/.test(format)) {
        format = format.replace(RegExp.$1, (this.getFullYear() + "")
                .substr(4 - RegExp.$1.length));
    }
    for ( var k in o) {
        if (new RegExp("(" + k + ")").test(format)) {
            format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k]
                    : ("00" + o[k]).substr(("" + o[k]).length));
        }
    }
    return format;
}

//pod资源使用曲线
dashboard.DynamicLine = function () {
	var runInterval = 2000; //1s
	//var ndayz = new Date(new Date().toLocaleDateString()).format("yyyy-MM-ddThh:mm:ss");
	var myDtLineChart = echarts.init(document.getElementById('MemDynamicLine'),'walden');
	lineoption = null;	
	var dtime = [];
	var datas = [];
	myDtLineChart.setOption(lineoption = {
		backgroundColor: '',
		title:{
			text:'内存使用[MBi]',
			top: '20',
			left:'center'
		},	
		tooltip: {
			trigger: 'axis'
		},
		grid: {
			top: '60',
			left: '9%',
			right: '3%',
			bottom: '50',
			containLabel: false
		},			
		xAxis: {
			type: 'category',
			//type: 'time',
			boundaryGap: false,
			data: dtime,
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		yAxis: {
			type: 'value',
			boundaryGap: [0, '100%'],
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		series: [{
			data: datas,
			type: 'line'
		}]
	});
	var myDtLineChart2 = echarts.init(document.getElementById('CpuDynamicLine'),'walden');
	lineoption2 = null;	
	var datac = [];
	myDtLineChart2.setOption(lineoption2 = {
		backgroundColor: '',
		title:{
			text:'CPU使用[毫核m]',
			top: '20',
			left:'center'
		},	
		tooltip: {
			trigger: 'axis'
		},
		grid: {
			top: '60',
			left: '9%',
			right: '3%',
			bottom: '50',
			containLabel: false
		},			
		xAxis: {
			type: 'category',
			//type: 'time',
			boundaryGap: false,
			data: dtime,
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		yAxis: {
			type: 'value',
			boundaryGap: [0, '100%'],
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		series: [{
			data: datac,
			type: 'line'
		}]
	});	
	setInterval(function () {
		//var nowtimestamp = new Date().getTime();
		//var sTime = nowtimestamp - runInterval
		//{name: "16:12",value: 22746}
		$.get('/xkube/metrics/PodUsage'+ location.search, function (data) {
         	dtime.push(data.time);
		 	datas.push(data.mem);
			datac.push(data.cpu);
			myDtLineChart.setOption(lineoption,true);	
			myDtLineChart2.setOption(lineoption2,true);
		});
	},runInterval);		
}

//节点资源使用曲线
dashboard.NodeDynamicLine = function () {
	var runInterval = 2000; //1s
	//var ndayz = new Date(new Date().toLocaleDateString()).format("yyyy-MM-ddThh:mm:ss");
	var myDtLineChart = echarts.init(document.getElementById('NodeMemDynamicLine'),'walden');
	lineoption = null;	
	var dtime = [];
	var datas = [];
	myDtLineChart.setOption(lineoption = {
		backgroundColor: '',
		title:{
			text:'内存使用[GBi]',
			top: '20',
			left:'center'
		},	
		tooltip: {
			trigger: 'axis'
		},
		grid: {
			top: '60',
			left: '9%',
			right: '3%',
			bottom: '50',
			containLabel: false
		},			
		xAxis: {
			type: 'category',
			//type: 'time',
			boundaryGap: false,
			data: dtime,
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		yAxis: {
			type: 'value',
			boundaryGap: [0, '100%'],
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		series: [{
			data: datas,
			type: 'line'
		}]
	});
	var myDtLineChart2 = echarts.init(document.getElementById('NodeCpuDynamicLine'),'walden');
	lineoption2 = null;	
	var datac = [];
	myDtLineChart2.setOption(lineoption2 = {
		backgroundColor: '',
		title:{
			text:'CPU使用[核]',
			top: '20',
			left:'center'
		},	
		tooltip: {
			trigger: 'axis'
		},
		grid: {
			top: '60',
			left: '7%',
			right: '3%',
			bottom: '50',
			containLabel: false
		},			
		xAxis: {
			type: 'category',
			//type: 'time',
			boundaryGap: false,
			data: dtime,
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		yAxis: {
			type: 'value',
			boundaryGap: [0, '100%'],
			splitLine: {
				show: true,
				lineStyle:{
					opacity:0.2
				}
			}				
		},
		series: [{
			data: datac,
			type: 'line'
		}]
	});	
	setInterval(function () {
		//var nowtimestamp = new Date().getTime();
		//var sTime = nowtimestamp - runInterval
		//{name: "16:12",value: 22746}
		$.get('/xkube/metrics/NodeUsage'+ location.search, function (data) {
         	dtime.push(data.time);
		 	datas.push(data.mem);
			datac.push(data.cpu);
			myDtLineChart.setOption(lineoption,true);	
			myDtLineChart2.setOption(lineoption2,true);
		});
	},runInterval);		
}

