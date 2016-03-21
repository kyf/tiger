String.prototype.replaceTpl = function(data){
	var reg = /{([^}]+)}/g;
	var result = this.replace(reg, function(main, group){
		return data[group] ? data[group] : '';
	})
	return result;
}

String.prototype.trim = function(){
	return this.replace(/^\s+|\s+$/g, '');
};

