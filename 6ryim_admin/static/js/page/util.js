
var getQueryParam = function(name){
	var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
	var r = window.location.search.substr(1).match(reg);
	if(r!=null)return  unescape(r[2]); return null;
};


String.prototype.replaceTpl = function(data){
	var reg = /{([^}]+)}/g;
	var result = this.replace(reg, function(main, group){
		return data[group] ? data[group] : '';
	})
	return result;
}

String.prototype.trim = function(){
	return this.replace(/^\s+|\s+$/g, '');
};



function ajaxBeforeSend(R){
	R.setRequestHeader("Connection", "keep-alive");
};

function ts2time(timestamp){
	var d = new Date(timestamp * 1000);    //根据时间戳生成的时间对象
	var date = (d.getFullYear()) + "-" + 
		(d.getMonth() + 1) + "-" +
		(d.getDate()) + " " + 
		(d.getHours()) + ":" + 
		(d.getMinutes()) + ":" + 
		(d.getSeconds());
	return date;
}

var MSG_TYPE_TEXT = 1, 
	MSG_TYPE_IMAGE= 2,
	MSG_TYPE_AUDIO = 3;

var MSG_SOURCE_WX = 1;
var MSG_SOURCE_IOS = 2;
var MSG_SOURCE_Android = 3;
var MSG_SOURCE_PC = 4;

function toggleSource(source){
	var result = "全部";
	switch(parseInt(source, 10)){
		case MSG_SOURCE_WX:
			result = "微信";
			break;
		case MSG_SOURCE_PC:
			result = "PC";
			break;
		default:
			result = "全部";
			break;
	}

	return result;
}
