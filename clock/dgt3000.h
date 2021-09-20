/* functions to communicate to a DGT3000 using I2C
 * version 0.7
 *
 * Copyright (C) 2015 DGT
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */


/* find out wich pi
	returns:
	0 = error
	1 = Pi b+
	2 = Pi 2 */
int checkPiModel();

/* calculate checksum and put it in the last byte
	*buffer = pointer to buffer */
char crc_calc(char *buffer);

/* print hex values
	bytes = array of bytes
	length is number of bytes to print */
void hexPrint(char bytes[], int length);


/* configure IO pins and I2C Master and Slave
	*/
void i2cReset();

/* get message from I2C receive buffer
	m[] = message buffer of 256 bytes
	timeOut = time to wait for packet in us (0=dont wait)
	returns:
	-6 = CRC Error
	-5 = I2C buffer overrun, at least 16 bytes received succesfully. rest is missing.
	-4 = our buffer overrun (should not happen)
	-3 = timeout
	-2 = I2C Error
	>0 = packet length*/
int i2cReceive(char m[]);

/* send message using I2CMaster
	 message[] = the message to send
	 returns:
	 -7 = message not Acked, probably clock off
	 -3 = message not Acked, probably collision
	 -2 = I2C Error
	 0 = Succes */
int i2cSend(char message[], char ackAdr);


/* Get direct access to BCM2708/9
	returns:
	-2 = I2C connection error no dgtpi?
	-1 = fail, no root access to /dev/mem
	0 = succes
	1 = fail (run as root!) */
int dgt3000Init(void);

/* configure dgt3000 to on, central controll and mode 25
	returns:
	-3 = sending commands failed
	-2 = I2C Error
	0 = succes */
int dgt3000Configure();

/* send a wake command to the dgt3000
	returns:
	-3 = wake ack error
	-1 = no hello message received
	0 = succes */
int dgt3000Wake();

/* send set central controll command to dgt3000
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000SetCC();

/* send set mode 25 to dgt3000
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000Mode25();

/* send end display to dgt3000 to clear te display
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000EndDisplay();

/* send set display command to dgt3000
	dm = display message to send
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000SetDisplay(char dm[]);

/* send set and run command to dgt3000
	lr/rr = left/right run mode, 0=stop, 1=count down, 2=count up
	lh/rh = left/right hours
	lm/rm = left/right minutes
	ls/rs = left/right seconds
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000SetNRun(char lr, char lh, char lm, char ls,
					char rr, char rh, char rm, char rs);

/* send set and run command to dgt3000 with current clock values
	lr/rr = left/right run mode, 0=stop, 1=count down, 2=count up
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000Run(char lr, char rr);

/* try three times to end and set de display
	text = message to display
	beep = beep length (/62.5ms) max 48 (3s)
	ld/rd = left/right buttons:
		1=flag,
		2=white king,
		4=black king,
		8=colon,
		16=dot,
		32=extra dot (left only)
	returns:
	-3 = sending failed
	0 = succes */
int dgt3000Display(char text[], char beep, char ld, char rd);

/* check for messages from dgt3000
	returns:
	0 = nothing is received
	1 = something is received
	2 = off button message is received */
void *dgt3000Receive(void *);

/* wait for an Ack message
	adr = adress to listen for ack
	cmd = command to ack
	timeOut = time to wait for ack
	returns:
	-3 = no Ack
	0 = Ack */
int dgt3000GetAck(char adr, char cmd, long long int timeOut);

/* return last received time
	time[] = 6 byte time descriptor */
void dgt3000GetTime(char time[]);

/* return buttons pressed
	buttons = buttons pressed
	  binary:
		0x01 < back
		0x02 - minus
		0x04   play/pause
		0x08 + plus
		0x10 > forward
	  special:
		0x20   off
		0xa0   on
		0x40   Lever changed, right side down
		0xc0   Lever changed, left side down */
int dgt3000GetButton(char *buttons, char *time);

/* return current button state
	returns:
	  binary:
		0x01 < back
		0x02 - minus
		0x04   play/pause
		0x08 + plus
		0x10 > forward
		0x20   on/off button
		0x40   Lever changed, (right side down = 1) */
int dgt3000GetButtonState();

/* turn off dgt3000
	returnMode = Mode clock will start in when turned on
	returns:
	-3 = sending failed, clock off (or collision)
	-2 = sending failed, I2C error
	-1 = no (positive)ack received, not in CC
	0 = succes */
int dgt3000Off(char returnMode);

/* stop receiving
	*/
int dgt3000Stop();
