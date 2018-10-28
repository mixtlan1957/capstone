from app import app
from flask import render_template, flash, redirect, request, jsonify #for rendering html templates, request for dealing with POST
from pymongo import MongoClient #for MongoDB / python
from app.forms import CrawlerForm

import requests #for sending POST requests to crawler and graph

"""
database, pull crawl data from DB and send to graphical UI
"""
client = MongoClient() #localhost, default port: 7018
db = client.pymongo_test #creates a new database
nodes = db.nodes


"""
routes
"""

@app.route('/', methods=['GET', 'POST'])
def index():
	form = CrawlerForm()
	if form.validate_on_submit():
		
		url = request.form.get('url')
		traversal = request.form.get('traversal')
		depth = request.form.get('depth')
		keyword = request.form.get('keyword') #need to use .get here because an optional key

		#data to send to crawler
		data = {"url": url, "traversal": traversal, "depth": depth, "keyword": keyword}

		res = requests.post('http://localhost:12345/crawl', json=data)

		"""
		res = requests.post('http://localhost:5000/json', json=data)

		print "resulting json"
		print res.json()

		# return res.json()

		return render_template('json.html', jsonresult = res.json()) #basically an echo
		"""
	return render_template('index.html', form=form)

@app.route("/about_us")
def aboutUs():
	return "Josh, Keane, and Mario"

@app.route("/faq")
def faq():
	return "Frequently Asked Questions"

#echos the json object sent by form, the sub for the crawler

@app.route("/json", methods=['POST'])
def json():
	input_json = request.get_json(force=True)
	return jsonify(input_json)
