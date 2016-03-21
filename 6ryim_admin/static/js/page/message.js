(function($, window){
	var pageNavigator = $('.pageNavigator');
	var listContainer = $('#listContainer');

	var listtpl = [
				'<li data-id="577999267" id="msgListItem577999267" class="message_item ">',
					'<div class="message_opr">',
						'<a title="快捷回复" class="icon18_common reply_gray js_reply" data-tofakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" data-id="577999267" href="javascript:;">快捷回复</a>',
					'</div>',
					'<div class="message_info">',
						'<div class="message_status"><em class="tips">已回复</em></div>',
						'<div class="message_time">10:42</div>',
						'<div class="user_info">',
							'<a class="remark_name" data-id="577999267" data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" target="_blank" href="/cgi-bin/singlesendpage?tofakeid=oq2yCs5Fhta4ygg8hbpludhO9PgI&amp;t=message/send&amp;action=index&amp;quickReplyId=577999267&amp;token=189298351&amp;lang=zh_CN">风筝</a>',
							'<span data-id="577999267" data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" class="nickname"></span>',
							'<a style="display:none;" title="修改备注" data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" class="icon14_common edit_gray js_changeRemark" href="javascript:;"></a>',
							'<a data-id="577999267" data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" class="avatar" href="/cgi-bin/singlesendpage?tofakeid=oq2yCs5Fhta4ygg8hbpludhO9PgI&amp;t=message/send&amp;action=index&amp;quickReplyId=577999267&amp;token=189298351&amp;lang=zh_CN" target="_blank">',
								'<img data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" src="/misc/getheadimg?token=189298351&amp;fakeid=oq2yCs5Fhta4ygg8hbpludhO9PgI&amp;msgid=577999267">',
							'</a>',
						'</div>',
					'</div>',

					'<div class="message_content text">',
						'<div class="wxMsg " data-id="577999267" id="wxMsg577999267">好多话好多话</div>',
					'</div>',

					'<div class="js_quick_reply_box quick_reply_box" id="quickReplyBox577999267">',
						'<div class="emoion_editor_wrp js_editor"></div>',
						'<div class="verifyCode"></div>',
						'<p class="quick_reply_box_tool_bar">',
						'<span data-id="577999267" class="btn btn_primary btn_input">',
							'<button data-fakeid="oq2yCs5Fhta4ygg8hbpludhO9PgI" data-id="577999267" class="js_reply_OK">发送(Enter)</button>',
						'</span><a href="javascript:;" data-id="577999267" class="js_reply_pickup btn btn_default pickup">收起</a>',
						'</p>',
					'</div>',

				'</li>'
	];



	var loadMsgList = function(toIndex){
		listContainer.hideLoading();
		listContainer.showLoading();
		listContainer.html('');
		$.ajax({
			url : '/trip_service/getHotelList',
			data:{
				page:toIndex,
				s_name:$('#stay_keyword').val(),
				city:city
			},
			dataType:'json',
			type:'POST',
			success:function(data, status, response){
				if(data.status == 1){
					$.each(data.data.list, function(i, d){
						if(d.status != 1){
							d.jumpclass = 'jumpclass';
							d.displayPub = '';
						}else{
							d.displayPub = 'display:none;';
						}
						pn_stay_list.append(staytpl.replaceTpl(d));	
					});
					listContainer.hideLoading();

					if(data.data.count == 0){
						pn_stay_list.html('<div style="text-align:center;"> No Data </div>')
						pageNavigator.hide();
						return;
					}
					pageNavigator.show();

					if(!isInitPageNavs.stay){
						PageNavCtls.stay = pageNavigator.PN({
							recordCount:data.data.count,
							pageSize:data.data.size,
							showPageNum:4,
							jump:function(to){
								listContainer.showLoading();
								loadStayList(to);
							}
						});
						isInitPageNavs.stay = true;
					}else{
						PageNavCtls.stay.defaults.pageIndex = data.data.page;
						PageNavCtls.stay.defaults.recordCount = data.data.count;
						PageNavCtls.stay.defaults.pageSize = data.data.size;
						PageNavCtls.stay.update();
					}
				}
			}
		});
	};

})(jQuery, window)
