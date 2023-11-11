#! /usr/bin/env bash

set -xue

QEMU=qemu-system-riscv32
CC=clang

CFLAGS="-std=c11 -O2 -g3 -Wall -Wextra --target=riscv32 -ffreestanding -nostdlib"

MYL="--ld-path=/nix/store/1jfx3fx1hvdvsahhnpjppj1cpm956g08-clang-wrapper-16.0.6/bin/ld"

# build kernel
$CC $CFLAGS $MYL -Wl,-Tkernel.ld -Wl,-Map=kernel.map -o kernel.elf kernel.c

$QEMU -machine virt -bios default -nographic -serial mon:stdio --no-reboot