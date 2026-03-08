#include <elf.h>      // ELF の型や定数。次の段階で Elf_Ehdr などを使うため。
#include <fcntl.h>    // open() と O_RDONLY を使うため。
#include <stdio.h>    // printf() と fprintf() を使うため。
#include <stdlib.h>   // 今回は未使用。参考書どおりで、後で exit() を使う段階で必要。
#include <sys/mman.h> // 今回は未使用。後で mmap() / munmap() を使う段階で必要。
#include <sys/stat.h> // struct stat と fstat() を使うため。
#include <unistd.h>   // close() を使うため。

#define IS_ELF(ehdr) \
    ((ehdr).e_ident[EI_MAG0] == ELFMAG0 && \
    (ehdr).e_ident[EI_MAG1] == ELFMAG1 && \
    (ehdr).e_ident[EI_MAG2] == ELFMAG2 && \
    (ehdr).e_ident[EI_MAG3] == ELFMAG3)

int main(int argc, char *argv[])
{
    int fd, i;
    struct stat sb;
    char * head;
    Elf64_Ehdr *ehdr;
    Elf64_Shdr *shdr;
    Elf64_Shdr *shstr;

    if (argc != 2) {
        fprintf(stderr, "usage: %s <elf-file>\n", argv[0]);
        return 1;
    }

    fd = open(argv[1], O_RDONLY);
    if (fd == -1) {
        perror("open");
        return 1;
    }

    if (fstat(fd, &sb) == -1) {
        perror("fstat");
        close(fd);
        return 1;
    }

    head = mmap(NULL, sb.st_size, PROT_READ, MAP_SHARED, fd, 0);

    ehdr = (Elf64_Ehdr *)head;

    if (!IS_ELF(*ehdr)) {
        fprintf(stderr, "This is not ELF file. (%s)\n", argv[1]);
        return 1;
    }

    shstr = (Elf64_Shdr *)(head + ehdr->e_shoff + ehdr->e_shentsize * ehdr->e_shstrndx);

    for (i=0; i < ehdr->e_shnum; i++) {
         shdr = (Elf64_Shdr *)(head + ehdr->e_shoff + ehdr->e_shentsize * i);
         printf("%s\n", (char *)(head + shstr->sh_offset + shdr->sh_name));       
    }

    munmap(head, sb.st_size);
    if (close(fd) == -1) {
        perror("close");
        return 1;
    }

    return 0;
}
