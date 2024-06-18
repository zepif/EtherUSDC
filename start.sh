#!/bin/sh

until nc -z db 5432; do
	echo "Waiting for database..."
	sleep 1
done

EtherUSDC migrate up
EtherUSDC run service
