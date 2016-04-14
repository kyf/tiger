(function($, window){
	var formtpl = [
		'<div class="dialog_bd_inner" style="margin:30px 60px 20px;">',
		    '<div class="frm_control_group">',
		        '<label class="frm_label">用户名</label>',
		        '<div class="frm_controls">',
		            '<span class="frm_input_box with_counter counter_in append count" style="padding:0px;">',
		                '<input type="text" class="frm_input username" value="{user}">',
		            '</span>',
		            '<span class="frm_msg fail js_usernamefail">请输入用户名</span>',
		        '</div>',
		    '</div>',
		    '<div class="frm_control_group">',
		        '<label class="frm_label">密码</label>',
		        '<div class="frm_controls">',
		            '<span class="frm_input_box with_counter counter_in append count" style="padding:0px;">',
		                '<input type="password" class="frm_input password">',
		            '</span>',
		            '<span class="frm_msg fail js_passwordfail">请输入密码</span>',
		        '</div>',
		    '</div>',
		    '<div class="frm_control_group">',
		        '<label class="frm_label">确认密码</label>',
		        '<div class="frm_controls">',
		            '<span class="frm_input_box with_counter counter_in append count" style="padding:0px;">',
		                '<input type="password" class="frm_input repassword">',
		            '</span>',
		            '<span class="frm_msg fail js_repasswordfail">两次输入密码不一致</span>',
		        '</div>',
		    '</div>',
		    '<div class="frm_control_group">',
		        '<label class="frm_label">openid</label>',
		        '<div class="frm_controls">',
		            '<span class="frm_input_box with_counter counter_in append count" style="padding:0px;">',
		                '<input type="text" class="frm_input openid" value="{openid}">',
		            '</span>',
		            '<span class="frm_msg fail js_openidfail">请输入微信openid</span>',
		        '</div>',
		    '</div>',
		'</div>'
	];

	var listtpl = [
			'<tr>',
				'<td class="table_cell nickname">',
					'<div class="nickname_inner">',
						'<img src="http://admin.6renyou.com/statics/socketchat/img/six-service.jpg">',
						'<p class="ncik_name">{user}</p>',
					'</div>',
				'</td>',
				'<td class="table_cell info">',
					'<div class="info_inner">',
						'<p class="wx_account">{opid}</p>',
					'</div>',
				'</td>',
				'<td class="table_cell opr"><div class="opr_inner">',
					'<a href="javascript:;" class="js_kf_record" data-id="{id}" data-openid="{opid}">客服记录</a>',
					'<a href="javascript:;" class="js_kf_edit" data-user="{user}" data-openid="{opid}" data-id="{id}">编辑</a>',
					'<a href="javascript:;" class="js_kf_del" data-id="{id}">删除</a>',
				'</div></td>',
			'</tr>'
		];

	var emptytpl = [
			'<tr>',
				'<td colspan="3" style="height:75px;line-height:75px;">还没有添加账户</td>',
			'</tr>'
		];

	$(document.body).on('click','.js_kf_edit', function(){
		var lastdata = {
			id:$(this).attr('data-id'),
			user:$(this).attr('data-user'),
			openid:$(this).attr('data-openid')
		};
		Addfn('编辑帐号', 'edit', lastdata);
	});

	$(document.body).on('click','.js_kf_del', function(){
		if(confirm("确认删除？")){
			$.ajax({
				url : '/admin/remove',
				dataType:'json',
				data:{
					id:$(this).attr('data-id')	
				},
				type:'POST',
				success:function(data){
					if(data.status){
						loadAccount();
					}else{
						alert(data.msg);
					}
				}
			});
		}	
	});


	var loadAccount = function(){
		$.ajax({
			url : '/admin/list',
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					if(!data.data)data.data = [];
					data = data.data;
					if(data.length == 0){
						$('#js_list').html(emptytpl.join(''));
					}else{
						var tmp = [];
						$.each(data, function(i, d){
							tmp.push(listtpl.join('').replaceTpl(d));		
						});
						$('#js_list').html(tmp.join(''))
					}
				}else{
					alert(data.msg);
				}
			}
		});
	};

	loadAccount();

	$('.js_kf_add').click(function(){
		Addfn('添加帐号', 'add');
	});
	var Addfn = function(title, action, lastdata){
		var config = {
			'Title' : title,
			'Button' : '确定',
			'Content' : formtpl.join('').replaceTpl(lastdata ? lastdata : {})
		};
		var dialog = $(Template.dialog.join('').replaceTpl(config));
		dialog.find('.pop_closed').click(function(){
			dialog.hide();
		});
		dialog.find('.submitbt').click(function(){
			dialog.find('.frm_msg').hide();

			var username = dialog.find('.username').val().trim();	
			var password = dialog.find('.password').val().trim();	
			var repassword = dialog.find('.repassword').val().trim();	
			var openid = dialog.find('.openid').val().trim();	

			if(username == ""){
				dialog.find('.js_usernamefail').show();
				return;
			}

			if(password == ""){
				dialog.find('.js_passwordfail').show();
				return;
			}

			if(repassword == ""){
				dialog.find('.js_repasswordfail').show();
				return;
			}

			if(password != repassword){
				dialog.find('.js_repasswordfail').show();
				return;
			}

			if(openid == ""){
				dialog.find('.js_openidfail').show();
				return;
			}

			var data = {
					user:username,
					pwd:password,
					openid:openid
				};

			if(lastdata){
				data.id = lastdata.id;
			}

			$.ajax({
				url : '/admin/' + action,
				data : data,
				dataType:'json',
				type:'POST',
				success:function(data){
					if(data.status){
						dialog.hide();
						loadAccount();
					}else{
						alert(data.msg);
					}					
				}
			});
		});
		dialog.appendTo(document.body);
	};

})(jQuery,window)
