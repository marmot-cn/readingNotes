#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

int main(int argc, char *argv[])
{
    printf("hello world (pid:%d)\n", (int) getpid());
    
    int fd;
    fd = open("/tmp/temp", O_WRONLY|O_CREAT);

    int rc = fork();
    if (rc < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (rc == 0) {
        // child (new process)
        printf("in child fd is :%d\n", fd);

        char s[] = "child!\n";
        write(fd, s, sizeof(s));
        close(fd);

        printf("hello, I am child (pid:%d)\n", (int) getpid());
    } else {
        // parent goes down this path (original process)

        printf("in parent fd is :%d\n", fd);

        char s[] = "parent!\n";
        write(fd, s, sizeof(s));
        close(fd);

        printf("hello, I am parent of %d (pid:%d)\n",
	       rc, (int) getpid());
    }
    return 0;
}