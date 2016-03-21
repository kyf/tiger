(function($, window){
	$.fn.dropdown = function(){
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
