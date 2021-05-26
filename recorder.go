package fisher

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cfh0081/runutils"
)

type Storage struct {
	file *os.File
	ch   chan []byte
}

var storage *Storage
var once sync.Once

func (storage *Storage) SetTarget(file *os.File) {
	storage.file = file
}

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
			fmt.Printf("%v to the end.\n", runutils.RunFuncName())
			close(storage.ch)
			return
		}
	}
}

func GetStorageInstance() *Storage {
	once.Do(func() {
		storage = &Storage{file: nil, ch: make(chan []byte, 64)}
	})
	return storage
}

func RecordData(data []byte) {
	GetStorageInstance().WriteData(data)
}

// 如果要将数据存储到文件中，则先需要调用该接口，设置用于写入的文件句柄
func SetStorageTarget(file *os.File) {
	GetStorageInstance().SetTarget(file)
}

// 需要使用go关键字，单独起一个协程运行，只有开启了该服务，才能存储数据
func StartStorageServer(ctx context.Context) {
	GetStorageInstance().StoreData(ctx)
}
