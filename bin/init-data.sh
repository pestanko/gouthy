#!/usr/bin/env bash

source "bin/config.sh"
# CREATE STANDARD APPS

$APP apps create -c "admin_console" -n "Admin console" -C "admin_console" -D "Admin console application to manage the server" -T "internal"
$APP apps create -c "default" -n "Default application" -C "default" -D "Default application used for the password login" -T "internal"

