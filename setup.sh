# Set up Golang
wget https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz
sudo tar -xvf go1.11.1.linux.amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/usr/local/go
mkdir -p $HOME/Projects/Go
export GOPATH=$HOME/Projects/Go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH

# Move to GOPATH
cd $HOME/Projects/Go

# Clone the git repository here
git clone -b midpoint-progress https://github.com/supernimbus/capstone.git

# Set up Python
sudo apt install python
sudo apt install python-pip
pip install --upgrade -r requirements.txt

# Set up MongoDB
sudo apt-get install -y mongodb
mongo_db_path="grep dbpath /etc/mongodb.conf"
mongod --dbpath $mongo_db_path

# Start the web goat
bash ./goats/easy_goat/goat_server.sh

# Start the website


# Start the crawler
#bash ./src/crawl.sh
