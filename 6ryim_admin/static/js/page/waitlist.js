(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = 'http://im1.6renyou.com:8989';
	var SECOND = 1000;
	var LAST_ID = null;
	
	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;text-align:center;">',
						'<tr>',
							'<td style="width:70px;">',
								'<a href="/message/detail?orderid={orderid}" target="_blank"><img src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" /></a>',
							'</td>',
							'<td style="text-align:left;">',
								'<div><a href="/message/detail?orderid={orderid}" target="_blank" class="{from}_label">{from_name}</a></div>',
								'<div>{message}</div>',
							'</td>',
							'<td style="width:100px;">{msgtype_name}</td>',
							'<td style="width:150px;">{createtime}</td>',
							'<td style="width:200px;color:red" class="{orderid}_reply" jqid="{id}">',
								'<span class="btn btn_primary btn_input">',
									'<button class="js_fetch">接入</button>',
								'</span>',
							'</td>',
						'</tr>',
					'</table>',
				'</li>'
	];
	listtpl = listtpl.join('');

	var isInitPageNavs = false,
		PageNavCtls = null;

	var monitor = function(lastid){
		$.ajax({
			url : SERVICE_DOMAIN + '/message/new/number',
			data:{
				lastid:lastid,
				fromtype:"2"
			},
			dataType:'json',
			type:'POST',
			beforeSend:ajaxBeforeSend,
			success:function(data, status, response){
				if(data.data){
					var num = data.data;
					if(num > 0){
						$('#newMsgTip').show(true);
						$('#newMsgNum').text(num);
					}
				}	
				setTimeout(function(){monitor(lastid);}, SECOND * 10);
			}
		});
	};


	
	var loadUser = function(userids, source){
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
			url : SERVICE_DOMAIN + '/message/show',
			data:{
				page:toIndex,
				key:$('.jsSearchInput').val(),
				size:size,
				fromtype:"2"
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.data == null)data.data = [];
				if(data.data){

					var tmpkv = new Object(), userids = new Array(), source = new Array();

					$.each(data.data, function(i, d){
						if(!tmpkv[d.from] && d.from != 'system'){
							userids.push(d.from);
							source.push(d.source == "1" ? "weixin" : "app");
							tmpkv[d.from] = true;
						}
						switch(d.msgtype){
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

						switch(d.source){
							case '1':
								d.source_name = '微信';
								break;
							case '2':
								d.source_name = 'IOS';
								break;
							case '3':
								d.source_name = 'Android';
								break;
							default:
						}
						d.from_name = d.from;
						d.to_name = 'unknwon';

						listContainer.append(listtpl.replaceTpl(d));	
					});
					var hasOrderid = {};
					$.each(data.data, function(i, d){
						if(hasOrderid[d.orderid])return
						hasOrderid[d.orderid] = true;
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

	$('.search_gray').click(function(){
		loadMsgList(1);
	});
})(jQuery, window)
