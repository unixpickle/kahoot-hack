apt update
apt install clang git golang
export GOPATH=/data/data/com.termux/files/usr
printf "Downloading kahoot-hack ..."
go get github.com/unixpickle/kahoot-hack
printf "\nDownloading websocket ..."
go get github.com/gorilla/websocket
printf "\nProcessing ...\n"
mkdir /data/data/com.termux/files/usr/etc/profile.d
echo "export GOPATH=/data/data/com.termux/files/usr" >> /data/data/com.termux/files/usr/etc/profile.d/kh-conf.sh
mkdir /data/data/com.termux/files/usr/var/kahoot-hack
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-crash/main.go /data/data/com.termux/files/usr/var/kahoot-hack/crash.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-flood/main.go /data/data/com.termux/files/usr/var/kahoot-hack/flood.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-html/main.go /data/data/com.termux/files/usr/var/kahoot-hack/html.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-play/main.go /data/data/com.termux/files/usr/var/kahoot-hack/play.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-profane/main.go /data/data/com.termux/files/usr/var/kahoot-hack/profane.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-rand/main.go /data/data/com.termux/files/usr/var/kahoot-hack/rand.go
mv /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/kahoot-xss/main.go /data/data/com.termux/files/usr/var/kahoot-hack/xss.go
gcc /data/data/com.termux/files/usr/src/github.com/unixpickle/kahoot-hack/Android/kahoot-hack.c -o /data/data/com.termux/files/usr/bin/kahoot-hack
rm /data/data/com.termux/files/home/kh-tmux-setup.sh
clear
printf "Installation Successful.\nTo use the hack, type 'kahoot-hack' in the command line.\n"
