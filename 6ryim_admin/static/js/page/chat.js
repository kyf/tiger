(function($, window){
	$('.webuploader-element-invisible').css({'opacity':0});

	var text_tpl = {
		'rightplain' : [
                     '<div class="plain">',
                          '<pre  class="js_message_plain ng-binding">{content}</pre>',
                          '<img alt="" style="display:none" src="https://res.wx.qq.com/zh_CN/htmledition/v2/images/icon/ico_loading28a2f7.gif" class="ico_loading" > <i title="{error}"   class="ico_fail web_wechat_message_fail" style="display:none;"></i>',
					'</div>'
		],
		'rightpicture' : [
					'<div class="picture">',
	                       '<a href="{content}" target="_blank"><img class="msg-img" src="{content}" style="width: 100px;"></a>',
	                       '<i class="arrow"></i>',
	                       '<p class="loading ng-hide">',
	                            '<img src="https://res.wx.qq.com/zh_CN/htmledition/v2/images/icon/ico_loading28a2f7.gif"><i title="{error}"   class="ico_fail web_wechat_message_fail" style="display:none;"></i>',
	                        '</p>',
	              	'</div>'
		],
		'leftplain':[
					'<div class="plain">',
						'<pre class="js_message_plain ng-binding">{content}</pre>',
						//'<img alt="" src="/images/ico_loading28a2f7.gif" class="ico_loading ng-hide" >',
						//'<i title="重新发送" class="ico_fail web_wechat_message_fail ng-hide"></i>',
					'</div>',
			],
		'leftpicture':[
					'<div class="picture">',
			        	'<a href="{content}" target="_blank"><img class="msg-img" src="{content}" style="width: 100px;"></a>',
						'<i class="arrow"></i>',
						'<p class="loading ng-hide">',
							'<img alt="" src="https://res.wx.qq.com/zh_CN/htmledition/v2/images/icon/ico_loading28a2f7.gif">',
						'</p>',
					'</div>'
			],
		'leftaudio':[
					'<div class="voice" jqcontent="{content}" style="width: 47px;">',
			        	'<i class="voice_icon web_wechat_voice_gray"></i>',
						'<span class="duration ng-binding"><i class="web_wechat_noread ng-hide"></i></span>',
					'</div>'
			]
	};

	var CurrentVoice = null, CurrentContent = null;

	$(document.body).on('click', '.voice', function(){
		var amr = $(this).attr('jqcontent');
		if(CurrentVoice){
			CurrentVoice.stop();
			if(amr == CurrentContent)return;
		}
		CurrentContent = amr;
		var icon = $(this).find('.voice_icon');
		icon.removeClass('web_wechat_voice_gray');
		icon.addClass('web_wechat_voice_gray_playing');
		playRemoteVoice(amr, function(){
			icon.removeClass('web_wechat_voice_gray_playing');
			icon.addClass('web_wechat_voice_gray');
			CurrentVoice = null;
			CurrentContent = null;
		});
	});

	var playRemoteVoice = function (_file, cb) {
		if (!_file || "" == _file || (-1 == _file.indexOf(".amr"))) {
			return false;
		}
		var playVoice = function (url, cb) {
			var oReq = new XMLHttpRequest();
			oReq.onload = function (e) {
				var arraybuffer = oReq.response;
				var _array = new Uint8Array(arraybuffer);
				var samples = AMR.decode(_array);
				if (!samples) {
					alert('Failed to decode!');
					return;
				}
				else {
					var ctx = new AudioContext();
					var src = ctx.createBufferSource();
					src.onended = function(){
						cb();
					};
					var buffer = ctx.createBuffer(1, samples.length, 8000);
					if (buffer.copyToChannel) {
						buffer.copyToChannel(samples, 0, 0);
					} else {
						var channelBuffer = buffer.getChannelData(0);
						channelBuffer.set(samples);
					}

					src.buffer = buffer;
					src.connect(ctx.destination);
					src.start();
					CurrentVoice = src;
				}
			}
			oReq.open("GET", url, true);
			oReq.responseType = "arraybuffer";
			oReq.send();
		};
		playVoice(_file, cb);
	};


	var chattpl = [
		'<div  class="ng-scope chat_list_item">',
			'<div message-directive="" class="clearfix">',
				'<div  style="overflow: hidden;" on="message.MsgType" ng-switch="">',
					'<div  class="message ng-scope me" ng-switch-default="">',
						//'<p class="message_system ng-scope" ><span class="content ng-binding">10:57</span></p>',
						'<img src="http://admin.6renyou.com/statics/socketchat/img/six-service.jpg" class="avatar">',
						'<div class="content">',
							'<h4 class="nickname ng-binding ng-scope" style="width:auto;margin-right:10px;font-size:14px;font-weight:normal;">{opid}</h4>',
							'<div class="bubble js_message_bubble ng-scope bubble_primary right">',
								'<div  class="bubble_cont ng-scope">',
									'{main_content}',
								'</div>',
                    		'</div>',
                  		'</div>',
                	'</div>',
					//'<div style="float:right;font-size:14px;margin-top:-16px;margin-bottom:16px;">中国人的什么</div>',
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
					'<img src="{headurl}" class="avatar">',
					'<div class="content">',
						'<div class="bubble js_message_bubble ng-scope bubble_default left">',
							'<div class="bubble_cont ng-scope">',
								'{main_content}',
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
			'<div class="ng-scope useritem useritem_{openid}">',
                '<div class="chat_item slide-left ng-scope {active}" source="{source}" openid="{openid}" openid_name="{username}">',
                  '<div class="avatar"> ',
					  '<img src="{headurl}" class="img"> ',
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

	var addUserItem = function(username, lastmsg, active, number, openid, msgtype, headurl, source){
		switch(parseInt(msgtype)){
			case MSG_TYPE_IMAGE:
				lastmsg = "[图片]";
				break;
			case MSG_TYPE_AUDIO:
				lastmsg = "[语音]";
				break;
		}
		var data = {
			openid:openid,
			username : username,
			lastmsg : lastmsg,
			active : active ? "active" : '',
			number : number,
			source:source,
			headurl : headurl == '' ? 'http://admin.6renyou.com/statics/socketchat/img/default-user.jpg' : headurl,
			number_display:number > 0 ? 'display:' : 'display:none'
		};
		$('.UserContainer').before(usertpl.replaceTpl(data));	
	};

	var CurrentOpenid = '';
	var CurrentSource = '';
	var CurrentHeadurl = 'http://admin.6renyou.com/statics/socketchat/img/default-user.jpg';

	var addChatItem = function(content, media_id, msg_type){
		var user = "";
		if(!content){
			content = $('#editArea').val();
			$('#editArea').val('');
			user = USER;
			if(CurrentOpenid == ''){
				alert('当前没有用户');
				return;
			}
		}
		if(content.trim().length == 0)return;
		if(!msg_type){
			msg_type = MSG_TYPE_TEXT;
		}
		var main_content = "";
		switch(msg_type){
			case MSG_TYPE_TEXT:
				main_content = text_tpl["rightplain"].join('').replaceTpl({content:ToggleFace(content)});
				break;
			case MSG_TYPE_IMAGE:
				main_content = text_tpl["rightpicture"].join('').replaceTpl({content:content});
				break;
		}
		var item = $(chattpl.replaceTpl({main_content:main_content, opid:user}));
		$('.ChatContainer').before(item);
		$('.MainChatContainer').scrollTop($('.MainChatContainer').get(0).scrollHeight);
		$(item).find('.ico_loading').show();

		$.ajax({
			url:'/request/send',
			data:{
				openid:CurrentOpenid,
				msg_type:msg_type,
				message:content,
				source:CurrentSource,
				media_id:media_id
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
					var wait = data.data.wait;
					data = data.data;
					if(wait > 0){
						window.parent.$('#newMsgTip').show();
						window.parent.$('#newWaitNum').text(wait);
					}
					if(data.data.length == 0)return;
					$('.useritem').remove();
					data.data = data.data.sort(sort_user);
					if(CurrentOpenid == ''){
						CurrentOpenid = data.data[0].openid;
						CurrentSource = data.data[0].source;
						if(data.data[0].headurl != ''){
							CurrentHeadurl = data.data[0].headurl;
						}
						$('.title_name').text(data.data[0].openid_name);
						loadHistory();
					}
					$.each(data.data, function(i, d){
						addUserItem(d.openid_name, d.msg, CurrentOpenid == d.openid ? true : false, d.number, d.openid, d.msgType, d.headurl, d.source);
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
						var in_tpl;
						switch(d.msgType){
							case MSG_TYPE_TEXT:
								d.content = ToggleFace(d.content);
								in_tpl = text_tpl['leftplain'];
								break;
							case MSG_TYPE_IMAGE:
								in_tpl = text_tpl['leftpicture'];
								break;
							case MSG_TYPE_AUDIO:
								in_tpl = text_tpl['leftaudio'];
								break;
						}
						var main_content = in_tpl.join('').replaceTpl(d);
						$('.ChatContainer').before(chatlefttpl.replaceTpl({main_content:main_content, headurl:CurrentHeadurl}));
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
						var tpl = chatlefttpl, in_tpl;
						switch(d.msgType){
							case MSG_TYPE_TEXT:
								d.content = ToggleFace(d.content);
								in_tpl = text_tpl['leftplain'];
								break;
							case MSG_TYPE_IMAGE:
								in_tpl = text_tpl['leftpicture'];
								break;
							case MSG_TYPE_AUDIO:
								in_tpl = text_tpl['leftaudio'];
								break;
						}

						if(d.opid !== ""){
							tpl = chattpl;

							switch(d.msgType){
								case MSG_TYPE_TEXT:
									d.content = ToggleFace(d.content);
									in_tpl = text_tpl['rightplain'];
									break;
								case MSG_TYPE_IMAGE:
									in_tpl = text_tpl['rightpicture'];
									break;
							}
						}
						var main_content = in_tpl.join('').replaceTpl(d);
						$('.ChatContainer').before(tpl.replaceTpl({main_content:main_content, opid:d.opid, headurl:CurrentHeadurl}));
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
		var source = $(this).attr('source');
		var openid_name = $(this).attr('openid_name');
		var headurl = $(this).find('img').attr('src');
		CurrentOpenid = openid;
		CurrentSource = source;
		CurrentHeadurl = headurl;
		$('.title_name').text(openid_name);
		loadHistory();
	});

	window.callback = function(data){
		if(data.status){
			data = data.data;
			addChatItem(data.filepath, data.media_id, MSG_TYPE_IMAGE);
		}else{
			alert(data.msg);
		}
	}
	$('.webuploader-element-invisible').change(function(){
		$('#myform').submit();
	});

	$('.web_wechat_face').click(function(){
		$('#mmpop_emoji_panel').toggle();
	});

	$('.qq_face').find('.face').click(function(){
		var val = $('#editArea').val();	
		$('#editArea').val(val + '[' + this.title + ']');
	});


	$(document.body).on('contextmenu', '.chat_item', function(e){
		var openid = $(this).attr('openid');
		var source = $(this).attr('source');
		$('#contextMenu').css({top:e.clientY, left:e.clientX}).show();
		$('#contextMenu').find('.bookitem').attr('data-openid', openid);
		$('#contextMenu').find('.bookitem').attr('data-source', source);
		$('#contextMenu').find('.closeitem').attr('data-openid', openid);
		return false;
	});

	$(document.body).on('click', function(e){
		$('#contextMenu').hide();
		if(e.target != $('.web_wechat_reply').get(0)){
			ReplyPanel.hide();
		}

		//if(e.target != $('.web_wechat_face').get(0) && !$(e.target).hasClass('face')){
		if(e.target != $('.web_wechat_face').get(0)){
			$('#mmpop_emoji_panel').hide();
		}
	});

	$(document.body).on('click', '.bookitem', function(){
		var openid = $(this).attr('data-openid');
		var source = $(this).attr('data-source');
		if(source == MSG_SOURCE_PC){
			alert('匿名用户目前不支持代客下单');
			return;
		}
		window.open('/call/center/my/book?openid=' + openid + "&opid=" + OPID);
	});

	var clearScreen = function(openid){
		if(openid == CurrentOpenid){
			$('.chat_list_item').remove();
			CurrentOpenid = '';
			CurrentSource = '';
			$('.title_name').text('');
		}
		$('.useritem_' + openid).remove();
	};


	$(document.body).on('click', '.closeitem', function(){
		var openid = $(this).attr('data-openid');
		
		$.ajax({
			url:"/unbind",
			data:{
				openid:openid
			},
			dataType:'json',
			type:'POST',
			success:function(data){
				if(data.status){
					clearScreen(openid);	
				}else{
					alert(data.msg);
				}
			}
		});
	});

	var ReplyPanel = null;

	(function(){
		if(ReplyPanel == null){
			var reply_tpl = [
				'<div style="width:310px;position:absolute;right:20px;top:70px;background:white;border-radius:3px;border:1px solid #cfcfcf;">',
					'<div style="padding:5px 0px 0px 0px;background:#f2f2f2;">',
						'<div style="background:white;font-size:14px;width:90px;margin-left:15px;line-height:30px;height:30px;text-align:center;:">快捷回复</div>',
					'</div>',
					'<div class="ReplyContainer" style="padding:0px 8px;height:313px;overflow:auto;"><div style="text-align:center;line-height:80px;">loading...</div></div>',
				'</div>'
			];	

			reply_tpl = reply_tpl.join('');
			ReplyPanel = $(reply_tpl);
			$(document.body).append(ReplyPanel);
			var item_tpl = [
				'<div style="{border};" class="reply_item">',
					'{content}',
				'</div>'
			];
			item_tpl = item_tpl.join('');
			$.ajax({
				url:'/request/fastreply/list',
				type:'POST',
				dataType:'json',
				success:function(data){
					if(data.status){
						if(!data.data)data.data = [];
						ReplyPanel.find('.ReplyContainer').html('');
						$.each(data.data, function(index, d){
							if(index > 0){
								d.border = 'border-top:1px solid #f2f2f2;';
							}
							ReplyPanel.find('.ReplyContainer').append(item_tpl.replaceTpl(d));
						});
					}else{
						ReplyPanel.find('.ReplyContainer').find('div').html(data.msg);
					}
				}
			});
		}
	
		ReplyPanel.hide();
	})();

	$('.web_wechat_reply').click(function(){
		ReplyPanel.toggle();
	});

	$(document.body).on('click', '.reply_item', function(){
		var val = $('#editArea').val();
		$('#editArea').val(val + $(this).text());
		ReplyPanel.hide();
	});

})(jQuery, window)
