package main

import (
  "fmt"
  "os"

  //"./board"
  "./clock"
  "./client"
)

func main() {
	result := clock.Init()
  if result < 0 {
    os.Exit(result)
  }
	fmt.Println(result)
}

/*
res = lib.dgt3000Init()
if(res < 0):
	sys.exit(res)
res = lib.dgt3000Configure()
if(res < 0):
	sys.exit(res)

x = lib.dgt3000SetNRun(1, 4, 23, 56, 0, 3, 12, 45)
print(x)

while(run):
	if(lib.dgt3000GetButton(ctypes.pointer(but),ctypes.pointer(buttime))==1):
		if (but.value==1):
			run=0
		else:
			print(but,buttime)

		if (but.value==2):
                        lib.dgt3000Display(b"Hello World",0,0,0)
                        print(lib.dgt3000GetButtonState())

                if (but.value==4):
                        lib.dgt3000EndDisplay()

        lib.dgt3000GetTime(clktime)
        times = list(clktime.raw)
        print(times)

	time.sleep(0.1)
print(but,buttime)
print(lib.dgt3000Off(5))
lib.dgt3000Stop()

*/
