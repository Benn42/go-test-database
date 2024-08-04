package gotestdatabase

type UnicodeStyler struct {
	left  string
	right string
	line  string
	fork  string
	last  string
}

func NewUnicodeStyler() *UnicodeStyler {
	return &UnicodeStyler{
		left:  "─<",
		right: "─>",
		line:  "│",
		fork:  "├",
		last:  "└",
	}
}

func (s *UnicodeStyler) getLeft() string {
	return s.left
}

func (s *UnicodeStyler) getRight() string {
	return s.right
}

func (s *UnicodeStyler) getLine() string {
	return s.line
}

func (s *UnicodeStyler) getFork() string {
	return s.fork
}

func (s *UnicodeStyler) getLast() string {
	return s.last
}
