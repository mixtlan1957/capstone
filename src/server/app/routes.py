from app import app
from flask import render_template, session, redirect, request, jsonify, make_response #for rendering html templates, request for dealing with POST
from pymongo import MongoClient #for MongoDB / python
from app.forms import CrawlerForm

import requests #for sending POST requests to crawler and graph
import subprocess #to call crawler, starts from capstone folder as main
import pprint
import json
import ast
import copy #for deepcopy

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
		if not url.startswith("http") and not url.startswith("https"): #to correct data for crawler
			url = "http://" + url
		if not url.endswith("/"):
			url += "/"

		traversal = request.form.get('traversal')
		depth = request.form.get('depth')
		keyword = request.form.get('keyword') #need to use .get here because an optional key

		#setting session data
		# options = [ url, traversal, depth, keyword ]
		# session['options'] = options

		# data to send to crawler
		# data = {"url": url, "traversal": traversal, "depth": depth, "keyword": keyword}

		url_flag = "-url=" + url
		traversal_flag = "-search=" + traversal

		depth_flag =  "-depth=" + depth
		keyword_flag = "-keyword=" + keyword
		
		#use subprocess library to directly call crawler

		db_id = subprocess.check_output(["go", "run", "../main.go", url_flag, traversal_flag, depth_flag, keyword_flag])

		print db_id
		db_id = str(db_id)

		#find crawl result in database
		if traversal == "bfs":
			crawl_res = db.bfsCrawl.find({'crawlid': db_id}) #have to iterate otherwise get cursor obj
		else:
			crawl_res = db.dfsCrawl.find({'crawlid': db_id})

		# pprint.pprint(crawl_res[0])

		formatted_json = []

		node_counter = 0

		#reformat data to match cytoscape

		for node in crawl_res[0]['linkdata']:

			# print type(data)
			new = {}
			new['name'] = str(node['url'])
			new['children'] = [ str(c) for c in node['childlinks'] ]
			new['xss'] = node['xssvulnerable']
			new['sqli'] = node['sqlivulnerable']
			new['isCrawlRoot'] = node['iscrawlroot']
			new['keyword'] = node['haskeyword']
			new['title'] = node['title']
			
			node_counter += 1
			formatted_json.append(new)

		print node_counter
		print formatted_json #(no 'u')

		res = requests.post('https://capstone-graphics-portion.herokuapp.com/graphs', json=formatted_json)
		
		#cookie set-up

		resp = make_response(redirect('https://capstone-graphics-portion.herokuapp.com/'))

		resp.set_cookie('url', url) #save options in a cookie
		resp.set_cookie('traversal', traversal)
		resp.set_cookie('depth', depth)
		resp.set_cookie('keyword', keyword)

		return resp #adapted from https://www.tutorialspoint.com/flask/flask_cookies.htm

	# top 10 results setup
	# get top 10 from each crawl type

	crawlers = list(db.bfsCrawl.find().limit(10).sort("timestamp", -1)) #list() to turn cursor into list, sorts descending by time stamp
	crawlers2 = list(db.dfsCrawl.find().limit(10).sort("timestamp", -1)) #sorts descending by time stamp

	merged_crawlers = crawlers + crawlers2 #combine the lists for culling

	# print merged_crawlers

	culled_crawlers = sorted(merged_crawlers, key=lambda x: x['timestamp'])

	culled_crawlers = culled_crawlers[-10:] #get only top 10
	culled_crawlers.reverse() #so that most recent is on the top

	return render_template('index.html', form=form, data=culled_crawlers)

@app.route("/crawls/<id>")
def crawls(id):

	if id.startswith("bfs"):
		crawl_res = db.bfsCrawl.find({'crawlid': id}) #have to iterate otherwise get cursor obj
	else:
		crawl_res = db.dfsCrawl.find({'crawlid': id})

	formatted_json = []

	node_counter = 0

	#reformat data to match cytoscape

	for node in crawl_res[0]['linkdata']:

		new = {}
		new['name'] = str(node['url'])
		new['children'] = [ str(c) for c in node['childlinks'] ]
		new['xss'] = node['xssvulnerable']
		new['sqli'] = node['sqlivulnerable']
		new['isCrawlRoot'] = node['iscrawlroot']
		new['keyword'] = node['haskeyword']
		new['title'] = node['title']
		
		node_counter += 1
		formatted_json.append(new)

	print node_counter
	print formatted_json #(no 'u')

	res = requests.post('https://capstone-graphics-portion.herokuapp.com/graphs', json=formatted_json)
	
	return redirect('https://capstone-graphics-portion.herokuapp.com/')

@app.route("/about_us")
def aboutUs():

	# session check
	# if 'options' in session:
	# 	for option in session['options']:
	# 		print "session" + option

	# cookie check
	# print "cookie" + request.cookies.get('url')
	# print "cookie" + request.cookies.get('traversal')
	# print "cookie" + request.cookies.get('depth')		
	# print "cookie" + request.cookies.get('keyword')

	return render_template('links.html')

@app.route("/faq")
def faq():
	return "Frequently Asked Questions"