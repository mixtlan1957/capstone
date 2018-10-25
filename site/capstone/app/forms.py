from flask_wtf import FlaskForm
from wtforms import StringField, RadioField, IntegerField, SubmitField, validators
from wtforms.validators import DataRequired, NumberRange, ValidationError, URL, Optional
import httplib2

#custom validator for url template
class URLChecker(object):
	def __init__(self, message=None):
		if not message:
			message = 'invalid URL'
		self.message = message

	def __call__(self, form, field):

		head = httplib2.Http() 
		if not field.data.startswith("http"):  #have to prepend with http for httplib2
			field.data = "http://" + field.data 
		try:
			res = head.request(field.data, 'HEAD') #get the header
			if int(res[0]['status'] < 400):
				raise ValidationError(self.message) #check status code in header to see if it's good
		except: 
			raise ValidationError(self.message) # catch error, website is invalid

#main website form
class CrawlerForm(FlaskForm):
	url = StringField('Seed Site', default="http://", validators=[DataRequired(), URLChecker()]) #first arg is label
	traversal = RadioField('BFS/DFS', choices=[('bfs','Breadth-First'),('dfs','Depth-First')], default="bfs", validators=[DataRequired()]) #first arg is value for choices
	depth = IntegerField('Depth level', default=1, validators=[NumberRange(min=1, max=3), DataRequired()]) #limit integer field
	keyword = StringField('Keyword (optional)', validators=[Optional()]) #no validator needed here
	submit = SubmitField('Submit')

	#these var names are the "name=" attribute for the form data