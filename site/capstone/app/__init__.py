#needed to make this subfolder into a package
from flask import Flask
app = Flask(__name__)
app.config['SECRET_KEY'] = 'HIMITSU' #required for CSRF



from app import routes