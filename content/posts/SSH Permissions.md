---
title: SSH Permissions
tags: [ssh, permissions]
date: 07 Nov 2023
---

# How To Set Correct SSH Directory Permissions in Linux

### Set Correct SSH Directory Permissions in Linux

If you ever encounter errors while trying to SSH into a server, you can set correct ssh directory permissions on the **.ssh** directory using the **chmod** command.

```
# chmod u+rwx,go-rwx ~/.ssh
OR
# chmod  0700 ~/.ssh

```

If ssh complains of wrong permissions on any of the above files, you can set the correct permissions for any of the files like this:

```
# chmod u+rw,go-rwx .ssh/id_rsa
# chmod u+rw,go-rwx .ssh/id_rsa.pub
# chmod u+rw,go-rwx .ssh/authorized_keys
# chmod u+rw,go-rwx .ssh/known_hosts
# chmod u+rw,go-rwx .ssh/config
OR
# chmod 600 .ssh/id_rsa
# chmod 600 .ssh/id_rsa.pub
# chmod 600 .ssh/authorized_keys
# chmod 600 .ssh/known_hosts
# chmod 600 .ssh/config

```

To remove write permissions for group and others on the home directory, run this command:

```
# chmod go-w ~
OR
# chmod 755 ~

```

---

# References:
https://www.tecmint.com/set-ssh-directory-permissions-in-linux/
