package main

type DgtBoard struct {
	given_device            string
	device                  string
	reverse                 bool
	is_pi                   bool
	disable_end             bool
	field_factor            int
	serial                  ST0
	lock                    ST1
	incoming_board_thread   ST2
	lever_pos               ST3
	clock_lock              bool
	last_clock_command      List
	enable_ser_clock        ST4
	watchdog_timer          ST5
	wait_counter            int
	r_time                  int
	l_time                  int
	bconn_text              ST8
	field_timer             ST9
	field_timer_running     bool
	channel                 ST10
	in_settime              bool
	low_time                bool
}

func __init__(self DgtBoard, device string, disable_revelation_leds bool, is_pi bool, disable_end bool, field_factor T0) {
	__init__(super(DgtBoard, self))
	self.given_device = device
	self.device = device
	self.reverse = false
	self.is_pi = is_pi
	self.disable_end = disable_end
	self.field_factor = (field_factor % 10)
	self.serial = nil
	self.lock = Lock()
	self.incoming_board_thread = nil
	self.lever_pos = nil
	self.clock_lock = false
	self.last_clock_command = []None{}
	self.enable_ser_clock = nil
	self.watchdog_timer = RepeatedTimer(1, self._watchdog)
	self.wait_counter = 0
	self.r_time = (3600 * 10)
	self.l_time = (3600 * 10)
	self.bconn_text = nil
	self.field_timer = nil
	self.field_timer_running = false
	self.channel = nil
	self.in_settime = false
	self.low_time = false
}

func (b *dgtBoard) set_reverse(flag T0) {
	self.reverse = flag
}

func (b *dgtBoard) get_reverse() interface{} {
	return self.reverse
}

func (b *dgtBoard) expired_field_timer() {
	debug(logging, "board position now stable => ask for complete board")
	self.field_timer_running = false
	write_command(self, []None{DgtCmd.DGT_SEND_BRD})
}

func (b *dgtBoard) stop_field_timer() {
	debug(logging, "board position was unstable => ignore former field update")
	cancel(self.field_timer)
	join(self.field_timer)
	self.field_timer_running = false
}

func (b *dgtBoard) start_field_timer() {
	var wait float64
	if self.low_time {
		wait = 0.1 + 0.06*float64(self.field_factor)
	} else {
		wait = 0.25 + 0.03*float64(self.field_factor)
	}
	debug(logging, "board position changed => wait %.2fsecs for a stable result low_time: %s", wait, self.low_time)
	self.field_timer = Timer(wait, self.expired_field_timer)
	start(self.field_timer)
	self.field_timer_running = true
}

