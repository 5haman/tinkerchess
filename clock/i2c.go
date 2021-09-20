package main

type DGTi2c struct {
	device ST0
	serial ST1
}

func __init__(self DGTi2c, device string) *DGTi2c {
	self.device = device
}

func write_text_to_clock(self DGTi2c, message string, beep bool) {
	beep_code := 0
	if beep {
		beep_code = 3
	}
	dgt3000Display(self.lib, message, beep_code, 0, 0)
}

func write_stop_to_clock(self DGTi2c) {
	dgt3000SetNRun(self.lib, 0, 0, 0, 0, 0, 0, 0, 0)
}

func write_start_to_clock(self DGTi2c, l_hms byte, r_hms byte, side byte) {
	var lr int
	var rr int
	if side == 1 {
		lr = 1
		rr = 0
	} else {
		lr = 0
		rr = 1
	}
	dgt3000SetNRun(self.lib, lr, l_hms[0], l_hms[1], l_hms[2], rr, r_hms[0], r_hms[1], r_hms[2])
}

func startup_clock(self DGTi2c) {
	dgt3000Init(self.lib)
	dgt3000Configure(self.lib)
}

func process_incoming_clock_forever(self DGTi2c) {
	but := ctypes.c_byte(0)
	buttime := ctypes.c_byte(0)
	for true {
		if dgt3000GetButton(self.lib, ctypes.pointer(but), ctypes.pointer(buttime)) == 1 {
			ack3 := but.value
			if ack3 == 1 {
				info(logging, "Button 0 pressed")
				show(Display, Message.DGT_BUTTON, 0)
			}
			if ack3 == 2 {
				info(logging, "Button 1 pressed")
				show(Display, Message.DGT_BUTTON, 1)
			}
			if ack3 == 4 {
				info(logging, "Button 2 pressed")
				show(Display, Message.DGT_BUTTON, 2)
			}
			if ack3 == 8 {
				info(logging, "Button 3 pressed")
				show(Display, Message.DGT_BUTTON, 3)
			}
			if ack3 == 16 {
				info(logging, "Button 4 pressed")
				show(Display, Message.DGT_BUTTON, 4)
			}
			if ack3 == 32 {
				info(logging, "Button on/off pressed")
			}
			if ack3 == 64 {
				info(logging, "Lever pressed")
			}
		}
		sleep(time, 0.1)
	}
}

func run(self DGTi2c) {
	startup_clock(self)
	incoming_clock_thread := Timer(threading, 0, self.process_incoming_clock_forever)
	start(incoming_clock_thread)
}
