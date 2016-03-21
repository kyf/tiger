(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');
	var SERVICE_DOMAIN = 'http://im2.6renyou.com:8989';
	var SECOND = 1000;
	
	var msg_type = $('#msgtypeselect').dropdown();
	var msg_source = $('#msgsourceselect').dropdown();
	var searchBt = $('.js_reply_OK');

	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<table style="width:100%;">',
						'<tr>',
							'<td>{from}</td>',
							'<td>{to}</td>',
							'<td>{message}</td>',
							'<td>{msgtype}</td>',
							'<td>{createtime}</td>',
							'<td>{issystem}</td>',
							'<td>{source}</td>',
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
				lastid:lastid
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
				size:size
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.data){
					$.each(data.data, function(i, d){
						if(d.status != 1){
							d.jumpclass = 'jumpclass';
							d.displayPub = '';
						}else{
							d.displayPub = 'display:none;';
						}
						listContainer.append(listtpl.replaceTpl(d));	
					});
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
							pageSize:pageSize,
							showPageNum:4,
							jump:function(to){
								listContainer.showLoading();
								loadMsgList(to);
							}
						});
						isInitPageNavs = true;
						setTimeout(function(){monitor(data.data[0].id);}, SECOND * 10);
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
})(jQuery, window)
