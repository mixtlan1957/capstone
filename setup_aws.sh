# Setup for AWS linux

# Make sure git is installed
sudo yum install git

# Set up Golang
sudo yum install -y golang

# Be sure to add the next 3 exports to ~/.bash_profile as well
export GOROOT=/usr/lib/golang
export GOPATH=$HOME/projects
export PATH=$PATH:$GOROOT/bin

# Move to GOPATH
cd $GOPATH

# Clone the git repository here
git clone -b midpoint-progress https://github.com/supernimbus/capstone.git
mv ./capstone/* $GOPATH

# Set up MongoDB
# Follow these instructions to setup mongodb: https://docs.mongodb.com/manual/tutorial/install-mongodb-on-amazon/
# Then run this
# sudo service mongod start

# Start the web goat
cd $GOPATH/goats/easy_goat/
bash ./goat_server.sh &

# Start the website
cd $GOPATH/src/server/
bash ./server.py
