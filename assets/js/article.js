$(window).on('load', function(){
    if (/read/.test(window.location.href)){
        AJAXRouter("/api/article/get/", "POST", articleIDToJSON(localStorage.getItem("currentArticle")), "get");
    }
    if (/edit/.test(window.location.href)){
        AJAXRouter("/api/article/get/", "POST", articleIDToJSON(localStorage.getItem("currentArticle")), "edit");
    }
    deleteArticle();
    // $('table').addClass('uk-table uk-table-divider');
    mutationObserverHandler();
});

var monthNames_EN = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

function editArticle() {
    $.ajax({
        url: "/api/article/",
        contentType: "application/json",
        method: "POST",
        data: data,
        beforeSend: function (xhr) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
        },
        complete: function(succes){
            const $result = succes.responseJSON;
            console.log($result);
        }
    })
}

function deleteArticle(){
    var $btn = $('#itwiki_article_settings_delete_article');
    $btn.on('click', function(){
        var $id = $btn.data("id"),
        $data = articleIDToJSON($id);
        $.ajax({
            url: "/api/article/delete/",
            method: "DELETE",
            data: $data,
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
            },
            complete: function(succes){
                localStorage.setItem("currentArticle","");
                window.location.href = "/app";
            }
        })
    });
}

function Article(title, post, creation_date, content, category_name, article_id, category_id){
    this.article_id = article_id;
    this.title = title;
    this.post = post;
    this.creation_date = creation_date;
    this.content = content;
    this.category_name = category_name;
    this.category_id = category_id;
    this.get = function(){
        var $data = articleIDToJSON(this.article_id);
        console.log($data);
    }
}

function getJWTToken(name){
    var match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    if (match) return match[2];
}

function articleIDToJSON(article_id){
    $data = {"id" : String(article_id)},
    $data = JSON.stringify($data);
    return $data;
}

function indexArray($array) {
    var ua = $array,
    ia = {};
    $.map(ua, function(n, i){
        ia[n['name']] = n['value']
    });
    return ia;
}

function AJAXRouter(url, method, data, type){
    if (type == "get"){
        return $.when(
            $.ajax({
                url: url,
                contentType: "application/json",
                method: method,
                data: data,
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
                },
            })
        ).then(function(succes){
            let d = new Article(succes.data.title, succes.data.post, succes.data.creation_date, succes.data.content, succes.data.category_name, succes.data.article_id);
            generateArticle(d);
        })
    }
    if (type == "delete") {
        $.when(
            $.ajax({
                url: url,
                contentType: "application/json",
                method: method,
                data: data,
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
                },
            })
        ).then(function(succes){
            return succes;
        })
    }
    if (type == "edit") {
        $.when(
            $.ajax({
                url: url,
                contentType: "application/json",
                method: method,
                data: data,
                beforeSend: function (xhr) {
                    xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
                },
            })
        ).then(function(succes){
            let d = new Article();
            d.title = succes.data.title;
            d.post = succes.data.post;
            d.category_id = succes.data.category;
            d.article_id = succes.data.article_id;
            generateEditingArticle(d);
        }) 
    }
}

function generateArticle(resp) {
    $date = new Date(resp.creation_date);
    $dateString = "" + $date.getDate() + " " + monthNames_EN[$date.getMonth()] + " " + $date.getFullYear() + ", at " + $date.getHours() + ":" + $date.getMinutes() 
    $('.itwiki_article_title').text(resp.title);
    $('#itwiki_article_content_wrapper').append(resp.content);
    $('#itwiki_article_created_at').text($dateString);
    $('#itwiki_article_category').text(resp.category_name);
    $('#itwiki_article_settings_delete_article').attr({'data-id' : resp.article_id});
}

function generateEditingArticle(resp) {
    $('#itwiki_edit_article_title').val(resp.title);
    $('#itwiki_autogenerated_select').val(resp.category_id);
    $('#itwiki_edit_article_post').val(resp.post);
    $('#itwiki_article_article_id').val(resp.article_id);
    $('#itwiki_update_article').on('submit', function(e){
        e.preventDefault();
        let $form = $("#itwiki_update_article");
        let $data = $form.serializeArray();
        $data = indexArray($data),
        $data = JSON.stringify($data);
        $.ajax({
            url: "/api/article/update/",
            method: "POST",
            contentType: "application/json",
            data: $data,
            beforeSend: function (xhr) {
                xhr.setRequestHeader('Authorization', 'Bearer ' + getJWTToken("token"));
            },
            complete: function(succes){
                localStorage.setItem('currentArticle', succes.responseText);
                window.location.href = "/app/article/read/";
            }
        })
    });
}

function mutationObserverHandler() {
    var list = $('table');
    var MO = window.MutationObserver || window.WebKitMutationObserver || window.MozMutationObserver;
    var observer = new MO(function(mutations){
        mutations.forEach(function(mutation){
            if (mutation.type === 'childList'){
                console.log("mutation!");
            }
        });
    });
    observer.observe(list, {
        attributes: true,
        childList: true,
        characterData: true,
        subtree: true
    });
}