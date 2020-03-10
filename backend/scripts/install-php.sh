# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-php.lock

apt-get install php

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-php.lock
