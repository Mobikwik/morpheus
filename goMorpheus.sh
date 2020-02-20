# Download, Install and Test Go Setup as described here: https://golang.org/doc/install 
# This script installs Go in /usr/goInstallation directory. 

echo "***************Installing GoLang 1.13.8**********************"
cd /usr
mkdir /usr/goInstallation
cd goInstallation
wget https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz
tar -C /usr/goInstallation -xzf go1.13.8.linux-amd64.tar.gz

# Make GoLang workspace directory /usr/goLangProjects
cd /usr
mkdir goLangProjects

# Setup PATH and GOPATH environment variables.
export PATH=$PATH:/usr/goInstallation/go/bin
export GOPATH=/usr/goLangProjects

# Clone Morpheus repo from github in GoLang workspace
echo "***************cloning Morpheus github repo*********************"
cd /usr/goLangProjects
git clone git@github.com:Mobikwik/morpheus.git

# Now we are ready to run Morpheus. Switch to "morpheus" directory.
cd morpheus

# Below command will skip the _test file and run Morpheus on port 8080 (port mentioned in env.properties). 
# It will also send the logs to file morpheusLog.log
echo "****************Running Morpheus*****************"
go run $(ls -t | grep -v _test | grep .go) env.properties > morpheusLog.log 2>&1 &

# Test if Morpheus has run successfully.
echo "****************Checking if Morpheus has run successfully*******************"
# Adding sleep to wait for Morpheus to start. Then hitting cURL. 
sleep 5
curl -X GET -i http://localhost:8080


# If you face below error msg while running go test:
# go test exec: "gcc": executable file not found in $PATH
# Uncomment below line to Set CGO_ENABLED=0 environment variable
# export CGO_ENABLED=0
