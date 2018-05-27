#!/bin/bash

str1=xyz
patern='x*'
if [[ $str1 == $patern ]];then
  echo YES
else
  echo NO
fi
  
str1=/home/yuta
if [[ $str1 =~ ^/home/[^/]+$ ]]; then
  echo YES
else
  echo NO
fi

