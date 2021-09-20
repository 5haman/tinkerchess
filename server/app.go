package main

import (
	"io"
	"log"
	"os"
	"strings"

	"./board"

	"github.com/jacobsa/go-serial/serial"
	"github.com/notnil/chess"
	"github.com/jessevdk/go-flags"
)

type DgtAppArgs struct {
	Device string `short:"d" long:"device" description:"Board device (e.g. /dev/ttyUSB0, /dev/ttyACM0)" default:"/dev/ttyACM0"`

	Pargs []string
}

func GetParsedArguments() *DgtAppArgs {
	var args DgtAppArgs
	parser := flags.NewParser(&args, flags.Default)
	pargs, err := parser.ParseArgs(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	args.Pargs = pargs

	return &args
}

// DgtApp is the main class-like struct for the application.
type DgtApp struct {
	log      *log.Logger
	args     *DgtAppArgs
	startFEN string
	port     io.ReadWriteCloser
	board    *board.DgtBoard
	currGame *chess.Game
	prevGame *chess.Game
}

func NewDgtApp(args *DgtAppArgs) *DgtApp {
	return &DgtApp{
		args: args,
		log: log.New(os.Stdout, "[DGT] ",
			log.Ldate|log.Ltime),
	}
}

// Run is the main entry point for DgtApp.
func (a *DgtApp) Run() {
	a.getStartFEN()
	a.openPort()
	a.createBoard()
	a.initialiseBoard()
	a.runForever()
}

func (a *DgtApp) getStartFEN() {
	startFEN := chess.NewGame().Position().String()
	a.startFEN = strings.Replace(startFEN, " w KQkq - 0 1", "", -1)

}

func (a *DgtApp) openPort() {
	a.log.Println("Opening port ...")
	options := serial.OpenOptions{
		PortName:        a.args.Device,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	a.check(err)
	a.port = port
}

func (a *DgtApp) createBoard() {
	a.log.Println("Creating board ...")
	a.board = board.NewDgtBoard(a.port)
}

func (a *DgtApp) initialiseBoard() {
	a.log.Println("Starting byte reader ...")
	a.board.StartByteReader()
	a.log.Println("Starting command processor ...")
	a.board.StartCommandProcessor()
	a.log.Println("Resetting board ...")
	_, err := a.board.WriteSendResetCommand()
	a.check(err)
	_, err = a.board.WriteSendBoardCommand()
	a.check(err)
	_, err = a.board.WriteSendUpdateBoardCommand()
	a.check(err)
}

func (a *DgtApp) runForever() {
	for {
		select {
		case bm := <-a.board.GetBoardMessageChannel():
			a.processBoardMessage(bm)
		}
	}
}

func (a *DgtApp) processBoardMessage(bm *board.BoardMessage) {
	if bm.FieldUpdate != nil {
		a.handleFieldUpdate(bm.FieldUpdate)
	} else {
		a.log.Println(bm.UnhandledMessage)
	}
}

func (a *DgtApp) handleFieldUpdate(fieldUpdate *board.Update) {
	// Every time we receive a field update, ask for a complete
	// board update. This is simple but inefficient, and will do
	// for the time being.
	//
	// Note that it's tempting to limit such requests only to
	// "piece drop" messages, but that's actually incorrect, for
	// at least two reasons. Firstly, we can receive the piece
	// drop and piece lift messages out of order. Secondly, a
	// capture executed as "use capturing piece to push captured
	// piece off the square, then remove captured piece" might end
	// with a final piece lift message.
	// log.Printf("Field update %s to %s\n", fieldUpdate.piece,
	//	fieldUpdate.square)
	_, err := a.board.WriteSendBoardCommand()
	a.check(err)
}

func (a *DgtApp) check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
