(function($, window){
	var chattpl = [
		'<div  class="ng-scope">',
            '<div message-directive="" class="clearfix">',
              '<div  style="overflow: hidden;" on="message.MsgType" ng-switch="">',
                '<div  class="message ng-scope me" ng-switch-default="">',
                  '<p class="message_system ng-scope" ><span class="content ng-binding">10:57</span></p>',
                  '<img  title="风筝" src="http://admin.6renyou.com/statics/socketchat/img/six-service.jpg" class="avatar">',
                  '<div class="content">',
                    '<div class="bubble js_message_bubble ng-scope bubble_primary right">',
                      '<div  class="bubble_cont ng-scope">',
                        '<div class="plain">',
                          '<pre  class="js_message_plain ng-binding">{content}</pre>',
                          '<img alt="" src="https://res.wx.qq.com/zh_CN/htmledition/v2/images/icon/ico_loading28a2f7.gif" class="ico_loading ng-hide" > <i title="重新发送"   class="ico_fail web_wechat_message_fail ng-hide"></i> </div>',
                      '</div>',
                    '</div>',
                  '</div>',
                '</div>',
              '</div>',
            '</div>',
       '</div>'
	];

	chattpl = chattpl.join('')

	var usertpl = [
			'<div class="ng-scope">',
                '<div class="chat_item slide-left ng-scope">', //active
                  '<div class="avatar"> ',
					  '<img src="http://admin.6renyou.com/statics/socketchat/img/default-user.jpg" class="img"> ',
				  '</div>',
                  '<div class="info">',
                    '<h3 class="nickname"> <span  class="nickname_text ng-binding">{username}</span> </h3>',
                    '<p class="msg ng-scope"> <span class="ng-binding">{lastmsg}</span> </p>',
                  '</div>',
                '</div>',
              '</div>'
	];
	usertpl = usertpl.join('');

	var addUserItem = function(username, lastmsg){
		var data = {
			username : username,
			lastmsg : lastmsg
		};
		$('.UserContainer').before(usertpl.replaceTpl(data));	
	};


	var addChatItem = function(content){
		if(!content){
			content = $('#editArea').val();
		}
		if(content.trim().length == 0)return;
		$('#editArea').val('');
		var data = {
			content : content
		};
		$('.ChatContainer').before(chattpl.replaceTpl(data));
		$('.MainChatContainer').scrollTop($('.MainChatContainer').get(0).scrollHeight);
	};
	$('.btn_send').click(addChatItem);
	$('#editArea').on('keyup', function(ev){
		if(ev.keyCode == 13 && ev.ctrlKey){
			addChatItem();
		}
	});


	addUserItem("6人游", "说设么");
	addUserItem("6人游", "说设么");
	addUserItem("6人游", "说设么");


})(jQuery, window)
