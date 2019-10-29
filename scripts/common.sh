#!/bin/bash

function get_ip() {
  if [[ "$OSTYPE" == "darwin"* ]]; then
    ip="localhost" 
  else
    ip=$(docker inspect $1 --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')
  fi
  echo $ip
}

function remove() {
  container="$1"
  docker rm -f $container || true
}

function launch() {
  name="$1"
  image="$2"
  env=${3:-""}
  remove $name
  container=$(docker run -d --name $name --label "detection.test=attacker" --label "com.docker.compose.service=$name" --label "app=$name" $env $image)
  echo $container
  sleep 1
}

function launch_attacker() {
  remove "attacker"
  docker run -d --name attacker --label "detection.test=attacker" -p 8080:8080 kaizheh/attacker
}

function attack() {
  attacker_ip=$(get_ip "attacker")
  target="$1"
  tool=${2:-"metasploit"}

  target_ip=$(get_ip "$target")
  curl -X POST -d "$(sed -e "s/{{.RHOST}}/$target_ip/g" ./json/$target.json)" "http://$attacker_ip:8080/attack?tool=$tool"
  sleep 1
}

function init_couchdb() {
  ip=$(get_ip "couchdb")
  sleep 30
  curl -X PUT http://admin:password@$ip:5984/_users
  curl -X PUT http://admin:password@$ip:5984/_replicator
  curl -X PUT http://admin:password@$ip:5984/_global_changes
  sleep 5
}