func (b *dgtBoard) write_command(message list) bool {
	var mes None
	if message[0].value == DgtCmd.DGT_CLOCK_MESSAGE.value {
		mes = message[3]
	} else {
		mes = message[0]
	}
	if !(mes == DgtCmd.DGT_RETURN_SERIALNR) {
		debug(logging, "(ser) board put [%s] length: %i", mes, len(message))
		if mes.value == DgtClk.DGT_CMD_CLOCK_ASCII.value {
			//debug(logging, "sending text [%s] to (ser) clock", "".join(message[4..12].iter().map(|elem| chr(elem)).collect::<Vec<_>>()));
		}
		if mes.value == DgtClk.DGT_CMD_REV2_ASCII.value {
			//debug(logging, "sending text [%s] to (rev) clock", "".join(message[4..15].iter().map(|elem| chr(elem)).collect::<Vec<_>>()));
		}
	}
	var array List = []None{}
	//char_to_xl := map[string]int{"0": 63, "1": 6, "2": 91, "3": 79, "4": 102, "5": 109, "6": 125, "7": 7, "8": 127, "9": 111, "a": 95, "b": 124, "c": 88, "d": 94, "e": 123, "f": 113, "g": 61, "h": 116, "i": 16, "j": 30, "k": 117, "l": 56, "m": 85, "n": 84, "o": 92, "p": 115, "q": 103, "r": 80, "s": 109, "t": 120, "u": 62, "v": 42, "w": 126, "x": 100, "y": 110, "z": 91, " ": 0, "-": 64, "/": 82, "|": 54, "\": 100, "?": 83, "@": 101, "=": 72, "_": 8}
	for _, item := range message {
		if isinstance(item, int) {
			array = append(array, item)
		} else {
			if isinstance(item, enum.Enum) {
				array = append(array, item.value)
			} else {
				if isinstance(item, str) {
					for _, character := range item {
						array = append(array, char_to_xl[character.lower()])
					}
				} else {
					error(logging, "type not supported [%s]", type_(item))
					return false
				}
			}
		}
	}
	for true {
		if self.serial {
			{
				__tmp1 := self.lock
				write(self.serial, bytearray(array))
				break
			}
		}
		if mes == DgtCmd.DGT_RETURN_SERIALNR {
			break
		}
		sleep(time, 0.1)
	}
	if message[0] == DgtCmd.DGT_CLOCK_MESSAGE {
		self.last_clock_command = message
		if self.clock_lock {
			warning(logging, "(ser) clock is already locked. Maybe a \"resend\"?")
		} else {
			debug(logging, "(ser) clock is locked now")
		}
		self.clock_lock = time(time)
	} else {
		sleep(time, 0.1)
	}
	return true
}

