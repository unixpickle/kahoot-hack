apt update
apt install tmux
apt install git
apt install golang
echo "export GOPATH=/data/data/com.termux/files/usr" > tmp1
chmod +x tmp1
sleep 1
source tmp1
echo "Downloading kahoot-hack... Please wait"
go get github.com/unixpickle/kahoot-hack
echo "Downloading websocket... Please wait"
go get github.com/gorilla/websocket
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-crash/main.go ~/crash.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-flood/main.go ~/flood.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-html/main.go ~/html.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-play/main.go ~/play.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-profane/main.go ~/profane.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-rand/main.go ~/rand.go
mv ../usr/src/github.com/unixpickle/kahoot-hack/kahoot-xss/main.go ~/xss.go
mkdir ../usr/etc/profile.d
echo "export GOPATH=/data/data/com.termux/files/usr" > ../usr/etc/profile.d/start.sh
chmod +x ../usr/etc/profile.d/start.sh
rm tmp1
rm a.sh
echo "Installation Successful. Please restart the app by typing 'exit' and open the app again."