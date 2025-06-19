
function getQueryString(name) {
  	let reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
  	let r = window.location.search.substr(1).match(reg);
  	if (r != null) {
      	return unescape(r[2]);
  	};
  	return null;
};

//设置cookie
function setCookie(cname, cvalue, exdays) {
	if ( cvalue != "" && cvalue != null ) {
	 	var d = new Date();
	 	d.setTime(d.getTime() + (exdays*24*60*60*1000));
	 	var expires = "expires="+d.toUTCString();
	 	document.cookie = cname + "=" + cvalue + "; Path=/; " + expires;
	}else{
		console.log("cvalue is null");
	}
}

//获取cookie
function getCookie(cname) {
 	var name = cname + "=";
	//console.log(document.cookie);
 	var ca = document.cookie.split(';');
 	for(var i=0; i<ca.length; i++) {
  		var c = ca[i];
  		while (c.charAt(0)==' ') c = c.substring(1);
  		if (c.indexOf(name) != -1) return c.substring(name.length, c.length);
 	}
 	return "";
}

//删除cookie 
function delCookie(name) {
  	document.cookie = name + '=; Max-Age=-99999999;';
}

	//设置默认集群
function SetDefaultCluster() {
	var clusterId = getQueryString("clusterId");
	if (clusterId == null) {
		clusterId = getCookie("clusterId")
		console.log(clusterId);
		if  (clusterId == "") {
			$.ajax({
			   url: "/rbac/cluster/MyClusterList",
			   type: "GET",
			   headers:{'X-Requested-With':'XMLHttpRequest'},
			   success: function (resp, status, xhr) {
					if (resp.count > 0) {
						var myClusterId = resp.data[0].cluster_id;
						if (myClusterId != 'all') {
							setCookie('clusterId',myClusterId,30);
						}
					}else{
						layer.msg('读取集群列表失败',{icon:2});
					}
				},
				error: function(xhr, status, error) {
    				layer.msg('获取集群列表失败',{icon:2});
  				}
		    });			
		}
	}
}

//function SetDefaultCluster() {
//	var clusterId = getQueryString("clusterId");
//	if (clusterId == null) {
//		clusterId = getCookie("clusterId")
//		if  (clusterId == "") {
//			$.get('/xkube/cluster/v1/List', function (resp) {
//				if (resp.count > 0) {
//					setCookie('clusterId',resp.data[0].cluster_id,30);
//				}else{
//					layer.msg('读取集群信息失败，或先添加k8s配置',{icon:2});
//				}
//		    });			
//		}
//	}
//}

// 版本比较使用示例
//console.log(compareVersions("1.2.3", "1.2.4")); // -1
//console.log(compareVersions("1.2.3", "1.2.3")); // 0
//console.log(compareVersions("1.2.4", "1.2.3")); // 1
function compareVersions(v1, v2) {
  const parts1 = v1.split('.').map(Number);
  const parts2 = v2.split('.').map(Number);
  
  for (let i = 0; i < Math.max(parts1.length, parts2.length); i++) {
    const num1 = parts1[i] || 0;
    const num2 = parts2[i] || 0;
    
    if (num1 > num2) return 1;
    if (num1 < num2) return -1;
  }
  
  return 0;
}


function GetCurrClusterId() {
	var clusterId = getQueryString("clusterId");
	if (clusterId == null) {
		clusterId = getCookie("clusterId")
	}
	$.get('/xkube/cluster/v1/List', function (resp) {
		$.each(resp.data,function(i,item){
			if (item.cluster_id == clusterId) {
				var html = '<option value="'+item.cluster_id+'" selected="">'+item.cluster_name+'</option>'
			}else{
				var html = '<option value="'+item.cluster_id+'">'+item.cluster_name+'</option>'	
			}
			$("#cluster_Id").append(html);
			setCookie(item.cluster_id,item.kube_version,30);
			//$('#cluster_Id').append(new Option(item.cluster_id,item.cluster_id));
		});
	   layui.use('form', function(){
	           var form = layui.form;
	           form.render();
	           //form.render('select', 'clusterId');
	   });		
    });		
}

function GetNamespace() {
	var clusterId = getQueryString("clusterId");
	if (clusterId == null) {
		clusterId = getCookie("clusterId")
	}
	if(clusterId != "") {
		$.get('/xkube/ns/v1/List?clusterId='+clusterId, function (resp) {
			$.each(resp.data,function(i,item){
				var html2 = '<option value="'+item.nameSpace+'">'+item.nameSpace+'</option>'
				$("#name_Space").append(html2);
			});
		    layui.use('form', function(){
		             var form = layui.form;
		            //form.render('select', 'name_Space');
					form.render('');
		    });					
	    });	
			
	}
}

function ExportRaw(filname, data) {
　　　　var urlObject = window.URL || window.webkitURL || window;
　　　　var export_blob = new Blob([data]);
　　　　var save_link = document.createElementNS("http://www.w3.org/1999/xhtml", "a")
　　　　save_link.href = urlObject.createObjectURL(export_blob);
　　　　save_link.download = filname;
　　　　save_link.click();
}