func (b *dgtBoard) processBoardMessage(message_id int, message tuple, message_length int) {
	if false {
		// pass

	} else {
		if message_id == DgtMsg.DGT_MSG_VERSION {
			if message_length != 2 {
				warning(logging, "illegal length in data")
			}
			var board_version string = ((String(message[0]) + ".") + String(message[1]))
			debug(logging, "(ser) board version %0.2f", float(board_version))
			write_command(self, []None{DgtCmd.DGT_SEND_BRD})
			var text_l, text_m, text_s = "USB e-Board", "USBboard", "ok usb"
			self.channel = "USB"
			self.bconn_text = DISPLAY_TEXT(Dgt, text_l, text_m, text_s, true, false, 1.1, map[string]bool{"i2c": true, "web": true})
			show(DisplayMsg, Message.DGT_EBOARD_VERSION(self.bconn_text, self.channel))
			startup_serial_clock(self)
			if is_running(self.watchdog_timer) {
				warning(logging, "watchdog timer is already running")
			} else {
				debug(logging, "watchdog timer is started")
				start(self.watchdog_timer)
			}
		} else {
			if message_id == DgtMsg.DGT_MSG_BWTIME {
				if message_length != 7 {
					warning(logging, "illegal length in data")
				}
				if (message[0] && 15) == 10 || (message[3] && 15) == 10 {
					var ack0 int = ((message[1] && 127) || ((message[3] << 3) && 128))
					var ack1 int = ((message[2] && 127) || ((message[3] << 2) && 128))
					var ack2 int = ((message[4] && 127) || ((message[0] << 3) && 128))
					var ack3 int = ((message[5] && 127) || ((message[0] << 2) && 128))
					if ack0 != 16 {
						warning(logging, "(ser) clock ACK error %s", ack0, ack1, ack2, ack3)
						if self.last_clock_command {
							debug(logging, "(ser) clock resending failed message [%s]", self.last_clock_command)
							write_command(self, self.last_clock_command)
							self.last_clock_command = []None{}
						}
						return
					} else {
						debug(logging, "(ser) clock ACK okay [%s]", DgtAck(ack1))
						if self.last_clock_command {
							cmd := self.last_clock_command[3]
							if cmd.value != ack1 && ack1 < 128 {
								warning(logging, "(ser) clock ACK [%s] out of sync - last: [%s]", DgtAck(ack1), cmd)
							}
						}
					}
					if ack1 == DgtAck.DGT_ACK_CLOCK_BUTTON.value {
						if ack3 == 49 {
							info(logging, "(ser) clock button 0 pressed - ack2: %i", ack2)
							show(DisplayMsg, Message.DGT_BUTTON(0, "ser"))
						}
						if ack3 == 52 {
							info(logging, "(ser) clock button 1 pressed - ack2: %i", ack2)
							show(DisplayMsg, Message.DGT_BUTTON(1, "ser"))
						}
						if ack3 == 51 {
							info(logging, "(ser) clock button 2 pressed - ack2: %i", ack2)
							show(DisplayMsg, Message.DGT_BUTTON(2, "ser"))
						}
						if ack3 == 50 {
							info(logging, "(ser) clock button 3 pressed - ack2: %i", ack2)
							show(DisplayMsg, Message.DGT_BUTTON(3, "ser"))
						}
						if ack3 == 53 {
							if ack2 == 69 {
								info(logging, "(ser) clock button 0+4 pressed - ack2: %i", ack2)
								show(DisplayMsg, Message.DGT_BUTTON(17, "ser"))
							} else {
								info(logging, "(ser) clock button 4 pressed - ack2: %i", ack2)
								show(DisplayMsg, Message.DGT_BUTTON(4, "ser"))
							}
						}
					}
					if ack1 == DgtAck.DGT_ACK_CLOCK_VERSION.value {
						self.enable_ser_clock = true
						var main int = (ack2 >> 4)
						var sub int = (ack2 && 15)
						debug(logging, "(ser) clock version %0.2f", float(((String(main) + ".") + String(sub))))
						var dev string
						if self.bconn_text {
							self.bconn_text.devs = map[string]bool{"ser": true}
							dev = "ser"
						} else {
							dev = "err"
						}
						show(DisplayMsg, Message.DGT_CLOCK_VERSION(main, sub, dev, self.bconn_text))
					}
				} else {
					//if(any(message[..7])) {
					if 1 {
						var r_hours int = (message[0] && 15)
						var r_mins int = (((message[1] >> 4) * 10) + (message[1] && 15))
						var r_secs int = (((message[2] >> 4) * 10) + (message[2] && 15))
						var l_hours int = (message[3] && 15)
						var l_mins int = (((message[4] >> 4) * 10) + (message[4] && 15))
						var l_secs int = (((message[5] >> 4) * 10) + (message[5] && 15))
						var r_time int = (((r_hours * 3600) + (r_mins * 60)) + r_secs)
						var l_time int = (((l_hours * 3600) + (l_mins * 60)) + l_secs)
						errtim := r_hours > 9 || l_hours > 9 || r_mins > 59 || l_mins > 59 || r_secs > 59 || l_secs > 59
						if errtim {
							warning(logging, "(ser) clock illegal new time received %s", message)
						} else {
							if r_time > self.r_time || l_time > self.l_time {
								warning(logging, "(ser) clock strange old time received %s l:%s r:%s", message, hms_time(self.l_time), hms_time(self.r_time))
								if self.in_settime {
									info(logging, "(ser) clock still in set mode, ignore received time")
									errtim = true
								} else {
									if (r_time-self.r_time) > 3600 || (l_time-self.l_time) > 3600 {
										info(logging, "(ser) clock new time over 1h difference, ignore received time")
										errtim = true
									}
								}
							} else {
								info(logging, "(ser) clock new time received l:%s r:%s", hms_time(l_time), hms_time(r_time))
								var status int = (message[6] && 63)
								var connect int = !(status && 32)
								if connect {
									var right_side_down int
									if status && 2 {
										right_side_down = -64
									} else {
										right_side_down = 64
									}
									if self.lever_pos != right_side_down {
										debug(logging, "(ser) clock button status: 0x%x old lever: %s", status, self.lever_pos)
										if self.lever_pos != nil {
											show(DisplayMsg, Message.DGT_BUTTON(right_side_down, "ser"))
										}
										self.lever_pos = right_side_down
									}
								} else {
									info(logging, "(ser) clock not connected, sending old time l:%s r:%s", hms_time(self.l_time), hms_time(self.r_time))
									l_time = self.l_time
									r_time = self.r_time
								}
								if self.in_settime {
									info(logging, "(ser) clock still in set mode, sending old time l:%s r:%s", hms_time(self.l_time), hms_time(self.r_time))
									l_time = self.l_time
									r_time = self.r_time
								}
								show(DisplayMsg, Message.DGT_CLOCK_TIME(l_time, r_time, connect, "ser"))
								if !(self.enable_ser_clock) {
									var dev string
									dev = "ser"
									if is_running(self.watchdog_timer) {
										info(logging, "(%s) clock restarting setup", dev)
										startup_serial_clock(self)
									} else {
										info(logging, "(%s) clock sends messages already but (%s) board still not found", dev, dev)
									}
								}
							}
						}
						if !(errtim) {
							self.r_time = r_time
							self.l_time = l_time
						}
					} else {
						debug(logging, "(ser) clock null message ignored")
					}
				}
				if self.clock_lock {
					debug(logging, "(ser) clock unlocked after %.3f secs", (time.time() - self.clock_lock))
					self.clock_lock = false
				}
			} else {
				if message_id == DgtMsg.DGT_MSG_BOARD_DUMP {
					if message_length != 64 {
						warning(logging, "illegal length in data")
					}
					piece_to_char := map[int]string{1: "P", 2: "R", 3: "N", 4: "B", 5: "K", 6: "Q", 7: "p", 8: "r", 9: "n", 10: "b", 11: "k", 12: "q", 13: "$", 14: "%", 15: "&", 0: "."}
					var board string = ""
					for _, character := range message {
						board += piece_to_char[(character && 15)]
					}
					//debug(logging, ("
					//" + "
					//".join(iter.NewIntSeq(iter.Start(0), iter.Stop(len(board)), iter.Step(8)).All().iter().map(|i| board[(0 + i)..(8 + i)]).collect::<Vec<_>>())));
					var fen string = ""
					var empty int = 0
					for _, square := range iter.NewIntSeq(iter.Start(0), iter.Stop(64)).All() {
						if message[square] != 0 && message[square] < 13 {
							if empty > 0 {
								fen += String(empty)
								var empty int = 0
							}
							fen += piece_to_char[(message[square] && 15)]
						} else {
							empty += 1
						}
						if ((square + 1) % 8) == 0 {
							if empty > 0 {
								fen += String(empty)
								empty = 0
							}
							if square < 63 {
								fen += "/"
							}
						}
					}
					debug(logging, "raw fen [%s]", fen)
					show(DisplayMsg, Message.DGT_FEN(fen, true))
				} else {
					if message_id == DgtMsg.DGT_MSG_FIELD_UPDATE {
						if message_length != 2 {
							warning(logging, "illegal length in data")
						}
						if self.field_timer_running {
							stop_field_timer(self)
						}
						start_field_timer(self)
					} else {
						if message_id == DgtMsg.DGT_MSG_SERIALNR {
							if message_length != 5 {
								warning(logging, "illegal length in data")
							}
							//show(DisplayMsg, Message.DGT_SERIAL_NR("".join(message.iter().map(|elem| chr(elem)).collect::<Vec<_>>())));
						} else {
							if message_id == DgtMsg.DGT_MSG_LONG_SERIALNR {
								if message_length != 10 {
									warning(logging, "illegal length in data")
								}
								//number := join("", message.iter().map(|elem| chr(elem)).collect::<Vec<_>>())
							}
						}
					}
				}
			}
		}
	}
}

