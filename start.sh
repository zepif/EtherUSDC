#!/bin/sh

until ping -c1 db >/dev/null 2>&1; do
	echo "Waiting for database..."
	sleep 1
done

echo "Database is up, starting migrations..."
EtherUSDC migrate up

echo "Starting service..."
EtherUSDC run service
