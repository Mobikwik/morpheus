# Download, Setup and Test Go as described here: https://golang.org/doc/install in directory /home/goInstallation. 
cd /home
mkdir /home/goInstallation
cd goInstallation
wget https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz
tar -C /home/goInstallation -xzf go1.13.8.linux-amd64.tar.gz

# Make GoLang workspace directory goLangProjects
cd /home
mkdir goLangProjects

# Setup PATH and GOPATH environment variables.
export PATH=$PATH:/home/goInstallation/go/bin
export GOPATH=/home/goLangProjects

# Clone Morpheus repo from github in GoLang workspace
cd /home/goLangProjects
git clone git@github.com:Mobikwik/morpheus.git

# Now we are ready to run Morpheus. Switch to "morpheus" directory.
cd src/github.com/Mobikwik/morpheus

# Below command will skip the _test file and run Morpheus on port 8080 (port mentioned in env.properties). 
# It will also send the logs to file /tmp/morpheusLog.log
go run $(ls -t | grep -v _test | grep .go) env.properties > /tmp/morpheusLog.log 2>&1 &

# Test if Morpheus has run successfully.
curl -X GET -i http://localhost:8080


# If you face below error msg while running go test:
# go test exec: "gcc": executable file not found in $PATH
# Set below environment variable
export CGO_ENABLED=0

