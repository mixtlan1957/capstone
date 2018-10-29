from app import app
from flask import render_template, flash, redirect, request, jsonify #for rendering html templates, request for dealing with POST
from pymongo import MongoClient #for MongoDB / python
from app.forms import CrawlerForm

import requests #for sending POST requests to crawler and graph
import subprocess #to call crawler, starts from capstone folder as main
import pprint
import json
import ast

"""
database, pull crawl data from DB and send to graphical UI
"""
client = MongoClient() #localhost, default port: 7018
db = client.crawlResults #creates a new database



"""
routes
"""

@app.route('/', methods=['GET', 'POST'])
def index():
	form = CrawlerForm()
	if form.validate_on_submit():
		url = request.form.get('url')
		if not url.startswith("http"): #to correct data for crawler
			print url
			url = "http://" + url
			print url
		if not url.endswith("/"):
			print url
			url += "/"
			print url

		traversal = request.form.get('traversal')
		depth = request.form.get('depth')
		keyword = request.form.get('keyword') #need to use .get here because an optional key

		# data to send to crawler
		# data = {"url": url, "traversal": traversal, "depth": depth, "keyword": keyword}

		url_flag = "-url=" + url
		traversal_flag = "-search=" + traversal

		# res = requests.post('http://localhost:12345/crawl', json=data)
		# cat = subprocess.check_output(["ls"])

		db_id = subprocess.check_output(["go", "run", "../main.go", url_flag, traversal_flag])

		print db_id
		db_id = str(db_id)

		if traversal == "bfs":
			crawl_res = db.bfsCrawl.find({'crawlid': db_id}) #have to iterate otherwise get cursor obj
		else:
			crawl_res = db.dfsCrawl.find({'crawlid': db_id})

		# pprint.pprint(crawl_res)

		# pprint.pprint(crawl_res[0])
		# print(crawl_res[0])

		formatted_json = []

		node_counter = 0

		for node in crawl_res[0]['linkdata']:

			# print type(data)
			new = {}
			new['name'] = str(node['url'])
			new['children'] = [ str(c) for c in node['childlinks'] ]
			new['xss'] = "false" #just because not in the linkNode struct yet
			new['sqli'] = "false" #just because not in the linkNode struct yet

			node_counter += 1
			formatted_json.append(new)

		print node_counter
		# print formatted_json (no 'u')

		res = requests.post('https://capstone-graphics-portion.herokuapp.com/graphs', json=formatted_json)
		return redirect('https://capstone-graphics-portion.herokuapp.com/')


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

