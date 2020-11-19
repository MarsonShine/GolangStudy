package thumbnails_test

import (
	"log"
	"sync"
	"tutorials/goroutineAndchannel/thumbnails"
)

func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnails.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

func makeThumbnails2(filenames []string) {
	for _,f ：= range filenafilenames{
		go thumbnails.ImageFile(f)	// 开启协程执行，没有等待完成程序就结束了
	}
}

func makeThumbnails3(filenames []string) {
	ch:=make(chan struct{})
	for _,f:= range filenames {
		go func(f string) {	// 注意，这里是将 f 做每个协程的局部变量
			thumbnails.ImageFile(f)
			ch <- struct{}{}
		}(f)
		// ！！！！这种闭包是很危险的，是不正确的
		// go func() {
		// 	thumbnails.ImageFile(f)	// 这里引用了外部变量
		// }()
	}
	// 等待所有协程任务完成
	for range filenames {
		<-ch
	}
}

func makeThumbnails4(filenames []string) error{
	errors:=make(chan error)

	for _,f:= range filenames {
		go func(f string) {
			_,err := thumbnails.ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
        if err := <-errors; err != nil {
            return err // NOTE: incorrect: goroutine leak!	协程泄露，因为当遇到非nil时，会直接返回给调用方法，整个方法结束，剩下的 goroutine 没有处理，会永远阻塞下去
        }
	}
	
	return nil
}

func makeThumbnails5(filenames []string)(thumbfiles []string,err error) {
	type item struct {
        thumbfile string
        err       error
    }

	ch := make(chan item, len(filenames))
	
	for _,f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnails.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil,it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}
	return thumbfiles, nil
}

func makeThumbnails6(filenames <- chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup	// 工作协程计数
	for f:= range filenames {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			thumb,err:=thumbnails.ImageFile(f)
			if err!=nil {
				log.Println(err)
                return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
            sizes <- info.Size()
		}(f)
	}
	// 关闭
	go func(){
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size:= range sizes {
		total +=size
	}

	return total
}
