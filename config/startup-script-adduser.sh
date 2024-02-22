#!/bin/bash
export NEW_USER=#{NEW_USER}
export NEW_USER_PASSWD=#{NEW_USER_PASSWD}
export SSH_PUB_KEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6zdgFiLr1uAK7PdcchDd+LseA5fEOcxCCt7TLlr7Mx6h8jUg+G+8L9JBNZuDzTZSF0dR7qwzdBBQjorAnZTmY3BhsKcFr8Gt4KMGrS6r3DNmGruP8GORvegdWZuXgASKVpXeI7nCIjRJwAaK1x+eGHwAWO9Z8ohcboHbLyffOoSZDSIuk2kRIc47+ENRjg0T6x2VRsqX27g6j4DfPKQZGk0zvXkZaYtr1e2tZgqTBWqZUloMJK8miQq6MktCKAS4VtPk0k7teQX57OGwD6D7uo4b+Cl8aYAAwhn0hc0C2USfbuVHgq88ESo2/+NwV4SQcl3sxCW21yGIjAGt4Hy7J $NEW_USER@localhost.localdomain"

sudo adduser -U -m "$NEW_USER"
echo "$NEW_USER:$NEW_USER_PASSWD" | chpasswd
sudo mkdir /home/"$NEW_USER"/.ssh
# shellcheck disable=SC2024
sudo echo "$SSH_PUB_KEY" > /home/"$NEW_USER"/.ssh/authorized_keys
sudo chown -R "${NEW_USER}": /home/"$NEW_USER"/.ssh