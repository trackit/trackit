#!/bin/bash

echo "===> checking vm.max_map_count setting"

min_mmap=262144
export PATH=$PATH:/sbin
current_mmap=$(sysctl -n vm.max_map_count)
if [[ -n $current_mmap ]] && [[ $current_mmap -lt $min_mmap ]]
then
  echo "your virtual memory max map count is too low to allow elasticsearch to run properly, please run the following command before restarting the installation:"
  echo "sudo sysctl -w vm.max_map_count=262144"
  exit 1
fi
exit 0
