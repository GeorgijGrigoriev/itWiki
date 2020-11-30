$(window).on('load', function(){
    var x = getJWTToken("username");
    console.log(x);
    if (x !== undefined) {
        window.location.href = "/app";
    }
    handleLogin();
});

function handleLogin(){
    $('#itwiki_login_form').on('submit', function(e){
    e.preventDefault();
    var $data = $(this).serializeArray(),
    $data = indexArray($data),
    $data = JSON.stringify($data);
    $.ajax({
        url: "/auth/login",
        method: "POST",
        contentType: "application/json",
        data: $data,
        complete: function(succes){
            var $result = succes.responseJSON;
            if ($result.Status == false){
                alert($result.Message);
            }
            if ($result.Status == true){
                document.cookie = "username=" + $result.account.username;
                document.cookie = "token=" + $result.account.token;
                window.location.href = "/app";
            }
        }
    });
    });
}

function indexArray($array) {
    var ua = $array,
    ia = {};
    $.map(ua, function(n, i){
        ia[n['name']] = n['value']
    });
    return ia;
}

function getJWTToken(name){
    var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    if (match) return match[2];
}