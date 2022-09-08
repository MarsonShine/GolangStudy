#ifndef CLIBRARY_H
#define CLIBRARY_H

typedef void (*StartCallbackFn)(void* user_data, int i);
typedef void (*EndCallbackFn)(void* user_data, int a, int b);

typedef struct {
    StartCallbackFn start; // start 回调go函数
    EndCallbackFn end;// end 回调go函数
} Callbacks;


void traverse(char* filename, Callbacks cbs, void* user_data);

#endif
