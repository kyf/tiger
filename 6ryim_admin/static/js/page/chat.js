(function($, window){
	var MESSAGE_TYPE_TEXT = 1, MESSAGE_TYPE_IMAGE = 2, MESSAGE_TYPE_AUDIO = 3;


	var chattpl = [
		'<div  class="ng-scope chat_list_item">',
            '<div message-directive="" class="clearfix">',
              '<div  style="overflow: hidden;" on="message.MsgType" ng-switch="">',
                '<div  class="message ng-scope me" ng-switch-default="">',
                  //'<p class="message_system ng-scope" ><span class="content ng-binding">10:57</span></p>',
                  '<img src="http://admin.6renyou.com/statics/socketchat/img/six-service.jpg" class="avatar">',
                  '<div class="content">',
                    '<div class="bubble js_message_bubble ng-scope bubble_primary right">',
                      '<div  class="bubble_cont ng-scope">',
                        '<div class="plain">',
                          '<pre  class="js_message_plain ng-binding">{content}</pre>',
                          '<img alt="" style="display:none" src="https://res.wx.qq.com/zh_CN/htmledition/v2/images/icon/ico_loading28a2f7.gif" class="ico_loading" > <i title="{error}"   class="ico_fail web_wechat_message_fail" style="display:none;"></i> </div>',
                      '</div>',
                    '</div>',
                  '</div>',
                '</div>',
              '</div>',
            '</div>',
       '</div>'
	];

	chattpl = chattpl.join('')


	var chatlefttpl = [
	'<div class="ng-scope chat_list_item">',
		'<div class="clearfix">',
			'<div style="overflow: hidden;" on="message.MsgType" ng-switch="">',
				'<div class="message ng-scope you" ng-switch-default="">',
					'<img src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" class="avatar">',
					'<div class="content">',
						'<div class="bubble js_message_bubble ng-scope bubble_default left">',
							'<div class="bubble_cont ng-scope">',
								'<div class="plain">',
									'<pre class="js_message_plain ng-binding">{content}</pre>',
									//'<img alt="" src="/images/ico_loading28a2f7.gif" class="ico_loading ng-hide" >',
									//'<i title="重新发送" class="ico_fail web_wechat_message_fail ng-hide"></i>',
								'</div>',
							'</div>',
						'</div>',
					'</div>',
				'</div>',
			'</div>',
		'</div>',
	'</div>'
	];
	chatlefttpl = chatlefttpl.join('');




	var usertpl = [
			'<div class="ng-scope useritem">',
                '<div class="chat_item slide-left ng-scope {active}" openid="{openid}" openid_name="{username}">',
                  '<div class="avatar"> ',
					  '<img src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" class="img"> ',
					  '<i style="{number_display}" class="unread_number {openid}_unread_number icon web_wechat_reddot_middle ng-binding ng-scope">{number}</i>',
				  '</div>',
                  '<div class="info">',
                    '<h3 class="nickname"> <span  class="nickname_text ng-binding">{username}</span> </h3>',
                    '<p class="msg ng-scope"> <span class="ng-binding">{lastmsg}</span> </p>',
                  '</div>',
                '</div>',
              '</div>'
	];
	usertpl = usertpl.join('');

	var addUserItem = function(username, lastmsg, active, number, openid){
		var data = {
			openid:openid,
			username : username,
			lastmsg : lastmsg,
			active : active ? "active" : '',
			number : number,
			number_display:number > 0 ? 'display:' : 'display:none'
		};
		$('.UserContainer').before(usertpl.replaceTpl(data));	
	};

	var CurrentOpenid = '';

	var addChatItem = function(content){
		if(!content){
			content = $('#editArea').val();
		}
		if(content.trim().length == 0)return;
		$('#editArea').val('');
		var data = {
			content : content
		};
		var item = $(chattpl.replaceTpl(data));
		$('.ChatContainer').before(item);
		$('.MainChatContainer').scrollTop($('.MainChatContainer').get(0).scrollHeight);
		$(item).find('.ico_loading').show();

		$.ajax({
			url:'/request/send',
			data:{
				openid:CurrentOpenid,
				msg_type:MESSAGE_TYPE_TEXT,
				message:content
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					$(item).find('.ico_loading').hide();
				}else{
					$(item).find('.web_wechat_message_fail').show();
					$(item).find('.web_wechat_message_fail').attr('title', data.msg);
				}	
			}
		});

	};

	$('.btn_send').click(function(){
		addChatItem();
	});

	$('#editArea').on('keyup', function(ev){
		if(ev.keyCode == 13 && ev.ctrlKey){
			addChatItem();
		}
	});


	var sort_user = function(a, b){
		return a.ts < b.ts;
	};


	var loadMyUser = function(){
		$.ajax({
			url:'/request/cc',
			data:{},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					if(data.data.length == 0)return;
					$('.useritem').remove();
					data.data = data.data.sort(sort_user);
					if(CurrentOpenid == ''){
						CurrentOpenid = data.data[0].openid;
						$('.title_name').text(data.data[0].openid_name);
						loadHistory();
					}
					$.each(data.data, function(i, d){
						addUserItem(d.openid_name, d.msg, CurrentOpenid == d.openid ? true : false, d.number, d.openid);
						if(CurrentOpenid == d.openid && d.number > 0){
							getUnread();
						}
					});	
				}else{
					alert(data.msg);
				}
			}
		});
	};

	(function(){
		var cb = arguments.callee;
		loadMyUser();
		setTimeout(cb, 10000);
	})();

	var getUnread = function(isload){
		var openid = CurrentOpenid;
		$('.' + CurrentOpenid + '_unread_number').hide();
		$.ajax({
			url : '/request/fetch',
			data:{
				openid:CurrentOpenid
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					if(isload)return;
					if(openid != CurrentOpenid)return;
					data = data.data.unread;
					if(data.length == 0)return;
					$.each(data, function(i, d){
						$('.ChatContainer').before(chatlefttpl.replaceTpl(d));
						$('.MainChatContainer').scrollTop($('.MainChatContainer').get(0).scrollHeight);
					});

				}			
			}
		});

	};

	var loadHistory = function(){
		$.ajax({
			url : '/request/message/list',
			data:{
				openid:CurrentOpenid
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					data = data.data;
					if(data.length == 0)return;
					$('.chat_list_item').remove();
					data = data.reverse();
					$.each(data, function(i, d){
						var tpl = chatlefttpl;
						if(d.opid !== ""){
							tpl = chattpl;
						}
						$('.ChatContainer').before(tpl.replaceTpl(d));
						$('.MainChatContainer').scrollTop($('.MainChatContainer').get(0).scrollHeight);
					});

					getUnread("unload");
				}				
			}
		});

	};

	$(document.body).on('click', '.chat_item', function(){
		if($(this).hasClass('active'))return;
		$('.chat_item').removeClass('active');
		$(this).addClass('active');

		$(this).find('.unread_number').hide();

		var openid = $(this).attr('openid');
		var openid_name = $(this).attr('openid_name');
		CurrentOpenid = openid;
		$('.title_name').text(openid_name);
		loadHistory();
	});


})(jQuery, window)
