package main

type Side uint

const (
	White Side = iota
	Black
)

type Clock struct {
	LeftTime    uint
	RightTime   uint
	SetTime     bool
	SideRunning Side
}

func New() Clock {
	clock := Clock{
		LeftTime: 1800,
		RightTime: 1800,
		SetTime: false,
		SideRunning: White
	}

	/*
	for dgtpicom_init(self.lib) < 0 {
		warning(logging, "Init() failed - Jack half connected?")
		show(DisplayMsg, Message.DGT_JACK_CONNECTED_ERROR())
		sleep(time, 0.5)
	}
	if dgtpicom_configure(self.lib) < 0 {
		warning(logging, "Configure() failed - Jack connected back?")
		show(DisplayMsg, Message.DGT_JACK_CONNECTED_ERROR())
	}
	show(DisplayMsg, Message.DGT_CLOCK_VERSION(2, 2, "i2c", nil))

	incoming_clock_thread := Timer(0, self._process_incoming_clock_forever)
	start(incoming_clock_thread)
	*/

	return clock
}

func _process_incoming_clock_forever(self DgtPi) {
	but := c_byte(0)
	buttime := c_byte(0)
	clktime := create_string_buffer(6)
	var counter int = 0
	info(logging, "incoming_clock ready")
	for true {
		{
			__tmp1 := self.lib_lock
			res := dgtpicom_get_button_message(self.lib, pointer(but), pointer(buttime))
			if res > 0 {
				ack3 := but.value
				if ack3 == 1 {
					info(logging, "(i2c) clock button 0 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(0, "i2c"))
				}
				if ack3 == 2 {
					info(logging, "(i2c) clock button 1 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(1, "i2c"))
				}
				if ack3 == 4 {
					info(logging, "(i2c) clock button 2 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(2, "i2c"))
				}
				if ack3 == 8 {
					info(logging, "(i2c) clock button 3 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(3, "i2c"))
				}
				if ack3 == 16 {
					info(logging, "(i2c) clock button 4 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(4, "i2c"))
				}
				if ack3 == 32 {
					info(logging, "(i2c) clock button on/off pressed")
					dgtpicom_configure(self.lib)
					show(DisplayMsg, Message.DGT_BUTTON(32, "i2c"))
				}
				if ack3 == 17 {
					info(logging, "(i2c) clock button 0+4 pressed")
					show(DisplayMsg, Message.DGT_BUTTON(17, "i2c"))
				}
				if ack3 == 64 {
					info(logging, "(i2c) clock lever pressed > right side down")
					show(DisplayMsg, Message.DGT_BUTTON(64, "i2c"))
				}
				if ack3 == -64 {
					info(logging, "(i2c) clock lever pressed > left side down")
					show(DisplayMsg, Message.DGT_BUTTON(-64, "i2c"))
				}
			}
			if res < 0 {
				warning(logging, "GetButtonMessage returned error %i", res)
			}
			dgtpicom_get_time(self.lib, clktime)
		}
		times := list(clktime.raw)
		counter = ((counter + 1) % 10)
		if counter == 0 {
			//l_hms := times[0,3]
			//r_hms := times[3,0]
			info(logging, "(i2c) clock new time received l:%s r:%s", l_hms, r_hms)
			if self.in_settime {
				info(logging, "(i2c) clock still not finished set time, sending old time")
			} else {
				if self.side_running == ClockSide.LEFT {
					self.l_time = (((l_hms[0] * 3600) + (l_hms[1] * 60)) + l_hms[2])
				}
				if self.side_running == ClockSide.RIGHT {
					self.r_time = (((r_hms[0] * 3600) + (r_hms[1] * 60)) + r_hms[2])
				}
			}
			text := DGT_CLOCK_TIME(Message, self.l_time, self.r_time, true, "i2c")
			show(DisplayMsg, text)
		}
		sleep(time, 0.1)
	}
}

func _run_configure(self DgtPi) {
	res := dgtpicom_configure(self.lib)
	if res < 0 {
		warning(logging, "Configure() also failed %i, resetting the dgtpi clock", res)
		dgtpicom_stop(self.lib)
		dgtpicom_init(self.lib)
	}
}

func _display_on_dgt_pi(self DgtPi, text string, beep bool, left_icons byte, right_icons byte) bool {
	if len(text) > 11 {
		warning(logging, "(i2c) clock message too long [%s]", text)
	}
	debug(logging, "[%s]", text)
	text = string(bytes(text, "utf-8"))
	{
		__tmp2 := self.lib_lock
		//res := dgtpicom_set_text(self.lib, text, (beep? ({ 3; }) : ({ 0; })), left_icons.value, right_icons.value)
		if res < 0 {
			warning(logging, "SetText() returned error %i, running configure", res)
			_run_configure(self)
			//res = dgtpicom_set_text(self.lib, text, (beep? ({ 3; }) : ({ 0; })), left_icons.value, right_icons.value)
		}
	}
	if res < 0 {
		warning(logging, "finally failed %i", res)
		return false
	} else {
		return true
	}
}

func display_text_on_clock(self DgtPi, message string) bool {
	text := message.l
	if text == nil {
		text = message.m
	}
	var left_icons None
	if hasattr(message, "ld") {
		left_icons = message.ld
	} else {
		left_icons = ClockIcons.NONE
	}
	var right_icons None
	if hasattr(message, "rd") {
		right_icons = message.rd
	} else {
		right_icons = ClockIcons.NONE
	}
	return _display_on_dgt_pi(self, text, message.beep, left_icons, right_icons)
}

func display_move_on_clock(self DgtPi, message string) bool {
	var bit_board, text = get_san(self, message)
	if get_new_rev2_mode(Rev2Info) {
		var text string = (". " + text)
	}
	text := format("{:3d}{:s}", bit_board.fullmove_number, text)
	var left_icons None
	if hasattr(message, "ld") {
		left_icons = message.ld
	} else {
		left_icons = ClockIcons.DOT
	}
	var right_icons None
	if hasattr(message, "rd") {
		right_icons = message.rd
	} else {
		right_icons = ClockIcons.NONE
	}
	return _display_on_dgt_pi(self, text, message.beep, left_icons, right_icons)
}

func display_time_on_clock(self DgtPi, message string) bool {
	if self.side_running != ClockSide.NONE || message.force {
		{
			__tmp3 := self.lib_lock
			res := dgtpicom_end_text(self.lib)
			if res < 0 {
				warning(logging, "EndText() returned error %i, running configure", res)
				_run_configure(self)
				res = dgtpicom_end_text(self.lib)
			}
		}
		if res < 0 {
			warning(logging, "finally failed %i", res)
			return false
		}
	} else {
		debug(logging, "(i2c) clock isnt running - no need for endText")
	}
	return true
}

func stop_clock(self DgtPi, devs set) bool {
	debug(logging, "(%s) clock sending stop time to clock l:%s r:%s", ",".join(devs), hms_time(self.l_time), hms_time(self.r_time))
	return _resume_clock(self, ClockSide.NONE)
}

func _resume_clock(self DgtPi, side ClockSide) bool {
	if self.l_time >= (3600*10) || self.r_time >= (3600*10) {
		warning(logging, "time values not set - abort function")
		return false
	}
	var l_run int = 0
	var r_run int = 0
	if side == ClockSide.LEFT {
		l_run = 1
	}
	if side == ClockSide.RIGHT {
		r_run = 1
	}
	{
		__tmp4 := self.lib_lock
		res := dgtpicom_run(self.lib, l_run, r_run)
		if res < 0 {
			warning(logging, "Run() returned error %i, running configure", res)
			_run_configure(self)
			res = dgtpicom_run(self.lib, l_run, r_run)
		}
	}
	if res < 0 {
		warning(logging, "finally failed %i", res)
		return false
	} else {
		self.side_running = side
		return true
	}
}

func start_clock(self DgtPi, side ClockSide, devs set) bool {
	l_hms := hms_time(self.l_time)
	r_hms := hms_time(self.r_time)
	debug(logging, "(%s) clock sending start time to clock l:%s r:%s", ",".join(devs), l_hms, r_hms)
	var l_run int = 0
	var r_run int = 0
	if side == ClockSide.LEFT {
		l_run = 1
	}
	if side == ClockSide.RIGHT {
		r_run = 1
	}
	{
		__tmp5 := self.lib_lock
		res := dgtpicom_set_and_run(self.lib, l_run, l_hms[0], l_hms[1], l_hms[2], r_run, r_hms[0], r_hms[1], r_hms[2])
		if res < 0 {
			warning(logging, "SetAndRun() returned error %i, running configure", res)
			_run_configure(self)
			res = dgtpicom_set_and_run(self.lib, l_run, l_hms[0], l_hms[1], l_hms[2], r_run, r_hms[0], r_hms[1], r_hms[2])
		}
	}
	if res < 0 {
		warning(logging, "finally failed %i", res)
		return false
	} else {
		self.side_running = side
		start(Timer(0.9, self.out_settime))
		return true
	}
}

func out_settime(self DgtPi) {
	self.in_settime = false
}

func (c *Clock) Set(time_left int, time_right int, devs set) bool {
	l_hms := hms_time(time_left)
	r_hms := hms_time(time_right)
	debug(logging, "(%s) clock received last time from clock l:%s r:%s [ign]", ",".join(devs), hms_time(self.l_time), hms_time(self.r_time))
	debug(logging, "(%s) clock sending set time to clock l:%s r:%s [use]", ",".join(devs), l_hms, r_hms)
	self.in_settime = true
	self.l_time = time_left
	self.r_time = time_right
	return true
}
