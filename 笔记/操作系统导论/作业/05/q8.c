#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
   
    //int pipe(int filedes[2]);
    //ssize_t read(int fd,void*buf,size_t count)
    //ssize_t write(int fd,void*buf,size_t count)

    int fd[2], nbytes;
    int result, processWrite, processRead;

    int *write_fd = &fd[1];
    int *read_fd = &fd[0];

    char output[] = "hello world, pipe";
    char input[100];

    // printf("strlen is %d", sizeof(output));

    result=pipe(fd);
    if(result == -1) {
        printf("fail to create pipe\n");
        return -1;
    }

    processWrite=fork();

    if (processWrite < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (processWrite == 0) {
        //write
        printf("write\n");
        close(*read_fd);
        printf("write %s \n", output);

        int ret;
        ret = write(*write_fd, output, sizeof(output));
        if (ret < 0 ) {
            printf("write wrong\n");
            return 0;
        }
        return 0;
    } else {
        wait(NULL);        
    }

    processRead=fork();
    if (processRead < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (processRead == 0) {
        //read
        printf("read\n");
        close(*write_fd);

        nbytes=read(*read_fd, input, sizeof(input)-1);
        printf("receive %d bytes %s \n", nbytes, input);
        return 0;
    }

    return 0;
}