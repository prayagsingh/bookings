#!/bin/bash

# run booking iff go build is successful
go build -o bookings cmd/web/*.go && ./bookings