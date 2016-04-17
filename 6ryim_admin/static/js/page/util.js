
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


var MSG_TYPE_TEXT = 1, 
	MSG_TYPE_IMAGE= 2,
	MSG_TYPE_AUDIO = 3;
