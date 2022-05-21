#!/usr/bin/env sh
OUTFILE="data/in/input.json"
N=$1
URL=""

mkdir -p data/in && mkdir -p data/out # If data dirs don't exist yet, create them

# If the argument is empty or negative, just download all records
if [ -z "$N" ] || [ "$N" -lt 0 ]
then
URL="https://replicate.npmjs.com/_all_docs?include_docs=true"
else # If N was positive and specified, download that many records
URL="https://replicate.npmjs.com/_all_docs?include_docs=true&limit=${N}"
fi
# The following script downloads the npm data dump and then projects so that we only have the rows entries
echo "Starting download..."
wget -O "$OUTFILE" "$URL" && echo "Starting preprocessing..." && sed --regexp-extended -i '1s/.+/[/;$s/.+/]/' $OUTFILE && echo "Done!"

