var Template = {
	'dialog' : [
		'<div style="width: 726px; margin-left: -363px; margin-top: -202px;" class="dialog_wrp align_edge ui-draggable">',
			'<div class="dialog">',
				'<div class="dialog_hd">',
					'<h3>{Title}</h3>',
					'<a class="icon16_opr closed pop_closed" onclick="return false" href="javascript:;">关闭</a>',
				'</div>',
				'<div class="dialog_bd">',
					'<div class="whitelist_dialog">',
						'{Content}',
					'</div>',
				'</div>',
				'<div class="dialog_ft">',
					'<span class="btn btn_primary btn_input js_btn_p" style="display: inline-block;">',
						'<button data-index="0" class="js_btn submitbt" type="button">',
							'{Button}',
						'</button>',
					'</span>',
				'</div>',
			'</div>',
		'</div>',
		'<div class="mask ui-draggable">',
			'<iframe frameborder="0" src="about:blank" style="filter:progid:DXImageTransform.Microsoft.Alpha(opacity:0);position:absolute;top:0px;left:0px;width:100%;height:100%;">',
			'</iframe>',
		'</div>'
	]
};
