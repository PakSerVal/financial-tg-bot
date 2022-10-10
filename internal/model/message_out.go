package model

type MessageOut struct {
	Text     string
	KeyBoard *KeyBoard
}

type KeyBoard struct {
	OneTime bool
	Rows    []KeyBoardRow
}

type KeyBoardRow struct {
	Buttons []KeyBoardButton
}

type KeyBoardButton struct {
	Text string
}
