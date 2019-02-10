#!/bin/sh

# define constand variables
DIR=/etc/X11
FILE=xorg.conf
FILE_EGPU=xorg.conf.egpu
FILE_INTERNAL=xorg.conf.internal
EGPU_PCI_ID=$(cat $DIR/$FILE_EGPU | grep -Ei "BusID" | grep -oEi '[0-9]+\:[0-9]+\:[0-9]+')
TEST=$(nvidia-xconfig --query-gpu-info | grep -c $EGPU_PCI_ID)

# check if the PCI of the egpu is listed in the "lspci" output
if([ $TEST -eq 1 ]); then
	# egpu is connected
	echo "egpu detected... linking $FILE_EGPU to $FILE"
	ln -sf $DIR/$FILE_EGPU $DIR/$FILE
else
	# egpu is not connected
	echo "NO egpu detected... linking $FILE_INTERNAL to $FILE"
	ln -sf $DIR/$FILE_INTERNAL $DIR/$FILE
fi

# systemctl restart display-manager.service