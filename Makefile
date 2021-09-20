

CC=gcc

CFLAGS=-pthread -Wall

all:
	$(CC) $(CFLAGS) -shared -fPIC -o dgtpicom.so dgt3000/dgtpicom.c
