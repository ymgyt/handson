#! /usr/bin/env bash

set -xue

# QEMU=qemu-system-riscv32
QEMU=/opt/homebrew/bin/qemu-system-riscv32
# CC=/usr/local/Cellar/llvm/17.0.4/bin/clang
CC=/opt/homebrew/Cellar/llvm/17.0.4/bin/clang
OBJCOPY=/opt/homebrew/Cellar/llvm/17.0.4/bin/llvm-objcopy

CFLAGS="-std=c11 -O2 -g3 -Wall -Wextra --target=riscv32 -ffreestanding -nostdlib"


# シェルをビルド
$CC $CFLAGS -Wl,-Tuser.ld -Wl,-Map=shell.map -o shell.elf shell.c user.c common.c
$OBJCOPY --set-section-flags .bss=alloc,contents -O binary shell.elf shell.bin
$OBJCOPY -Ibinary -Oelf32-littleriscv shell.bin shell.bin.o

$CC $CFLAGS -Wl,-Tkernel.ld -Wl,-Map=kernel.map -o kernel.elf kernel.c common.c shell.bin.o

$QEMU -machine virt -bios default -nographic -serial mon:stdio --no-reboot -kernel kernel.elf