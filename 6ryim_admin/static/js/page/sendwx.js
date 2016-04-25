(function($, window){

	var sender = $('#sender').dropdown();
	var msg_type = $('#msg_type').dropdown();
	var msg_source = $('#msg_source').dropdown();


	$('.js_reply_OK').click(function(){
		var content = $('#content').val();
		if(content.trim() == ""){
			alert('不允许发送空的内容');
			return;
		}

		$.ajax({
			url:'/request/receive',
			data:{
				content:content,
				openid:sender.getValue(),
				msgType:msg_type.getValue(),
				source:msg_source.getValue()
			},
			type:'POST',
			dataType:'json',
			success:function(data){
				if(data.status){
					alert('发送成功!');
					$('#content').val('');
				}else{
					alert(data.msg);
				}	
			}
		});
	});
})(jQuery, window)
