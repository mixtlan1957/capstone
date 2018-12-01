// https://code.tutsplus.com/tutorials/creating-a-web-app-from-scratch-using-python-flask-and-mysql--cms-22972

$(function() {
    $('#login').click(function() {
        $.ajax({
            url: '/loginAttempt',
            data: $('form').serialize(),
            type: 'POST',
            dataType: 'json',
            success: function(response) {
                document.getElementById('loginStatus').innerHTML = response.html;
            },
            error: function(error) {
                document.getElementById('loginStatus').innerHTML = "<div class='errorMsg'>" + error + "</div>";
            }
        });
    });
});

$(function() {
    $('#createAccount').click(function() {
        $.ajax({
            url: '/accountMade',
            data: $('form').serialize(),
            type: 'POST',
            dataType: 'json',
            success: function(response) {
                document.getElementById('accountCreationStatus').innerHTML = response.html;
            },
            error: function(error) {
                document.getElementById('accountCreationStatus').innerHTML = "<div class='errorMsg'>" + error + "</div>";
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
                document.getElementById('messageContent').innerHTML = response.html;
            },
            error: function(error) {
                document.getElementById('messageContent').innerHTML = "<div class='errorMsg'>" + error + "</div>";
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
                document.getElementById('searchQuery').textContent = response.query;
                document.getElementById('searchResults').innerHTML = response.html;
            },
            error: function(error) {
                document.getElementById('searchResults').innerHTML = "<div class='errorMsg'>" + error + "</div>";
            }
        });
    });
});

$(function() {
    $('#passwordChange').click(function() {
        $.ajax({
            url: '/passwordChange',
            data: $('form').serialize(),
            type: 'POST',
            dataType: 'json',
            success: function(response) {
                console.log(response.html);
                document.getElementById('pwordChangeSuccess').innerHTML = response.html;
            },
            error: function(error) {
                console.log(error);
                document.getElementById('pwordChangeSuccess').innerHTML = "<div class='errorMsg'>" + error + "</div>";
            }
        });
    });
});