import os, getpass, random
from app import app
from flask import current_app, Flask, json, g, render_template, flash, redirect, request, jsonify
from flask import url_for, session
import hashlib
import sqlite3 as sql
import click
from flask.cli import with_appcontext
app.secret_key = "afiendirtme3o,lkd"
#from flaskext.mysql import MySQL

#mysql = MySQL()

DATABASE = "app/goatdb.db"


#if os.path.exists('goatdb.db'):
#    os.remove('goatdb.db')


'''
#connect to sql db
connection = sql.connect()
connection.row_factory = sql.Row

cursor = connection.cursor()
'''

#source: http://flask.pocoo.org/docs/1.0/tutorial/database/
def get_db():
    if 'db' not in g:
        g.db = sql.connect(
            DATABASE,
            detect_types=sql.PARSE_DECLTYPES
        )
        g.db.row_factory = sql.Row
        #g.db.row_factory = dict_factory

    return g.db

def close_db(e=None):
    db = g.pop('db', None)

    if db is not None:
        db.close()

def init_db():
    db = get_db()

    with current_app.open_resource('schema.sql') as f:
        db.executescript(f.read().decode('utf8'))
    db.commit()
    print("initizlied the database...")


@click.command('init-db')
@with_appcontext
def init_db_command():
    """Clear the existing data and create new tables."""
    init_db()
    click.echo('Initialized the database.')

#source: https://docs.python.org/2/library/sqlite3.html
def dict_factory(cursor, row):
    d = {}
    for idx, col in enumerate(cursor.description):
        d[col[0]] = row[idx]
    return d


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


@app.route('/')
def index():
    if 'name' in session:
        name = session['username']
        flash("Logged in as: \"" + name + "\"!!!")     


    init_db()
    return render_template('home.html')



@app.route('/login', methods = ['POST'])
def login():
    if request.method == 'POST':
        try:
            #get form input
            userName = request.form['login']
            pw = request.form['pw']
           
            #check if user and pw is in db
            esc_user = sqlFilter(userName)
            esc_user = str(esc_user)
            print("escaped username: " + esc_user)
            lQuery = "SELECT salt, user_id FROM main.users WHERE eid = \""+esc_user+"\";"
            print(lQuery)
            
            cursor = get_db().cursor()
            cursor.execute(lQuery)
            result = cursor.fetchone()
            #result = cursor.fetchall()
            #cursor.close()

              
            if (result == None) :
                flash("Error logging in. Incorrect credencials supplied...")
                cursor.close()
                raise Exception("User not found!")

            #verify the passowrd hash
            salt = result[0]
            hash = hashlib.md5((salt+pw).encode('utf-8'))
            #hash.update((salt+pw).encode('utf-8'))
            

            lQuery = "SELECT user_id, name, eid FROM main.users WHERE eid = \""+esc_user+"\""
            lQuery = lQuery + " AND password = \""+str(hash.digest())+"\";"
            print(lQuery)

            cursor.execute(lQuery)
            result = cursor.fetchone()

            if (result == None):
                flash("Error logging in. Incorrect credencials supplied...")
                cursor.close()
                raise Exception("Login error.")

            #otherwise store the session
            session['username'] = userName
            session['eid'] = result['eid']
            session['name'] = result['name']
            cursor.close()

        except Exception as e:
            
            print(e)

        return redirect('/')

@app.route('/register', methods = ['POST'])
def register():
    if request.method == 'POST':
        try:
            u_name = request.form['name']
            u_accountName = request.form['login']
            u_pw1 = request.form['pw1']
            u_pw2 = request.form['pw2']

            #connect to db
            db_connection = get_db()
            #cursor = get_db().cursor()
            cursor = db_connection.cursor()

            #check if user already exists
            squery = "SELECT user_id FROM main.users WHERE eid= \""+u_accountName+"\";"
            cursor.execute(squery)
            result = cursor.fetchone()

            if result != None:
                cursor.close()
                flash("User already exists!")
                raise Exception("User already exists!")

            #check if user field is empty
            if u_accountName == None or u_name == None:
                cursor.close()
                flash("Could not fulfil request, user and/or user name fields are empty.")


            #check if passwords do not match
            if u_pw1 != u_pw2:
                flash("Passwords do not match.")
                cursor.close()
                raise Exception("Passwords do not match.")

            #create user
            salt = str(random.randint(10000, 99999))
            hash = hashlib.md5()
            hash.update((salt+u_pw1).encode('utf-8'))
            h = str(hash.digest())

            squery = ("INSERT INTO main.users (name, eid, password, salt) VALUES "
                "(\""+u_name+"\", \""+u_accountName+"\", \""+h+"\", \""+salt+"\");")

            print("we made it here")
            print(squery)
            cursor.execute(squery)
            db_connection.commit()

            flash("User: "+u_accountName+" has been created. Please log in.")
            cursor.close()

        except Exception as e:
            print(e)

        return redirect('/')

            

@app.route('/logout', methods = ['POST'])
def logout():
    if request.method == 'POST':
        session.clear()
        flash("Successfully logged out.")
        return redirect('/')