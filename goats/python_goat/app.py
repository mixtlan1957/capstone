# SOURCE: https://code.tutsplus.com/tutorials/creating-a-web-app-from-scratch-using-python-flask-and-mysql--cms-22972

import getpass
import traceback
from flask import escape, Flask, json, render_template, request
from flaskext.mysql import MySQL
app = Flask(__name__)

mysql = MySQL()

app.config['MYSQL_DATABASE_USER'] = 'root'
app.config['MYSQL_DATABASE_PASSWORD'] = getpass.getpass('Please enter your password: ')
app.config['MYSQL_DATABASE_DB'] = 'FlaskGoat'
app.config['MYSQL_DATABASE_HOST'] = 'localhost'
mysql.init_app(app)

connection = mysql.connect()

cursor = connection.cursor()

@app.route("/")
def main():
    return render_template('index.html')

@app.route('/createAccount')
def createAccount():
    return render_template('account.html')

@app.route('/login')
def login():
    return render_template('login.html')

@app.route('/messages')
def forum():
    return render_template('messages.html')

@app.route('/search')
def src():
    return render_template('search.html')

@app.route('/loginAttempt', methods=['POST'])
def loginAttempt():
    username = request.form['username']
    password = request.form['password']

    if username and password:
        try:
            cursor.execute("SELECT `password` from `FlaskGoat`.`users` WHERE `username` = %s;", (username,))
            data = cursor.fetchall()
            connection.commit()

            if data[0][0] == password:
                return json.dumps({'html':'<div id="formValid">Login info correct</div>'})
            else:
                return json.dumps({'html':'<div id="formValid">Login info NOT correct</div>'})

        except Exception as e:
            return json.dumps({
                'html':'<div id="formValid">Login Error</div>',
                'error':str(e)
            })
    
    else:
        return json.dumps({'html':'<div id="formValid">All fields are not valid</div>'})

@app.route('/accountMade', methods=['POST'])
def accountMade():
    name = request.form['name']
    username = request.form['username']
    password = request.form['password']

    if name and username and password:
        try:
            query = "INSERT INTO `FlaskGoat`.`users` (`name`, `username`, `password`) VALUES ('%s', '%s', '%s');" % (name, username, password)
            cursor.execute(query)
            connection.commit()

            return json.dumps({
                'html':'<div id="formValid">All fields are valid</div>',
                'error':None
            })
        except Exception as e:
            return json.dumps({
                'html':'<div id="formValid">Account Creation Error</div>',
                'error':str(e)
            })
    else:
        return json.dumps({'html':'<div id="formInvalid">Please fill in all fields</div>'})

@app.route('/sendMessage', methods=['POST'])
def messageReceived():
    name = request.form['name']
    message = request.form['message']
    
    if name and message:
        try:
            cursor.execute("INSERT INTO `FlaskGoat`.`messages` (`name`, `message`) VALUES (%s, %s);", (name, message,))
            connection.commit()

            try:
                cursor.execute("SELECT name, message, msg_time FROM FlaskGoat.messages ORDER BY msg_time DESC;")
                msgData = cursor.fetchall()

                msgDiv = "<div>"
                for msg in range(0, len(msgData)):
                    msgDiv += "<div><p><b>" + str(escape(msgData[msg][0])) + " " + str(escape(msgData[msg][2])) + "</b></p>"
                    msgDiv += "<p>" + str(msgData[msg][1]) + "</p></div>"
                
                msgDiv += "</div>"

                return json.dumps({
                    'html':msgDiv,
                    'error':None
                })
            
            except Exception as e:
                return json.dumps({
                    'html':'<div id="formValid">Messages could not be retrieved</div>',
                    'error':str(e)
                })

            return json.dumps({
                'html':'<div id="formValid">Message Submitted</div>',
                'error':None
            })
        except Exception as e:
            return json.dumps({
                'html':'<div id="formValid">Message Insert Error</div>',
                'error':str(e)
            })
    
    else:
        return json.dumps({'html':'<div id="formInvalid">Please fill in all fields</div>'})

@app.route('/searchUsername', methods=['POST'])
def searchQueryReceived():
    username = request.form['username']

    if username:
        try:
            cursor.execute("SELECT `username` from `FlaskGoat`.`users` WHERE `username` REGEXP %s;", (username,))
            searchResults = cursor.fetchall()
            connection.commit()

            srchDiv = "<div>"
            for uname in range(0, len(searchResults)):
                srchDiv += "<p>" + str(escape(searchResults[uname][0])) + "</p>"
            
            srchDiv += "</div>"

            return json.dumps({
                'html': srchDiv,
                'query': str(escape(username)),
                'error': None
            })

        except Exception as e:
            return json.dumps({
                'html':'<div id="formValid">Search error</div>',
                'error':str(e)
            })
    
    else:
        return json.dumps({'html':'<div id="formValid">All fields need to be filled in</div>'})


if __name__ == "__main__":
    app.run(port=8080)