package store
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)
type OpType int
const(
	Put OpType=1
	Del OpType=-1
)
type WAL struct {
    file *os.File
}
// avoid 
// // only to be used by the store controller, just encoding the string to string mapping into binary format
// type Record struct{
// 	opcode OpType
// 	key []byte
// 	val []byte
// }

func OpenWAL(path string) (*WAL, error){
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err!=nil{
		return nil, err
	}
	return &WAL{file:f},nil
}

func (w *WAL)Append(opcode OpType,key,value string) error{
	line:=fmt.Sprintf("%d\t%s\t%s\n",opcode,key,value);
	if _, err:=w.file.WriteString(line); err!=nil{
		return err
	}
	// force mf os to actaully write on disk like fsync
	return w.file.Sync()
}

func (w *WAL) Replay(s *Store) error{
	if _, err:=w.file.Seek(0,0); err!=nil{
		return err
	}
	scanner:=bufio.NewScanner(w.file)
	for scanner.Scan(){
		line:=scanner.Text()
		parts:=strings.SplitN(line,"\t",3)
		if(len(parts)!=3){
			continue
		}
		opcode:=parts[0]
		key:=parts[1]
		if opcode=="1" && len(parts)==3{
			value:=parts[2]
			s.Put(key, value)
		} else if opcode=="-1"{
			s.Delete(key)
		}
	}
	return scanner.Err()
}