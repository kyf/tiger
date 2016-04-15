(function($, window){

	var sender = $('#sender').dropdown();

	$('.js_reply_OK').click(function(){
		var content = $('#content');
		if(content.trim() == ""){
			alert('不允许发送空的内容');
			return;
		}

		$.ajax({
			url:'/handleReceive',
			data:{
				
			},
			type:'POST',
			dataType:'json',
			success:function(data){
			
			}
		});
	});
})(jQuery, window)