func (b *dgtBoard) _read_serial(bytes_toread uint) interface{} {
	return read(self.serial, bytes_toread)
}

func (b *dgtBoard) _read_board_message(head []uint8) interface{} {
	//message :=
	var header_len int = 3
	var header []uint8 = (head + _read_serial(self, (header_len-1)))
	header = unpack(struct_, ">BBB", header)
	message_id := header[0]
	var message_length int = (((header[1] << 7) + header[2]) - header_len)
	var counter int = (((header[1] << 7) + header[2]) - header_len)
	if message_length <= 0 || message_length > 64 {
		if message_id == 143 && message_length == 7936 {
			warning(logging, "falsely DGT_SEND_EE_MOVES send before => receive and ignore EE_MOVES result")
			stop(self.watchdog_timer)
			now := time(time)
			for counter > 0 {
				ee_moves := _read_serial(self, counter)
				info(logging, "EE_MOVES 0x%x bytes read - inWaiting: 0x%x", len(ee_moves), self.serial.inWaiting())
				counter -= len(ee_moves)
				if (time(time) - now) > 15 {
					warning(logging, "EE_MOVES needed over 15secs => ignore not readed 0x%x bytes now", counter)
					break
				}
			}
			start(self.watchdog_timer)
		} else {
			warning(logging, "illegal length in message header 0x%x length: %i", message_id, message_length)
		}
		return message
	}
	for counter {
		byte := _read_serial(self)
		if byte {
			data := unpack(struct_, ">B", byte)
			counter -= 1
			if data[0] && 128 {
				warning(logging, "illegal data in message 0x%x found", message_id)
				warning(logging, "ignore collected message data %s", message)
				return _read_board_message(self, byte)
			}
			message += data
		} else {
			warning(logging, "timeout in data reading")
		}
	}
	processBoardMessage(self, message_id, message, message_length)
	return message
}

