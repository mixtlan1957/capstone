// https://code.tutsplus.com/tutorials/creating-a-web-app-from-scratch-using-python-flask-and-mysql--cms-22972

$(function() {
    $('#createAccount').click(function() {
        $.ajax({
            url: '/accountMade',
            data: $('form').serialize(),
            type: 'POST',
            success: function(response) {
                console.log(response);
            },
            error: function(error) {
                console.log(error);
            }
        });
    });
});

$(function() {
    $('#login').click(function() {
        $.ajax({
            url: '/loginAttempt',
            data: $('form').serialize(),
            type: 'POST',
            success: function(response) {
                console.log(response);
            },
            error: function(error) {
                console.log(error);
            }
        });
    });
});

$(function() {
    $('#sendMessage').click(function() {
        $.ajax({
            url: '/sendMessage',
            data: $('form').serialize(),
            type: 'POST',
            dataType: 'json',
            success: function(response) {
                console.log(response);
                document.getElementById('messageContent').innerHTML = response.html;
            },
            error: function(error) {
                console.log("ERROR", error);
            }
        });
    });
});

$(function() {
    $('#searchBtn').click(function() {
        $.ajax({
            url: '/searchUsername',
            data: $('form').serialize(),
            type: 'POST',
            dataType: 'json',
            success: function(response) {
                console.log(response);
                document.getElementById('searchQuery').textContent = response.query;
                document.getElementById('searchResults').innerHTML = response.html;
            },
            error: function(error) {
                console.log("ERROR", error);
            }
        });
    });
});