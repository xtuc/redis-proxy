get () {
  echo 
}

put () {
  curl -s -XPUT http://localhost:8080/$1 -d "$2" > /dev/null
}

for i in {1..110}; do
  put "item$i" "value$i"
  printf .
done
