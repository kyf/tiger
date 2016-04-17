(function($, window){
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = '';
	
	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;text-align:center;">',
						'<tr>',
							'<td style="width:70px;">',
								'<img src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" />',
							'</td>',
							'<td style="text-align:left;">',
								'<div class="{from}_label">{from_name}</div>',
								'<div>{msg}</div>',
							'</td>',
							'<td style="width:100px;">{msgtype_name}</td>',
							'<td style="width:150px;">{ts}</td>',
							'<td style="width:200px;color:red" >',
								'<span class="btn btn_primary btn_input">',
									'<button class="js_fetch" jqopenid="{openid}">接入</button>',
								'</span>',
							'</td>',
						'</tr>',
					'</table>',
				'</li>'
	];
	listtpl = listtpl.join('');

	var loadUser = function(userids, source){
		if(!userids || !source)return;
		if(userids.length == 0 || source.length == 0)return;
		$.ajax({
			url:"/user/get",
			data:{
				openids:userids.join(","),
				source:source.join(",")
			},
			type:'POST',
			dataType:'json',
			success:function(data, status, response){
				if(data.status != 0){
					alert(data.info);
					return;
				}
				data = data.data;
				if(data.length > 0){
					$.each(data, function(i, d){
						$('.' + d.userid + "_label").text(d.realname + '(' + d.mobile + ')');
					})
				}
			}
		});
	};

	var loadMsgList = function(toIndex){
		listContainer.html('');
		var size = 20;
		$.ajax({
			url : SERVICE_DOMAIN + '/request/wait',
			data:{
				size:size,
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.data == null)data.data = [];
				if(data.data){

					var tmpkv = new Object(), userids = new Array(), source = new Array();

					$.each(data.data, function(i, d){
						d.from = d.openid;
						if(!tmpkv[d.from] && d.from != 'system'){
							userids.push(d.from);
							source.push("weixin");
							tmpkv[d.from] = true;
						}
						switch(d.msgType){
							case '2':
								d.msgtype_name = '文本';
								break;
							case '3':
								d.msgtype_name = '图片';
								d.message = '<a href="' + SERVICE_DOMAIN + d.message + '" target="_blank"><img style="width:100px;height:100px;" src="' + SERVICE_DOMAIN + d.message + '"/></a>';
								break;
							case '4':
								d.msgtype_name = '语音';
								break;
							default:
						}

						d.from_name = d.from;
						d.to_name = 'unknwon';

						listContainer.append(listtpl.replaceTpl(d));	
					});
					loadUser(userids, source);

					if(data.total == 0){
						listContainer.html('<div style="text-align:center;"> No Data </div>')
						return;
					}

				}
			}
		});
	};

	loadMsgList(1);

	$(document.body).on('click', '.js_fetch', function(){
		var openid = $(this).attr('jqopenid');	

		$.ajax({
			url : '/request/bind',
			data:{
				openid:openid
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					window.location.href = '/call/center/my';
				}else{
					alert(data.msg);
				}
			}
		});
	});

})(jQuery, window)
