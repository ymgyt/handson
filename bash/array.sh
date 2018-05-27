#!/bin/bash

ary=(aaa bbb ccc)
echo ${ary[@]}

ary=("${ary[@]}" ddd eee)
echo ${ary[@]}

ary=(xxx yyy "${ary[@]}")
echo ${ary[@]}
