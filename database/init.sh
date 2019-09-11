#!/bin/sh

PGPASSWORD=testpw psql -h localhost -U buildacoin buildacoin -f database/init.psql