func (b *dgtBoard) _process_incoming_board_forever() {
	var counter int = 0
	info(logging, "incoming_board ready")
	for true {
		//var byte []uint8 = b""
		if self.serial {
			byte = _read_serial(self)
		} else {
			_setup_serial_port(self)
			if self.serial {
				debug(logging, "sleeping for 0.5 secs. Afterwards startup the (ser) board")
				sleep(time, 0.5)
				counter = 0
				_startup_serial_board(self)
			}
		}
		if byte && (byte[0] && 128) {
			_read_board_message(self, byte)
		} else {
			counter = ((counter + 1) % 10)
			if counter == 0 && !(is_running(self.watchdog_timer)) {
				_watchdog(self)
			}
			sleep(time, 0.1)
		}
	}
}

func (b *dgtBoard) ask_battery_status() {
	write_command(self, []None{DgtCmd.DGT_SEND_BATTERY_STATUS})
}

func (b *dgtBoard) startup_serial_clock() {
	self.clock_lock = false
	self.enable_ser_clock = false
	command := []None{DgtCmd.DGT_CLOCK_MESSAGE, 3, DgtClk.DGT_CMD_CLOCK_START_MESSAGE, DgtClk.DGT_CMD_CLOCK_VERSION, DgtClk.DGT_CMD_CLOCK_END_MESSAGE}
	write_command(self, command)
}

func (b *dgtBoard) _startup_serial_board() {
	write_command(self, []None{DgtCmd.DGT_SEND_UPDATE_NICE})
	write_command(self, []None{DgtCmd.DGT_SEND_VERSION})
}

func (b *dgtBoard) _watchdog() {
	if self.clock_lock && !(self.is_pi) {
		if (time(time) - self.clock_lock) > 2 {
			warning(logging, "(ser) clock is locked over 2secs")
			debug(logging, "resending locked (ser) clock message [%s]", self.last_clock_command)
			self.clock_lock = false
			write_command(self, self.last_clock_command)
		}
	}
	write_command(self, []None{DgtCmd.DGT_RETURN_SERIALNR})
}

func (b *dgtBoard) _open_serial(device string) bool {
	if !(!(self.serial)) {
		panic("assert")
	}
	self.serial = Serial(device, STOPBITS_ONE, PARITY_NONE, EIGHTBITS, 0.5)
	return true
}

