$(window).on('load', function(){
    populateFrontpageCards();
});

function getCategories(){

}

function populateFrontpageCards(){
    $.when(
        $.ajax({
            url: "/api/articles/all",
            dataType: "json"
        })
    ).then(function(succes, fail){
            if (succes.length !== 0){
                console.log(succes);
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
                    $cardTitle = $('<h3 />', {
                        class: "uk-card-title",
                        text: succes.data[i].title
                    }).appendTo($cardBody),
                    $cardHtmlBody = $('<div />', {
                        class: "uk-card-body"
                    }).appendTo($cardBody),
                    $cardPost = $('<p />', {
                        text: $postPreview,
                    }).appendTo($cardHtmlBody),
                    $cardFooter = $('<div />', {
                        class: "uk-card-footer",
                        html: "<a href='/read/article/" + succes.data[i].article_id + "' class='uk-button uk-button-text'>Read more</a>"
                    }).appendTo($cardBody);

            }
        }
    });
}

function getArticles(){
    console.log("Func started");
    $.when(
        $.ajax({
            url: "/api/articles/all",
            dataType: "json"
        })
    ).then(function(succes, fail){
        if (succes.length !== 0){
            console.log("recived data");
            console.log(succes);
        } else {
            console.log("Empty return");
        }
    });
}