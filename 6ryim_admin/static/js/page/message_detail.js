(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = 'http://im1.6renyou.com:8989';
	var SECOND = 1000;
	var LAST_ID = null;


	if(PageData.data.status == 1){
		var pd = PageData.data;
		$('#realname').text(pd.realname + '(' + pd.mobile + ')');
		$('#order_title').text(pd.order_title);
		$('#start_date').text(pd.start_date);
		$('#days').text(pd.detail.days + '天');
		$('#operator').text('Op: ' + pd.operator.name);
	}

	var ORDER_ID = getQueryParam('orderid');
	$('#order_label').text(ORDER_ID);
	$('#order_label').css('color', 'red');
	
	var msg_type = $('#msgtypeselect').dropdown();
	var msg_source = $('#msgsourceselect').dropdown();
	var searchBt = $('.js_reply_OK');

	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;text-align:center;">',
						'<tr>',
							'<td style="width:70px;">',
								'<img src="{from_icon}" />',
							'</td>',
							'<td style="text-align:left;">',
								'<div class="{from}_label">{from_name}</div>',
								'<div>{message}</div>',
							'</td>',
							'<td style="width:100px;">{msgtype_name}</td>',
							'<td style="width:150px;">{createtime}</td>',
							'<td style="width:100px;">{source_name}</td>',
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
				orderid:ORDER_ID
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
			url:"/user/get",
			data:{
				openids:userids.join(","),
				source:source.join(",")
			},
			type:'POST',
			dataType:'json',
			success:function(data, status, response){
				if(data.status != 0){
					//alert(data.info);
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
		listContainer.hideLoading();
		listContainer.showLoading();
		listContainer.html('');
		var size = 20;
		$.ajax({
			url : SERVICE_DOMAIN + '/message/show',
			data:{
				page:toIndex,
				key:$('.jsSearchInput').val(),
				msgtype:msg_type.getValue(),
				msgsource:msg_source.getValue(),
				size:size,
				orderid:ORDER_ID
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
							if(d.fromtype == "1"){
								source.push("weixin");
							}else{
								source.push(d.source == "1" ? "weixin" : "app");
							}
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
							case '4':
								d.source_name = 'DCloud';
								break;
							default:
						}

						if(d.fromtype == "1"){
							d.from_icon = "http://admin.6renyou.com/statics/socketchat/img/six-service.jpg";
						}else{
							d.from_icon = "http://admin.6renyou.com/statics/socketchat/img/default-user.jpg";
						}

						d.from_name = d.from;
						d.to_name = 'unknwon';

						listContainer.append(listtpl.replaceTpl(d));	
					});
					loadUser(userids, source);
					listContainer.hideLoading();

					if(data.total == 0){
						listContainer.html('<div style="text-align:center;"> No Data </div>')
						pageNavigator.hide();
						return;
					}
					pageNavigator.show();

					if(!isInitPageNavs){
						var pageSize = Math.ceil(data.total/size) + data.total%size;
						PageNavCtls = pageNavigator.PN({
							recordCount:data.total,
							pageSize:size,
							showPageNum:8,
							jump:function(to){
								listContainer.showLoading();
								loadMsgList(to);
							}
						});
						isInitPageNavs = true;
						if(LAST_ID == null){
							setTimeout(function(){monitor(data.data[0].id);}, SECOND * 10);
							LAST_ID = data.data[0].id;
						}
					}else{
						PageNavCtls.defaults.pageIndex = toIndex;
						PageNavCtls.defaults.recordCount = data.total;
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
