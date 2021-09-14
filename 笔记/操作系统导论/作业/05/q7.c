#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
   
    int x = 100;
    int rc = fork();

    //pid_t waitpid(pid_t pid, int *status, int options);

    if (rc < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (rc == 0) {
        // child (new process)
        printf("child (pid:%d)\n", (int) getpid());
        printf("STDOUT_FILENO is %d\n", (int)STDOUT_FILENO);
        printf("before close STDOUT_FILENO \n");
        close(STDOUT_FILENO);
        printf("after close STDOUT_FILENO \n");
    } else {
        // parent goes down this path (original process)
        printf("parent (pid:%d)\n", (int) getpid());
    }
    return 0;
}