(function($, window){
	var formtpl = [
		'<div class="dialog_bd_inner" style="margin:30px 60px 20px;">',
		    '<div class="frm_control_group">',
		        '<div class="frm_controls">',
		            '<span class="frm_input_box with_counter counter_in append count" style="padding:0px;width:auto;height:auto;line-height:auto;">',
		                '<textarea style="width:600px;height:150px;" class="frm_input content" >{content}</textarea>',
		            '</span>',
		            '<span class="frm_msg fail js_openidfail"></span>',
		        '</div>',
		    '</div>',
		'</div>'
	];

	var listtpl = [
			'<tr>',
				'<td class="table_cell info">',
					'<div class="info_inner">',
						'<p class="wx_account">{content}</p>',
					'</div>',
				'</td>',
				'<td class="table_cell opr" style="width:120px;"><div class="opr_inner">',
					'<a href="javascript:;" class="js_kf_edit" data-id="{id}">编辑<input type="hidden" class="edit_value" value="{content}" /></a>',
					'<a href="javascript:;" class="js_kf_del" data-id="{id}">删除</a>',
				'</div></td>',
			'</tr>'
		];

	var emptytpl = [
			'<tr>',
				'<td colspan="3" style="height:75px;line-height:75px;">还没有添加快捷回复</td>',
			'</tr>'
		];

	$(document.body).on('click','.js_kf_edit', function(){
		var lastdata = {
			id:$(this).attr('data-id'),
			content:$(this).find('.edit_value').val()
		};
		Addfn('文字回复', 'update', lastdata);
	});

	$(document.body).on('click','.js_kf_del', function(){
		if(confirm("确认删除？")){
			$.ajax({
				url : '/request/fastreply/remove',
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
			url : '/request/fastreply/list',
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
		Addfn('文字回复', 'add');
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

			var content = dialog.find('.content').val().trim();	

			if(content == ""){
				alert('回复内容为空');
				return;
			}

			var data = {
					content:content
				};

			if(lastdata){
				data.id = lastdata.id;
			}

			$.ajax({
				url : '/request/fastreply/' + action,
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
