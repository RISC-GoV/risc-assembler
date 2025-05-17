package assembler

type Program struct {
	machinecode []byte
	variables   []byte
	constants   []byte
	strings     []byte
	entrypoint  [4]byte
}
