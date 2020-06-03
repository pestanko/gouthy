#!/usr/bin/env bash

source "bin/config.sh"

# CREATE TEST USERS
$APP users create -u admin -e "admin@localhost" -n "Admin User" -p "$PASSWORD"
$APP users create -u admin2 -e "admin2@localhost" -n "Admin User2" -p "$PASSWORD"
