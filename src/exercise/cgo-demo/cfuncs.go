// 传递给C语言的Go回调函数的包装器。
package main

/*
extern void goStart(void*, int);
extern void goEnd(void*, int, int);
void startCgo(void* user_data, int i) {
  goStart(user_data, i);
}
void endCgo(void* user_data, int a, int b) {
  goEnd(user_data, a, b);
}
*/
import "C"
