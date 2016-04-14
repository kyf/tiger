(function($, window){
	$('.menu_item').click(function(){
		$('.menu_item').removeClass('selected');
		$(this).addClass('selected');
	});


	var prefix = 'www.6renyou.com';
	var path = prefix + window.location.pathname.toLowerCase();
	$('.menu_item').each(function(){
		var it = prefix + $(this).find('a').attr('href');

		if(path.indexOf(it) > -1){
			$(this).addClass('selected');
		}else{
			$(this).removeClass('selected');
		}
	});

})(jQuery, window)

