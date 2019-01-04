var Login = {

	init: function() {
		Login.setForm();
		Login.setRememberMe();
		Login.retrieveRememberedUsername();
	},

	setForm: function() {
		var form = document.getElementById('frmLogin')

    if(form) {
			/*form.addEventListener('submit', function(event){
				Login.validateForm(form);
				event.preventDefault();
			});*/
		}		
	},

	validateForm: function(form) {
		if(!Login.validateElement(form.username)) return false;
		if(!Login.validateElement(form.password)) return false;
	},

	validateElement: function(element) {
		var elementGroup = $(element)
		.parents('.control-group');
		
		if(element.value === "") {
			$(elementGroup).addClass('error');
			element.focus();
			return false;
		}

		$(elementGroup).removeClass('error');
		return true;
	},

	setRememberMe: function() {
		var inputCheckbox = document.getElementById('rememberMe');
		inputCheckbox.addEventListener('click', function(event){
			if(inputCheckbox.checked) {
				var username = $('#username').val();
				if(username) {
					Cookie.set('username', username, 3);
				}
			}
			else {
				Cookie.unset('username');
			}
		});
	},

	retrieveRememberedUsername: function() {
		var username = Cookie.get('username');
		if(username) {
			$('#username').val(username);
			//$('#rememberMe').attr('checked', true);
			document.getElementById('rememberMe').checked = true;
		}
	}

};

//initialization
Login.init();


function login(){
	if(!Login.validateElement(frmLogin.username)) return false;
	if(!Login.validateElement(frmLogin.password)) return false;

	$("#loginload").css('display',"block");
	$('#loginbtn').attr('disabled','disabled');
	$('#loginbtn').css('color','#337ab7')
	
	$.post(
		"loginaction",
		{
			user:frmLogin.username.value,
			pass:frmLogin.password.value
		},
		function(data,status){
			$("#loginload").css('display',"none");
			$('#loginbtn').removeAttr('disabled');
			$('#loginbtn').css('color','#fff')
			$('#wrongnamelabel').css('display',"none")
			$('#wrongpasslabel').css('display',"none")
			if(data.Status == 'failure'){
				if(data.Msg == "Wrong password"){
					$('#password').css("border-color","orangered");
					$('#wrongpasslabel').css('display',"block")
				}else{
					$('#username').css('border-color',"orangered");
					$('#wrongnamelabel').css('display',"block")
				}
			}else{
				location.href = data.Msg
			}
		}
	)
}
document.onkeydown = function(event){
	var e  = event || window.event || arguments.callee.caller.arguments[0];
	if(e && e.keyCode == 13){
		login()
	}
}

