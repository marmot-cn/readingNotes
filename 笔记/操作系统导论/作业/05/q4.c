#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>


int main(int argc, char *argv[])
{
    printf("begin fork\n");
    int x = 100;
    int rc = fork();
    if (rc < 0) {
        // fork failed; exit
        fprintf(stderr, "fork failed\n");
        exit(1);
    } else if (rc == 0) {
        printf("in child pid is %d\n", (int) getpid());
        printf("before run exec \n");

        //成功返回0,失败返回-1
        //带l 的exec函数：execl,execlp,execle，表示后边的参数以可变参数的形式给出且都以一个空指针结束
        //不带 l 的exec函数：execv,execvp表示命令所需的参数以char *arg[]形式给出且arg最后一个元素必须是 Null
        //带 p 的exec函数：execlp,execvp，表示第一个参数path不用输入完整路径，只有给出命令名即可
        //带 e 的exec函数：execle表示，将环境变量传递给需要替换的进程

        //int execl(const char *path, const char *arg, ...);
        // execl("/bin/ls","ls","-l",NULL);

        // int execlp(const char *file, const char *arg, ...);
        execlp("ls", "ls", "-l", "-h", NULL);

        // int execle(const char *path, const char *arg, ..., char * const envp[]);
        // char * const envp[] = {""};
        // execle("/bin/ls", "ls", "-l", NULL, envp);

        // int execv(const char *path, char *const argv[]);
        // printf("execv\n ");
        // char *argv[] = {"ls","-l",NULL};
        // execv("/bin/ls", argv);

        // int execvp(const char *file, char *const argv[]);
        // printf("execvp\n ");
        // char *argv[] = {"ls","-l",NULL};
        // execvp("ls", argv);

        //mac 并不支持该函数
        //int execvpe(const char *file, char *const argv[],char *const envp[]);
        // printf("execvpe\n ");
        // char * const envp[] = {""};
        // char *argv[] = {"ls", "-l", NULL};
        // execvpe("ls", argv, envp);

        //这句话不会输出，因为exec系列的函数会将当前进程替换掉
        printf("after run exec ");
    } 

    return 0;
}