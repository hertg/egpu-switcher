#!/bin/bash

# define some constant variables
TMP_DIR=/tmp/egpu-switcher
TMP_FILE=$TMP_DIR/results
TEMPLATE_FILE=./xorg.conf.template
XORG_DIR=/etc/X11
XORG_EGPU=$XORG_DIR/xorg.conf.egpu
XORG_INTERNAL=$XORG_DIR/xorg.conf.internal

INITD_TEMPLATE=./egpu-switch.sh
INITD=/etc/init.d/egpu-switch
INITD_SYMLINK=/etc/rc5.d/S15egpu-switch

# check if the script is run as root
if [ "$EUID" -ne 0 ]
  then echo "You need to run the script with root privileges"
  exit
fi

# check if the template/script files can be found
if [ ! -f $TEMPLATE_FILE ]; then
    echo "The xorg.conf template file could not be found."
    echo "Have you run the setup script in the same order where its located?"
    exit
fi

if [ ! -f $INITD_TEMPLATE ]; then
    echo "The egpu-switch script file could not be found."
    echo "Have you run the setup script in the same order where its located?"
    exit
fi

# create the tmp dir
mkdir -p $TMP_DIR

# delete existing results file
rm -f $TMP_FILE

# search for GPUs and save them to the temp file
nvidia-xconfig --query-gpu-info | grep -i -e 'gpu #[0-9]' | while read -r line ; do
    
    bus=$(nvidia-xconfig --query-gpu-info | grep "$line" -A 3 | awk '/PCI BusID/{print$4}')
    name=$(nvidia-xconfig --query-gpu-info | grep "$line" -A 3 | grep -oP 'Name\s+\:\s+\K.*')
    echo "$name ($bus)" >> $TMP_FILE
done

# save the number of lines to a variable
NUM_OF_RESULTS=$(wc -l < "/tmp/egpu-switcher/results")

# additional check
if [ $NUM_OF_RESULTS -lt "2" ]; then
    echo "Only $NUM_OF_RESULTS GPUs found, there need to be at least 2."
    echo "Make sure to connect your EGPU for the setup."
    exit
fi

# print the GPUs found
echo "Found $NUM_OF_RESULTS possible GPUs..."
echo ""

i=0
while read -r line; do
    i=$((i+1))
    echo "  $i: $line"
done < "$TMP_FILE"

echo ""

# prompt to choose the internal gpu from the listnvidia-xconfig --query-gpu-info | grep $line -A 3
echo "Choose your preferred INTERNAL GPU [1-$NUM_OF_RESULTS]: "
read internal
PCI_INTERNAL=$(sed ''"$internal"'q;d' $TMP_FILE | grep -Eo 'PCI\:[0-9]+\:[0-9]+\:[0-9]+')

# prompt to choose the external gpu from the list
echo "Choose your preferred EXTERNAL GPU [1-$NUM_OF_RESULTS]: "
read external
PCI_EXTERNAL=$(sed ''"$external"'q;d' $TMP_FILE | grep -Eo 'PCI\:[0-9]+\:[0-9]+\:[0-9]+')

# create the internal xorg config file
cp $TEMPLATE_FILE $XORG_INTERNAL
sed -i -e 's/\$BUS/'$PCI_INTERNAL'/g' -e 's/\$DRIVER/nvidia/g' -e 's/\$ID/Device0/g' $XORG_INTERNAL

# create the external xorg config file
cp $TEMPLATE_FILE $XORG_EGPU
sed -i -e 's/\$BUS/'$PCI_EXTERNAL'/g' -e 's/\$DRIVER/nvidia/g' -e 's/\$ID/Device0/g' $XORG_EGPU

# setup startup script
if [ -f $INITD ]; then
    echo "The script $INITD does already exist. Removing it ..."
    rm -f $INITD_SYMLINK
    rm -f $INITD
fi

cp $INITD_TEMPLATE $INITD
ln -sf $INITD $INITD_SYMLINK

if [ -f $INITD ]; then
    echo "The file $INITD was created"
else
    echo "Something went wrong while creating the $INITD file"
    exit
fi


if [ -f $INITD_SYMLINK ]; then
    echo "The file $INITD_SYMLINK was created"
else
    echo "Something went wrong while creating the $INITD_SYMLINK file"
    exit
fi

echo "Script finished"