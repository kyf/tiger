(function($, window){
	$.fn.dropdown = function(config){

		this.init = function(data){
			var tpl = [];
			$.each(data, function(index, it){
				if(index == 0){
					tpl.push('<a class="btn dropdown_switch jsDropdownBt" href="javascript:;" data-index="' + index + '" data-value="'+it.value+'" ><label class="jsBtLabel">'+it.text+'</label><i class="arrow"></i></a>');
					tpl.push('<div class="dropdown_data_container jsDropdownList" style="display: none;">');
					tpl.push('<ul class="dropdown_data_list">');
				}else{
					tpl.push('<li  class="dropdown_data_item >');
					tpl.push('<a data-name="'+it.text+'" data-index="'+index+'" data-value="'+it.value+'" class="jsDropdownItem" href="javascript:;" onclick="return false;">'+it.text+'</a>');
					tpl.push('</li>');
				}
			});

			tpl.push('</ul>');
			tpl.push('</div>');
			this.html(tpl.join(''));
		};


		if(config && config.data){
			this.init(config.data);
		}

		var triggerBt = this.find('.jsDropdownBt'),
			listPanel = this.find('.jsDropdownList');

		var label = triggerBt.find('.jsBtLabel'),
			items = listPanel.find('.dropdown_data_item');

		triggerBt.click(function(){
			listPanel.toggle();
		});

		items.click(function(){
			var name = $(this).find('.jsDropdownItem').attr('data-name');
			var val = $(this).find('.jsDropdownItem').attr('data-value');

			label.text(name);
			triggerBt.attr('data-value', val);
		});

		this.getValue = function(){
			return triggerBt.attr('data-value');
		};

		

		$(document.body).click(function(evt){
			if(evt.target != label.get(0)){
				listPanel.hide();
			}
		});

		return this;
	}
})(jQuery, window)
