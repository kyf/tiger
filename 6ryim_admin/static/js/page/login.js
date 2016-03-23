(function($, window){

	var user = $('#user'), password = $('#password'), submitbt = $('.btnLogin');

	submitbt.click(function(){
		$.ajax({
			url:'/checklogin',
			data:{
				user:user.val(),
				password:password.val()
			},
			type:'POST',
			success:function(data, status, response){
				if(data == "success"){
					window.location.href = '/main';
				}else{
					alert('用户名或密码错误!');
				}	
			}
		});
	});
})(jQuery, window)
