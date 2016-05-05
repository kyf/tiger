(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = '';
	var SECOND = 1000;
	var LAST_ID = null;
	
	var msg_type = $('#msgtypeselect').dropdown();
	var searchBt = $('.js_reply_OK');

	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;text-align:center;">',
						'<tr>',
							'<td style="width:70px;">',
								'<a href="/call/center/message/detail?openid={openid}" target="_blank"><img class="{from}_avatar" src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" style="width:40px;height:40px;" /></a>',
							'</td>',
							'<td style="text-align:left;">',
								'<div><a href="/call/center/message/detail?openid={openid}" target="_blank" class="{from}_label">{from_name}</a></div>',
								'<div>{content}</div>',
							'</td>',
							'<td style="width:100px;">{msgtype_name}</td>',
							'<td style="width:150px;">{createtime}</td>',
							'<td style="width:150px;color:red" class="{openid}_reply" jqid="{id}"></td>',
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
				fromtype:1
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


	var loadOnline = function(cb){
			$.ajax({
				url : SERVICE_DOMAIN + '/request/online/list',
				dataType:'json',
				success:function(data, status, response){
					if(data.status){
						cb(data.data);
					}
				}
			});
	};

	var loadOpenidLastMessage = function(openid){
		loadOnline(function(onlinedata){
			$.ajax({
				url : SERVICE_DOMAIN + '/request/message/show',
				data:{
					page:1,
					size:1,
					fromtype:2,
					openid:openid
				},
				dataType:'json',
				type:'POST',
				success:function(data, status, response){
					data.data = data.data.data;
					if(!data.data)data.data = [0];
					var current = $('.' + openid + '_reply');
					var latestId = data.data[0].id;
					current.each(function(){
						var id = $(this).attr('jqid');
						if(latestId > id){
							$(this).html('[已回复]');
						}else{
							if(onlinedata[openid]){
								$(this).html('<span class="btn btn_disabled btn_input"><button >已接入</button></span>');
							}else{
								$(this).html('<span class="btn btn_primary btn_input"><button class="js_fetch" jqopenid="' + openid + '">接入</button></span>');
							}
						}	
					});
				}
			});
		});
	};


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
					alert('该用户已经接入');
				}
			}
		});
	});




	var loadUser = function(userids, source){
		if(!userids || !source)return;
		if(userids.length == 0 || source.length == 0)return;
		if(userids.length != source.length)return;
		$.ajax({
			url:"/wx/user/get",
			data:{
				openids:userids.join(",")
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
						$('.' + d.openid + "_label").text(d.nickname);
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
				fromtype:1
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.status){
					var tmpkv = new Object(), userids = new Array(), source = new Array();

					if(!data.data.data)data.data.data = [];
					$.each(data.data.data, function(i, d){
						if(!tmpkv[d.openid] && d.openid != ''){
							userids.push(d.openid);
							source.push("weixin");
							tmpkv[d.openid] = true;
						}
						d.message = d.content;
						switch(d.msgType){
							case MSG_TYPE_TEXT:
								d.msgtype_name = '文本';
								break;
							case MSG_TYPE_IMAGE:
								d.msgtype_name = '图片';
								d.content = '<a href="' + SERVICE_DOMAIN + d.message + '" target="_blank"><img style="width:100px;height:100px;" src="' + SERVICE_DOMAIN + d.message + '"/></a>';
								break;
							case MSG_TYPE_AUDIO:
								d.msgtype_name = '语音';
								break;
							default:
						}

						d.from = d.openid;
						d.from_name = d.openid;
						d.to_name = 'unknwon';
						d.createtime = ts2time(d.ts);

						listContainer.append(listtpl.replaceTpl(d));	
					});
					var hasOpenid = {};
					$.each(data.data.data, function(i, d){
						if(hasOpenid[d.openid])return
						loadOpenidLastMessage(d.openid);
						hasOpenid[d.openid] = true;
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
