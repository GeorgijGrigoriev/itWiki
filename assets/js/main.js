$(window).on('load', function(){
    setPersonName();
    populateFrontpageCards();
    getCategories();
    addNewArticle();
    getCategoriesTable();
    addNewCategory();
    $('#itwiki_logout').on('click', function(){
        deleteCookie("username");
        deleteCookie("token");
        window.location.href = "/";
    });
});



function getJWTToken(name){
    var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    if (match) return match[2];
}

function setPersonName(){
    var $field = $('#itwiki_logged_in_person'),
    $username = getJWTToken("username");
    $field.text($username);
}

function deleteCookie(cname) {
    var d = new Date(); //Create an date object
    d.setTime(d.getTime() - (1000*60*60*24)); //Set the time to the past. 1000 milliseonds = 1 second
    var expires = "expires=" + d.toGMTString(); //Compose the expirartion date
    window.document.cookie = cname+"="+"; "+expires;//Set the cookie with name and the expiration date
 
}

function populateFrontpageCards(){
    $.when(
        $.ajax({
            url: "/api/articles/all/",
            dataType: "json",
            contentType: "application/json",
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
            }
        })
    ).then(function(succes, fail){
            if (succes.length !== 0){
                $cardRow = $('#itwiki_frontpage_cards');
                for (i=0; i < succes.total; i++){
                    var $postPreview = succes.data[i].post;
                    if (succes.data[i].post.length > 50) $postPreview = $postPreview.substring(0,50);
                    var $clearfix = $('<div />', {
                        id: succes.data[i].article_id
                    }).appendTo($cardRow),
                    $cardBody = $('<div />', {
                        class: "uk-card uk-card-default uk-card-hover uk-card-body",
                        id: succes.data[i].article_id
                    }).appendTo($clearfix),
                    $cardHeader = $('<div />', {
                        class: "uk-card-header"
                    }).appendTo($cardBody),
                    $cardTitle = $('<h3 />', {
                        class: "uk-card-title",
                        text: succes.data[i].title
                    }).appendTo($cardHeader),
                    $cardHtmlBody = $('<div />', {
                        class: "uk-card-body"
                    }).appendTo($cardBody),
                    $cardPost = $('<p />', {
                        text: $postPreview,
                    }).appendTo($cardHtmlBody),
                    $cardFooter = $('<div />', {
                        class: "uk-card-footer",
                        html: "<button class='uk-button uk-button-text article-read-more' data-id="+ succes.data[i].article_id +" onclick='readMoreHandler("+ succes.data[i].article_id +")'>Read more</button>"
                    }).appendTo($cardBody);

            }
        }
        if (succes.total == -1) {
            $cardRow = $('#itwiki_frontpage_cards');
            var $clearfix = $('<div />', {
            }).appendTo($cardRow),
            $cardBody = $('<div />', {
                class: "uk-card uk-card-default uk-card-hover uk-card-body",
            }).appendTo($clearfix),
            $cardHeader = $('<div />', {
                class: "uk-card-header"
            }).appendTo($cardBody),
            $cardTitle = $('<h3 />', {
                class: "uk-card-title",
                text: "Статей нет"
            }).appendTo($cardHeader),
            $cardHtmlBody = $('<div />', {
                class: "uk-card-body"
            }).appendTo($cardBody),
            $cardPost = $('<p />', {
                text: "В данной WIKI еще нет статей, вы можете их добавить.",
            }).appendTo($cardHtmlBody);
        }
    });
}


function getCategories(){
    $.when(
        $.ajax({
            url: "/api/categories/get/",
            dataType: "json",
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
            }
        })
    ).then(function(succes, fail){
        var $select = $("#itwiki_autogenerated_select");
        for (var i=0; i<succes.total; i++){
            $option = $('<option />', {
                value: succes.data[i].id,
                text: succes.data[i].category_name
            }).appendTo($select);
        }
    })
}

function getCategoriesTable(){
    $.ajax({
        url: "/api/categories/get/",
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
        },
        dataType: "json",
        complete: function(e){
            var $table = $('#itwiki_category_table'),
            data = (e.responseJSON);
            $table.empty();
            for (i=0; i < e.responseJSON.total; i++){
                $table.append('<tr><td>' + data.data[i].id + '</td>\
                <td>' + data.data[i].category_name + '</td></tr>'
                );
            }
        }
    });
}

function addNewCategory(){
    var $input = $("#itwiki_add_new_category_input"),
    $btn = $("#itwiki_add_new_category");
    $btn.on('click', function(e){
        e.preventDefault();
        $data = $input.serializeArray(),
        $data = indexArray($data),
        $data = JSON.stringify($data);
        console.log($data);
        $.ajax({
            url: "/api/categories/add/",
            method: "POST",
            data: $data,
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
            },
            complete: function(succes){
                getCategoriesTable();
                $("#itwiki_add_new_category_input").empty();
            }
        })
    })
}

function addNewArticle(){
    $articleForm = $('#itwiki_add_new_article');
    $articleForm.on('submit', function(e){
        e.preventDefault();
        var data = $articleForm.serializeArray(),
        data = indexArray(data),
        data = JSON.stringify(data);
        $.when(
            $.ajax({
                url: "/api/articles/add/",
                method: "POST",
                contentType: "application/json",
                dataType: "json",
                data: data,
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
                }
            })
        ).then(function(succes, fail){
            window.location.href = "/app";
        })
    });
}

function readMoreHandler(article_id){
    localStorage.setItem("currentArticle", article_id);
    window.location.href = "/app/article/read/";
}

function indexArray($array) {
    var ua = $array,
    ia = {};
    $.map(ua, function(n, i){
        ia[n['name']] = n['value']
    });
    return ia;
}
