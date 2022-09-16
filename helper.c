#include "helper.h"

int poll_fd(int fd,int timeout){
    struct pollfd data[1];
    data[0].fd=fd;
    data[0].events=POLLIN|POLLPRI;
    data[0].revents=0;
    return poll(data,1,timeout);
}