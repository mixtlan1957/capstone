# SOURCE: https://code.tutsplus.com/tutorials/creating-a-web-app-from-scratch-using-python-flask-and-mysql--cms-22972

import getpass
import traceback
from flask import Flask, json, render_template, request
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

@app.route('/loginAttempt', methods=['POST'])
def loginAttempt():
    username = request.form['username']
    password = request.form['password']

    if username and password:
        try:
            cursor.execute("SELECT `password` from `FlaskGoat`.`users` WHERE `username` = %s;" % (username))
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

if __name__ == "__main__":
    app.run(port=8080)