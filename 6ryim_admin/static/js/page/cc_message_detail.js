(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = '';
	var SECOND = 1000;
	var LAST_ID = null;


	
	var ORDER_ID = getQueryParam('openid');
	$('#order_label').css('color', 'red');
	
	var msg_type = $('#msgtypeselect').dropdown();
	var searchBt = $('.js_reply_OK');

	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;text-align:center;">',
						'<tr>',
							'<td style="width:70px;">',
								'<img src="{from_icon}" style="width:40px;height:40px;" class="{from}_avatar" />',
							'</td>',
							'<td style="text-align:left;">',
								'<div class="{from}_label">{from_name}</div>',
								'<div>{message}</div>',
							'</td>',
							'<td style="width:100px;">{msgtype_name}</td>',
							'<td style="width:100px;">{source_name}</td>',
							'<td style="width:150px;">{createtime}</td>',
						'</tr>',
					'</table>',
				'</li>'
	];
	listtpl = listtpl.join('');

	var isInitPageNavs = false,
		PageNavCtls = null;

	var monitor = function(lastid){
		$.ajax({
			url : SERVICE_DOMAIN + '/request/message/new/number',
			data:{
				lastid:lastid,
				openid:ORDER_ID
			},
			dataType:'json',
			type:'POST',
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
		if(!userids || !source)return;
		if(userids.length == 0 || source.length == 0)return;
		$.ajax({
			url:"/wx/user/get",
			data:{
				openids:userids.join(",")
			},
			type:'POST',
			dataType:'json',
			success:function(data, status, response){
				if(data.status != 0){
					return;
				}
				data = data.data;
				
				if(data.length > 0){
					$.each(data, function(i, d){
						$('.' + d.openid + "_label").text(d.nickname);
						if(d.openid == ORDER_ID){
							$('#order_label').text(d.nickname);	
						}
						$('.' + d.openid + "_avatar").attr('src', d.headimgurl);
					})
				}
			}
		});
	};


	var loadMsgList = function(toIndex){
		listContainer.hideLoading();
		listContainer.showLoading();
		listContainer.html('');
		var size = 20;
		$.ajax({
			url : SERVICE_DOMAIN + '/request/message/show',
			data:{
				page:toIndex,
				keyword:$('.jsSearchInput').val(),
				msg_type:msg_type.getValue(),
				size:size,
				openid:ORDER_ID
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.status){
					if(data.data.total == 0)data.data.data = [];
					var tmpkv = new Object(), userids = new Array(), source = new Array();
					$.each(data.data.data, function(i, d){
						d.from = d.openid;
						if(!tmpkv[d.from] && d.from != 'system'){
							userids.push(d.from);
							source.push("weixin");
							tmpkv[d.from] = true;
						}
						d.message = d.content;
						switch(d.msgType){
							case MSG_TYPE_TEXT:
								d.msgtype_name = '文本';
								break;
							case MSG_TYPE_IMAGE:
								d.msgtype_name = '图片';
								d.message = '<a href="' + SERVICE_DOMAIN + d.message + '" target="_blank"><img style="width:100px;height:100px;" src="' + SERVICE_DOMAIN + d.message + '"/></a>';
								break;
							case MSG_TYPE_AUDIO:
								d.msgtype_name = '语音';
								break;
							default:
						}

						switch(d.source){
							case MSG_SOURCE_WX:
								d.source_name = "微信";
								break;
							case MSG_SOURCE_IOS:
								d.source_name = "IOS";
								break;
							case MSG_SOURCE_Android:
								d.source_name = "Android";
								break;
							case MSG_SOURCE_PC:
								d.source_name = "PC";
								break;
							default:
								d.source_name = "未知";
						}


						
						d.from_name = d.from;
						if(d.fromtype == 2){
							d.from_icon = "http://admin.6renyou.com/statics/socketchat/img/six-service.jpg";
							d.from = '';
							d.from_name = d.opid;
						}else{
							d.from_icon = "http://admin.6renyou.com/statics/socketchat/img/default-user.jpg";
						}

						d.to_name = 'unknwon';
						d.createtime = ts2time(d.ts);

						listContainer.append(listtpl.replaceTpl(d));	
					});
					loadUser(userids, source);
					listContainer.hideLoading();

					if(data.data.total == 0){
						listContainer.html('<div style="text-align:center;"> No Data </div>')
						pageNavigator.hide();
						return;
					}
					pageNavigator.show();

					if(!isInitPageNavs){
						var pageSize = Math.ceil(data.data.total/size) + data.data.total%size;
						PageNavCtls = pageNavigator.PN({
							recordCount:data.data.total,
							pageSize:size,
							showPageNum:8,
							jump:function(to){
								listContainer.showLoading();
								loadMsgList(to);
							}
						});
						isInitPageNavs = true;
						if(LAST_ID == null){
							setTimeout(function(){monitor(data.data.data[0].ts);}, SECOND * 10);
							LAST_ID = data.data.data[0].ts;
						}
					}else{
						PageNavCtls.defaults.pageIndex = toIndex;
						PageNavCtls.defaults.recordCount = data.data.total;
						PageNavCtls.defaults.pageSize = size;
						PageNavCtls.update();
					}
				}
			}
		});
	};

	loadMsgList(1);

	searchBt.click(function(){
		loadMsgList(1);
	});
	$('.search_gray').click(function(){
		loadMsgList(1);
	});
})(jQuery, window)
