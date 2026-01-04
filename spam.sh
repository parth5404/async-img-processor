#!/bin/bash

for i in {1..70}
do
  ./client_bin &
done

wait
