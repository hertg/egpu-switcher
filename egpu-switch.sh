#!/bin/sh

DIR=/etc/X11
FILE=xorg.conf
FILE_EGPU=xorg.conf.egpu
FILE_INTERNAL=xorg.conf.internal

TEST=$(lspci | grep -c " VGA ")

#rm $DIR/$FILE

if([ $TEST -eq 2 ]); then
	# EGPU IS CONNECTED
	echo "egpu detected... linking $FILE_EGPU to $FILE"
	ln -sf $DIR/$FILE_EGPU $DIR/$FILE
else
	# EGPU IS NOT CONNECTED
	echo "NO egpu detected... linking $FILE_INTERNAL to $FILE"
	ln -sf $DIR/$FILE_INTERNAL $DIR/$FILE
fi