func (b *dgtBoard) _success(device string) bool {
	self.device = device
	debug(logging, "(ser) board connected to %s", self.device)
	return true
}

func (b *dgtBoard) _setup_serial_port() bool {

	//var waitchars []string = []string{"/", "-", "\", "|"}
	if is_running(self.watchdog_timer) {
		debug(logging, "watchdog timer is stopped now")
		stop(self.watchdog_timer)
	}
	if self.serial {
		return true
	}
	__tmp2 := self.lock
	if self.given_device {
		if _open_serial(self, self.given_device) {
			return _success(self.given_device)
		}
	} else {
		for _, file := range listdir("/dev") {
			if startswith(file, "ttyACM") || startswith(file, "ttyUSB") || file == "rfcomm0" {
				dev := join(path, "/dev", file)
				if _open_serial(self, dev) {
					return _success(dev)
				}
			}
		}
		if _open_bluetooth(self) {
			return _success("/dev/rfcomm123")
		}
	}
	var bwait string = ("Board" + waitchars[self.wait_counter])
	//text := DISPLAY_TEXT(Dgt, ("no e-" + bwait), ("no" + bwait), bwait, true, false, 0.1, map[string]bool{"i2c": true, "web": true})
	//show(DisplayMsg, Message.DGT_NO_EBOARD_ERROR(text));
	//self.wait_counter = ((self.wait_counter + 1) % len(waitchars))
	return false
}

func (b *dgtBoard) _wait_for_clock(fn string) {
	var has_to_wait bool = false
	var counter int = 0
	for self.clock_lock {
		if !(has_to_wait) {
			has_to_wait = true
			debug(logging, "(ser) clock is locked => waiting to serve: %s", func_)
		}
		sleep(time, 0.1)
		counter = ((counter + 1) % 30)
		if counter == 0 {
			warning(logging, "(ser) clock is locked over 3secs")
		}
	}
	if has_to_wait {
		debug(logging, "(ser) clock is released now")
	}
}

func (b *dgtBoard) set_text_3k(text string, beep int) bool {
	_wait_for_clock(self, "SetText3K()")
	set_pi_mode(Rev2Info, true)
	var res bool = write_command(self, []None{DgtCmd.DGT_CLOCK_MESSAGE, 12, DgtClk.DGT_CMD_CLOCK_START_MESSAGE, DgtClk.DGT_CMD_CLOCK_ASCII, text[0], text[1], text[2], text[3], text[4], text[5], text[6], text[7], beep, DgtClk.DGT_CMD_CLOCK_END_MESSAGE})
	return res
}

func set_and_run(self DgtBoard, lr int, lh int, lm int, ls int, rr int, rh int, rm int, rs int) bool {
	_wait_for_clock(self, "SetAndRun()")
	side := ClockSide.NONE
	if lr == 1 && rr == 0 {
		side = ClockSide.LEFT
	}
	if lr == 0 && rr == 1 {
		side = ClockSide.RIGHT
	}
	var res bool = write_command(self, []None{DgtCmd.DGT_CLOCK_MESSAGE, 10, DgtClk.DGT_CMD_CLOCK_START_MESSAGE, DgtClk.DGT_CMD_CLOCK_SETNRUN, lh, lm, ls, rh, rm, rs, side, DgtClk.DGT_CMD_CLOCK_END_MESSAGE})
	return res
}

func end_text(self DgtBoard) bool {
	_wait_for_clock(self, "EndText()")
	var res bool = write_command(self, []None{DgtCmd.DGT_CLOCK_MESSAGE, 3, DgtClk.DGT_CMD_CLOCK_START_MESSAGE, DgtClk.DGT_CMD_CLOCK_END, DgtClk.DGT_CMD_CLOCK_END_MESSAGE})
	return res
}

func run(self DgtBoard) {
	self.incoming_board_thread = Timer(0, self._process_incoming_board_forever)
	start(self.incoming_board_thread)
}
