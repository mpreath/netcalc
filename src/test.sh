#!/bin/sh

CMD=./netcalc
IP=204.56.78.12
MASK=255.255.255.192

echo '*************** Running Test #1 ***************'
echo "$CMD $IP $MASK"
$CMD $IP $MASK

echo '*************** Running Test #2 ***************'
echo "$CMD -v -n 6 $IP $MASK"
$CMD -n 6 $IP $MASK

echo '*************** Running Test #3 ***************'
echo "$CMD -v -h 4 $IP $MASK"
$CMD -v -h 4 $IP $MASK

echo '*************** Running Test #4 ***************'
echo "$CMD -v -l 2,2,3,20 $IP $MASK"
$CMD -v -l 2,2,2 $IP $MASK

echo '*************** Running Test #5 ***************'
echo "$CMD -v -s < test.txt"
$CMD -v -s < test.txt