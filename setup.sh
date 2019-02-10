#!/bin/bash

# define some constant variables
TMP_DIR=/tmp/egpu-switcher
TMP_FILE=$TMP_DIR/results
TEMPLATE_FILE=./xorg.conf.template
XORG_DIR=/etc/X11
XORG_EGPU=$XORG_DIR/xorg.conf.egpu
XORG_INTERNAL=$XORG_DIR/xorg.conf.internal

# check if the script is run as root
if [ "$EUID" -ne 0 ]
  then echo "You need to run the script with root privileges"
  exit
fi

# create the tmp dir
mkdir -p $TMP_DIR

# delete existing results file
rm -f $TMP_FILE

# search for GPUs and save them to the temp file
lspci | grep VGA | while read -r line ; do
    echo "$line" >> $TMP_FILE
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

# prompt to choose the internal gpu from the list
echo "Choose your preferred INTERNAL GPU [1-$NUM_OF_RESULTS]: "
read internal
PCI_INTERNAL=$(sed ''"$internal"'q;d' $TMP_FILE | grep -Eo '^[^ ]+')

# prompt to choose the external gpu from the list
echo "Choose your preferred EXTERNAL GPU [1-$NUM_OF_RESULTS]: "
read external
PCI_EXTERNAL=$(sed ''"$external"'q;d' $TMP_FILE | grep -Eo '^[^ ]+')

# create the internal xorg config file
cp $TEMPLATE_FILE $XORG_INTERNAL
sed -i -e 's/\$BUS/PCI:'$PCI_INTERNAL'/g' -e 's/\$DRIVER/nvidia/g' -e 's/\$ID/Device0/g' $XORG_INTERNAL

# create the external xorg config file
cp $TEMPLATE_FILE $XORG_EGPU
sed -i -e 's/\$BUS/PCI:'$PCI_EXTERNAL'/g' -e 's/\$DRIVER/nvidia/g' -e 's/\$ID/Device0/g' $XORG_EGPU

