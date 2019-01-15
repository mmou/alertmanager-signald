# alertmanager-signald

Alertmanager webhook receiver that sends Signal messages.

# Docker notes
As written, the Dockerfile requires our signald image, which is of a modified
version of [signald](https://github.com/thefinn93/signald) that uses TCP
instead of UNIX sockets. 

# Run notes
You need to set the values for 3 env vars:
- `SIGNALD_BIND_ADDR`: e.g. 0.0.0.0:8888
- `SENDER_NUMBER`
- `RECEIVER_GROUP_IP`

