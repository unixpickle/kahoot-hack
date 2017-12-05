apt-get update
apt-get install tmux
apt-get install git
apt-get install golang
echo "export GOPATH=/data/data/com.termux/files/usr" > and_kh_tmp1
chmod +x and_kh_tmp1
sleep 1
source and_kh_tmp1
echo "Downloading kahoot-hack... Please wait"
go get github.com/unixpickle/kahoot-hack
echo "Downloading websocket... Please wait"
go get github.com/gorilla/websocket
echo "Downloading gopass... Please wait"
go get github.com/howeyc/gopass
mkdir ~/kahoot
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-auto/main.go ~/kahoot/auto.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-crash/main.go ~/kahoot/crash.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-flood/main.go ~/kahoot/flood.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-html/main.go ~/kahoot/html.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-play/main.go ~/kahoot/play.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-profane/main.go ~/kahoot/profane.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-rand/main.go ~/kahoot/rand.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-xss/main.go ~/kahoot/xss.go
mkdir /data/data/com.termux/files/usr/etc/profile.d
echo "export GOPATH=/data/data/com.termux/files/usr" >> /data/data/com.termux/files/usr/etc/profile.d/kahoot-hack-config.sh
chmod +x /data/data/com.termux/files/usr/etc/profile.d/kahoot-hack-config.sh
rm and_kh_tmp1
echo "Installation Successful. Please restart the app by typing 'exit' and open the app again."
