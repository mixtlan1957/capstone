import os
from app import app
import sqlite3 as sql
from flask import json, g, render_template, flash, redirect, request, jsonify #for rendering html templates, request for dealing with POST
import hashlib
import click
from flask.cli import with_appcontext
from flask import current_app

#conn = sqlite3.connect('database.db')
#print "Opened"
if os.path.exists('goatdb.db'):
    os.remove('goatdb.db')

conn = sql.connect('goatdb.db')



DATABASE = 'goatdb'

'''
def init_db():
    db = get_db()

    with current_app.open_resource('schema.sql') as f:
        db.executescript(f.read().decode('utf8'))
@click.command('init-db')
@with_appcontext
def init_db_command():
    init_db()
    click.echo('Initialized the database.')


def get_db():
    if 'db' not in g:
        g.db = sql.connect(
            DATABASE,
            detect_types=sql.PARSE_DECLTYPES
        )
        g.db.row_factory = sql.Row

    return g.db


def close_db(e=None):
    db = g.pop('db', None)

    if db is not None:
        db.close()

'''


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
    #result = cur.fetchall()
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
            esc_username = sqlFilter(userName)
            esc_username = str(esc_username)
            print("escaped username: " + esc_username)
            lQuery = "SELECT salt FROM users WHERE eid = %s;" % "thing"
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





           


