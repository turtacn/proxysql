export LD_LIBRARY_PATH=/usr/local/mysql/lib/mysql/
DBNAME=$1
WH=$2
HOST=127.0.0.1
STEP=100

./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 1 -m 1 -n $WH >> 1.out &

x=1

while [ $x -le $WH ]
do
 echo $x $(( $x + $STEP - 1 ))
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 2 -m $x -n $(( $x + $STEP - 1 ))  >> 2_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 3 -m $x -n $(( $x + $STEP - 1 ))  >> 3_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 4 -m $x -n $(( $x + $STEP - 1 ))  >> 4_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 5 -m $x -n $(( $x + $STEP - 1 ))  >> 5_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 6 -m $x -n $(( $x + $STEP - 1 ))  >> 6_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 7 -m $x -n $(( $x + $STEP - 1 ))  >> 7_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 8 -m $x -n $(( $x + $STEP - 1 ))  >> 8_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 9 -m $x -n $(( $x + $STEP - 1 ))  >> 9_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 10 -m $x -n $(( $x + $STEP - 1 ))  >> 10_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 11 -m $x -n $(( $x + $STEP - 1 ))  >> 11_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 12 -m $x -n $(( $x + $STEP - 1 ))  >> 12_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 13 -m $x -n $(( $x + $STEP - 1 ))  >> 13_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 14 -m $x -n $(( $x + $STEP - 1 ))  >> 14_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 15 -m $x -n $(( $x + $STEP - 1 ))  >> 15_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 16 -m $x -n $(( $x + $STEP - 1 ))  >> 16_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 17 -m $x -n $(( $x + $STEP - 1 ))  >> 17_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 18 -m $x -n $(( $x + $STEP - 1 ))  >> 18_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 19 -m $x -n $(( $x + $STEP - 1 ))  >> 19_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 20 -m $x -n $(( $x + $STEP - 1 ))  >> 20_$x.out &
./tpcc_load -h $HOST -d $DBNAME -u root -p "" -w $WH -l 21 -m $x -n $(( $x + $STEP - 1 ))  >> 21_$x.out &
 x=$(( $x + $STEP ))
done