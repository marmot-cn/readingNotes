#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
   
    int x = 100;
    int rc = fork();

    //pid_t wait(int * status);

    if (rc < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (rc == 0) {
        // child (new process)
        printf("child (pid:%d)\n", (int) getpid());
        pid_t pid;
        int status, i;
        pid = wait(&status);
        //用宏 WEXITSTATUS 来提取子进程的返回值
        printf("child wait return is %d, WIFEXITED(status) == %d\n", (int)pid, WEXITSTATUS(status));
    } else {
        // parent goes down this path (original process)
        printf("parent (pid:%d)\n", (int) getpid());
        
        pid_t pid;
        pid = wait(NULL);
        printf("parent wait return is %d\n", (int)pid);
    }
    return 0;
}