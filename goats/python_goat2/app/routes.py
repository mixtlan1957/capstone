import os
from app import app
import sqlite3 as sql
from flask import json, g, render_template, flash, redirect, request, jsonify #for rendering html templates, request for dealing with POST
import hashlib

#conn = sqlite3.connect('database.db')
#print "Opened"

if os.path.exists('goatdb.db'):
    os.remove('goatdb.db')

conn = sql.connect('goatdb.db')



DATABASE = 'goatdb.db'



def get_db():
    db = getattr(g, '_database', None)
    if db is None:
        db = g._database = sql.connect(DATABASE)
    try:
        cur = db.cursor()
        cur.execute('CREATE TABLE IF NOT EXISTS users('
            'user_id INTEGER PRIMARY KEY AUTOINCREMENT, ' 
            'name TEXT, '
            'eid TEXT, ' 
            'password TEXT, ' 
            'salt TEXT, accounting TEXT);')
        msg = "loaded thing"

    except:
        msg = "error in loading db here"
        
    print(msg)

    db.row_factory = sql.Row
    return db


'''
def init_db():
    with app.app_context():
        db = get_db()
        with app.open_resource('schema.sql', mode='r') as f:
            db.cursor().executescript(f.read())
        db.commit()

init_db()
'''


def sqlFilter(strIn):
    filtered_string = strIn
    filtered_string = filtered_string.replace("--", "")
    filtered_string = filtered_string.replace("","")
    filtered_string = filtered_string.replace("/*","")
    filtered_string = filtered_string.replace("*/","")
    filtered_string = filtered_string.replace("//","")
    filtered_string = filtered_string.replace(" ","")
    filtered_string = filtered_string.replace("#","")
    filtered_string = filtered_string.replace("||","")
    filtered_string = filtered_string.replace("admin'","")
    filtered_string = filtered_string.replace("UNION","")
    filtered_string = filtered_string.replace("COLLATE","")
    filtered_string = filtered_string.replace("DROP","")
    return filtered_string





@app.teardown_appcontext
def close_connection(exception_):
    db = getattr(g, '_database', None)
    if db is not None:
        db.close()

@app.route('/')
def index():
    db = get_db()
    cur = db.cursor()
    cur.execute("SELECT * FROM users")
    result = cur.fetchall()
    db.commit()
    return render_template('home.html')


@app.route('/login', methods = ['POST'])
def login():
    if request.method == 'POST':
        try:
            #get form input
            userName = request.form['login']
            pw = request.form['pw']
            #db = get_db()
            #cur = db.cursor()
            #db.commit()

            #check if user and pw is in db
            escaped_username = sqlFilter(userName)
            print("escaped username: " + escaped_username)
            escaped_username = "belbe"
            lQuery = "SELECT `salt` FROM `users` WHERE `eid` = `belbe`"  
            print(lQuery)

            con = sql.connect(DATABASE)
            print("we here")
            con.row_factory = sql.Row

            cur = con.cursor()
            cur.execute(lQuery)

            result = cur.fetchall()
            con.commit()
            
            print("result:  " + result)

            if result[0] != "password":
                return json.dumps({'html':'<div id="formValid">Login info NOT correct</div>'})

            #verify the salt
            salt = result['salt']
            m = hashlib.md5(salt + pw)
            hashedPW = m.digest()

            pQuery = ("SELECT user_id, name, eid FROM users " 
            "WHERE eid=" + userName + " AND password=" + hashedPW +";")
            cur.execute(pQuery, None)
            userData = cur.fetchall()
            
            if userData:
                print("Successful login")
            else:
                print("Invalid Password")


        except Exception as e:
            msg = "error executing login query"
            print(msg)
            print(e)

        return redirect('/')





           


