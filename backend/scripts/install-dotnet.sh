# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-dotnet.lock

wget -q https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
dpkg -i packages-microsoft-prod.deb
apt-get install -y apt-transport-https
apt-get update
apt-get install -y dotnet-sdk-2.1 dotnet-runtime-2.1 aspnetcore-runtime-2.1

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-dotnet.lock
