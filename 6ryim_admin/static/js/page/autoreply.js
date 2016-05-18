(function($, window){
	var ListContainer = $('#js_list');
	var listtpl = [
			'<div class="TimeItem">',
				'<table style="width:100%;">',
					'<tr>',
						'<td class="table_cell info" style="width:80px;padding:0px;text-align:center;">',
							'<div class="info_inner">',
								'<div class="setting_time dropdown_wrp dropdown_menu" style="width:auto;"></div>',
							'</div>',
						'</td>',
						'<td class="info" style="width:18px;">',
							'时',
						'</td>',
						'<td class="table_cell info" style="width:80px;padding:0px;text-align:center;">',
							'<div class="info_inner">',
								'<div class="setting_time dropdown_wrp dropdown_menu" style="width:auto;"></div>',
							'</div>',
						'</td>',
						'<td class="info" style="width:18px;">',
							'分',
						'</td>',
						'<td class="info" style="width:20px;">',
							'->',
						'</td>',
						'<td class="table_cell info" style="width:80px;padding:0px;text-align:center;">',
							'<div class="info_inner">',
								'<div class="setting_time dropdown_wrp dropdown_menu" style="width:auto;"></div>',
							'</div>',
						'</td>',
						'<td class="info" style="width:18px;">',
							'时',
						'</td>',
						'<td class="table_cell info" style="width:80px;padding:0px;text-align:center;">',
							'<div class="info_inner">',
								'<div class="setting_time dropdown_wrp dropdown_menu" style="width:auto;"></div>',
							'</div>',
						'</td>',
						'<td class="info" style="width:18px;">',
							'分',
						'</td>',

						'<td class="table_cell info" style="padding:10px 0px;text-align:center;">',
							'<textarea style="width:90%;height:90px;" class="replyContent">{content}</textarea>',
						'</td>',
						'<td class="table_cell opr" style="width:120px;text-align:center;padding:10px 0px;"><div class="opr_inner">',
							'<a href="javascript:;" class="js_kf_edit" data-id="{id}">保存</a>',
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


	var addTimeItem = function(currents, currentdata){
		if(!currents)currents = [];
		if(!currentdata)currentdata = {};
		var data = getTimeData();
		var item = $(listtpl.replaceTpl(currentdata));
		ListContainer.append(item);
		var its = item.find('.setting_time');
		its.each(function(index, it){
			var st = $(it);
			st = st.dropdown({data:index % 2 == 0 ? data.hours : data.minutes, current:currents[index]});
		});
	};


	var emptytpl = [
			'<div class="NoReplyTip">',
				'<div style="height:75px;line-height:75px;">还没有设置时间段自动回复</div>',
			'</div>'
		];

	var tiptpl = [
			'<div style="position:fixed;top:200px;left:200px;width:200px;height:50px;background:#44b549;border-radius:5px;text-align:center;display:none;">',
				'<span style="font-weight:bold;line-height:50px;"></span>',
			'</div>'
		];
	var tip = $(tiptpl.join(''));
	$(document.body).append(tip);
	
	var showTip = function(msg, status){
		var width = $(document.body).width();
		tip.css('color', status ? 'white' : 'red');
		tip.css('left', (width - 200) / 2);
		tip.find('span').text(msg);
		tip.fadeIn(600, function(){
			tip.fadeOut(1000);
		});
	};


	var updateCacheAutoReply = function(){
		$.ajax({
			url:"/request/cacheautoreply"
		});
	};

	$(document.body).on('click','.js_kf_edit', function(){
		var par = $(this).parents('.TimeItem');
		var data = {
			id:$(this).attr('data-id'),
			content:par.find('.replyContent').val(),
			fromhour:par.find('.dropdown_switch').eq(0).attr('data-value'),
			fromminute:par.find('.dropdown_switch').eq(1).attr('data-value'),
			tohour:par.find('.dropdown_switch').eq(2).attr('data-value'),
			tominute:par.find('.dropdown_switch').eq(3).attr('data-value')
		};

		var url = "/request/autoreply/timeitem/add";
		if(data.id != ''){
			url = "/request/autoreply/timeitem/update";
		}

		$.ajax({
			url:url,
			data:data,
			type:'POST',
			dataType:'json',
			success:function(data){
				if(data.status){
					showTip('保存成功', true);					
					updateCacheAutoReply();
				}else{
					showTip(data.msg, false);
				}	
			}
		});
	});

	$(document.body).on('click','.js_kf_del', function(){
		var _this = $(this);
		if(confirm("确认删除？")){
			$.ajax({
				url : '/request/autoreply/timeitem/remove',
				dataType:'json',
				data:{
					id:$(this).attr('data-id')	
				},
				type:'POST',
				success:function(data){
					if(data.status){
						updateCacheAutoReply();
					}else{
						//alert(data.msg);
					}
				}
			});
			_this.parents('.TimeItem').remove();
			if($('.TimeItem').size() == 0){
				$('#js_list').html(emptytpl.join(''));
			}
		}	
	});


	var loadAutoReply = function(){
		$.ajax({
			url : '/request/autoreply/timeitem/list',
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					if(!data.data)data.data = [];
					data = data.data;
					if(data.length == 0){
						$('#js_list').html(emptytpl.join(''));
					}else{
						$.each(data, function(i, d){
							var ops = [
								{text:d.fromhour, value:d.fromhour}, 
								{text:d.fromminute, value:d.fromminute}, 
								{text:d.tohour, value:d.tohour},
								{text:d.tominute, value:d.tominute}, 
							];
							addTimeItem(ops, d);
						});
					}
				}else{
					alert(data.msg);
				}
			}
		});
	};

	loadAutoReply();

	$('.js_kf_add').click(function(){
		$('.NoReplyTip').remove();
		addTimeItem({});
	});

	var loadFirstAutoReply = function(){
			$.ajax({
			url:'/request/autoreply/first/load',
			dataType:'json',
			success:function(data){
				if(data.status){
					data = data.data[0];
					$('#FirtAutoReply').val(data.content);
				}else{
					showTip(data.msg, false);
				}
			}
		});

	};

	loadFirstAutoReply();

	$('.js_save_bt_first').click(function(){
		var content	= $('#FirtAutoReply').val().trim();
		if(content == ''){
			showTip('回复内容不能为空!', false);
			$('#FirtAutoReply').focus();
			return;
		}

		$.ajax({
			url:'/request/autoreply/first/save',
			data:{
				content:content
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					showTip('保存成功', true);
					updateCacheAutoReply();
				}else{
					showTip(data.msg, false);
				}
			}
		});
	});

})(jQuery,window)
