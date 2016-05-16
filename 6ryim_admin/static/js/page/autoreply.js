(function($, window){
	var ListContainer = $('#js_list');
	var listtpl = [
			'<div>',
				'<table>',
					'<tr>',
						'<td class="table_cell info">',
							'<div class="info_inner">',
								'<div class="setting_time dropdown_wrp dropdown_menu"></div>',
								'<div class="setting_time dropdown_wrp dropdown_menu"></div>',
							'</div>',
						'</td>',
						'<td class="table_cell opr" style="width:120px;text-align:center;"><div class="opr_inner">',
							'<a href="javascript:;" class="js_kf_edit" data-id="{id}">保存<input type="hidden" class="edit_value" value="{content}" /></a>',
							'<a href="javascript:;" class="js_kf_del" data-id="{id}">删除</a>',
						'</div></td>',
					'</tr>',
				'</table>',
			'</div>'
		];
	listtpl = listtpl.join('');

	var getTimeData = function(){
		if(window.TIME_DATA){
			return window.TIME_DATA;
		}
		var hours = [];
		var minutes = [];
		$.each(new Array(24), function(index){
			hours.push({text:index, value:index});
		});
		$.each(new Array(60), function(index){
			minutes.push({text:index, value:index});
		});

		var result = {
			hours : hours,
			minutes : minutes
		};

		window.TIME_DATA = result;
		return result;
	};

	var SettingTimes = [];

	var addTimeItem = function(current){
		var data = getTimeData();
		var item = $(listtpl);
		ListContainer.append(item);
		var its = item.find('.setting_time');
		its.each(function(index, it){
			var st = $(it);
			st = st.dropdown({data:index == 0 ? data.hours : data.minutes, current:current});
			SettingTimes.push(st);
		});
	};


	var emptytpl = [
			'<div>',
				'<div style="height:75px;line-height:75px;">还没有设置时间段自动回复</div>',
			'</div>'
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

	//loadAccount();

	$('.js_kf_add').click(function(){
		addTimeItem({});
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
