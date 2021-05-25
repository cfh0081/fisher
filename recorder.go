package fisher

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cfh0081/runutils"
)

var toRecordFile *os.File = nil

type Storage struct {
	file *os.File
	ch   chan []byte
}

var storage *Storage
var once sync.Once

func (storage *Storage) WriteData(data []byte) {
	storage.ch <- data
}

func (storage *Storage) StoreData(ctx context.Context) {
	for {
		select {
		case v := <-storage.ch:
			if storage.file != nil {
				_, err := storage.file.Write(v)
				if err != nil {
					fmt.Printf("%v call file.Write with error %v.\n", runutils.RunFuncName(), err)
				}
				_, err = storage.file.WriteString("\n")
				if err != nil {
					fmt.Printf("%v call file.WriteString with error %v.\n", runutils.RunFuncName(), err)
				}
				err = storage.file.Sync()
				if err != nil {
					fmt.Printf("%v call file.Sync with error %v.\n", runutils.RunFuncName(), err)
				}
			} else {
				fmt.Println(v)
			}
		case <-ctx.Done():
			close(storage.ch)
			return
		}
	}
}

func GetInstance() *Storage {
	once.Do(func() {
		storage = &Storage{file: toRecordFile, ch: make(chan []byte, 64)}
	})
	return storage
}

func RecordData(data []byte) {
	GetInstance().WriteData(data)
}

func StartStorageServer(ctx context.Context) {
	GetInstance().StoreData(ctx)
}
