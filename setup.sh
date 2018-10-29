# Make sure git is installed
sudo apt install git

# Set up Golang
wget https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz
tar -xvf go1.11.1.linux.amd64.tar.gz
sudo mv go /usr/local
mkdir -p $HOME/Projects/Go
# Be sure to add the next 3 exports to ~/.bashrc as well
export GOROOT=/usr/local/go
export GOPATH=$HOME/Projects/Go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
rm -rf https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz

# Move to GOPATH
cd $HOME/Projects/Go

# Clone the git repository here
git clone -b midpoint-progress https://github.com/supernimbus/capstone.git
mv ./capstone/* $HOME/Projects/Go/

# Set up Python
sudo apt install python
sudo apt install python-pip

# Set up MongoDB
sudo apt-get install -y mongodb
mongo_db_path="grep dbpath /etc/mongodb.conf"
mongod --dbpath $mongo_db_path

# Start the web goat
cd $HOME/Projects/Go/goats/easy_goat/
bash ./goat_server.sh &

# Start the website
cd $HOME/Projects/Go/src/server/
bash ./server.py
