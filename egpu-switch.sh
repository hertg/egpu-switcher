#!/bin/sh

# define constand variables
DIR=/etc/X11
FILE=$DIR/xorg.conf
FILE_EGPU=$DIR/xorg.conf.egpu
FILE_INTERNAL=$DIR/xorg.conf.internal
EGPU_PCI_ID=$(cat $FILE_EGPU | grep -Ei "BusID" | grep -oEi '[0-9]+\:[0-9]+\:[0-9]+')
TEST=$(nvidia-xconfig --query-gpu-info | grep -c $EGPU_PCI_ID)
DATETIME=$(date '+%Y%m%d%H%M%S');


# Check if there is a xorg.conf file, and back it up
if [ -f $FILE ] && ! [ -L $FILE ]; then
	echo "The $FILE file already exists. Saving a backup to $FILE.backup.$DATETIME"
	cp "$FILE" "$FILE.backup.$DATETIME"
fi

# check if the PCI of the egpu is listed in the "nvidia-xconfig --query-gpu-info" output
if([ $TEST -eq 1 ]); then
	# egpu is connected
	echo "egpu detected... linking $FILE_EGPU to $FILE"
	ln -sf $FILE_EGPU $FILE
else
	# egpu is not connected
	echo "NO egpu detected... linking $FILE_INTERNAL to $FILE"
	ln -sf $FILE_INTERNAL $FILE
fi

# systemctl restart display-manager.service