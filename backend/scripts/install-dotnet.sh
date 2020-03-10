# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-dotnet.lock

apt-get install curl
curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
mv microsoft.gpg /etc/apt/trusted.gpg.d/microsoft.gpg
sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-ubuntu-artful-prod artful main" > /etc/apt/sources.list.d/dotnetdev.list'
apt-get install apt-transport-https
apt-get update
apt-get install dotnet-sdk-2.1

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-dotnet.lock
