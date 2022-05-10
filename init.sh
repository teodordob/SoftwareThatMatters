#!/usr/bin/env sh
OUTFILE="data/input.json"
N=$1
URL=""

if [ ! -d "data" ] # If data dir doesn't exist yet, create it
then
mkdir data && mkdir data/out
fi

# If the argument is empty or negative, just download all records
if [ -z "$N" ] || [ $N -lt 0 ]
then
URL="https://replicate.npmjs.com/_all_docs?include_docs=true"
else # If N was positive and specified, download that many records
URL="https://replicate.npmjs.com/_all_docs?include_docs=true&limit=${N}"
fi
# The following script downloads the npm data dump and then projects so that we only have the rows entries
echo "Starting download..."
curl $URL > $OUTFILE && echo "Starting preprocessing..." && sed --regexp-extended -i 's/\{"total_rows":.*,"offset":.*,"rows":\[/[/;s/\]}$/\]/' $OUTFILE && echo "Done!"