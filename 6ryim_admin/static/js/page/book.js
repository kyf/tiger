(function($, window){
	var openid = getQueryParam("openid");	
	var opid = getQueryParam("opid");	

	if(openid == null || openid == '' || opid == null || opid == ''){
		alert('参数错误');
		return;
	}

	$('#iframer').attr('src', "http://admin.6renyou.com/weixin_plugin/orderDetail?action=WX&openid=" + openid + "&op=" + opid);
	$('#iframer').show();
})(jQuery, window)
