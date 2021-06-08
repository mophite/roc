package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"roc/internal/x/bytesbuffpool"
	"roc/rlog/common"
)

const fileNameFormat = "2006.01.02.15:04:05.000"

type FileSysOut struct {
	sync.Mutex
	opts       Option
	file       *os.File
	fileWriter *bufio.Writer
	fileSize   int
	async      bool
	Bucket     chan *bytes.Buffer
	close      chan struct{}
	level      common.Level
	timestamp  int
	ticker     *time.Ticker
}

func NewFs(opts ...Options) *FileSysOut {
	f := &FileSysOut{}
	f.opts = newOpts(opts...)
	return f
}

func (l *FileSysOut) Init(string) {
	return
}

func (l *FileSysOut) Out(level common.Level, b *bytes.Buffer) {
	if level < l.level {
		return
	}

	l.Bucket <- b
}

func (l *FileSysOut) Level() common.Level {
	return l.level
}

func (l *FileSysOut) SetLevel(level common.Level) {
	l.level = level
}

func (l *FileSysOut) String() string {
	return "file"
}

func (l *FileSysOut) Link() string {
	return l.opts.link
}

func (l *FileSysOut) Name() string {
	return l.opts.name
}

func (l *FileSysOut) Path() string {
	return l.opts.path
}

func (l *FileSysOut) Async() bool {
	return l.opts.async
}

func (l *FileSysOut) Prefix() string {
	return l.opts.prefix
}

func (l *FileSysOut) MaxBucketSize() int {
	return l.opts.maxBucketSize
}

func (l *FileSysOut) Rotate() bool {
	return l.opts.rotate
}

func (l *FileSysOut) Interval() time.Duration {
	return l.opts.interval
}

func (l *FileSysOut) Filename() string {
	return l.opts.fileName
}

func (l *FileSysOut) Zone() *time.Location {
	return l.opts.zone
}

func (l *FileSysOut) FileModel() os.FileMode {
	return os.FileMode(l.opts.modePerm)
}

func (l *FileSysOut) MaxBufferSize() int {
	return l.opts.maxBufferSize
}

func (l *FileSysOut) loadLink() (err error) {
	l.Lock()
	defer l.Unlock()

	l.opts.fileName, err = isLink(l.Link())
	if err != nil {
		return err
	}

	l.file, err = open(l.Filename(), l.FileModel())
	if err != nil {
		return err
	}

	info, err := os.Stat(l.Filename())
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(
		fileNameFormat,
		getFilenamePrefix(info.Name()),
		l.Zone(),
	)

	if err != nil {
		return err
	}

	l.timestamp = convertTime(t)
	l.fileSize = int(info.Size())
	l.fileWriter = bufio.NewWriterSize(l.file, l.MaxBufferSize())

	return nil
}

func (l *FileSysOut) create() error {
	l.Lock()
	defer l.Unlock()

	var err error
	if !pathIsExist(l.Path()) {
		if err = os.MkdirAll(l.Path(), os.ModePerm); err != nil {
			return err
		}
	}

	var now = time.Now()

	l.timestamp = convertTime(now)

	l.opts.fileName = filepath.Join(
		l.Path(),
		now.Format(fileNameFormat)+".log")

	f, err := open(l.Filename(), l.FileModel())
	if err != nil {
		return err
	}

	l.file = f
	l.fileWriter = bufio.NewWriterSize(f, l.MaxBucketSize())

	_ = os.Remove(l.Link())
	return os.Symlink(l.Filename(), l.Link())
}

func (l *FileSysOut) Poller() {
	if l.loadLink() != nil {
		err := l.create()
		if err != nil {
			panic(err)
		}
	}

QUIT:
	for {
		select {
		case <-l.ticker.C:
			if l.fileWriter.Size() > 0 {
				l.fflush()
			}

		case n := <-l.Bucket:
			l.rotateWrite(n)

		case <-l.close:
			break QUIT
		}
	}
}

func convertTime(t time.Time) int {
	y, m, d := t.Date()
	return y*10000 + int(m)*100 + d*1
}

func (l *FileSysOut) fflush() {
	l.Lock()
	_ = l.fileWriter.Flush()
	l.Unlock()
}

func (l *FileSysOut) rotateWrite(b *bytes.Buffer) {
	n, _ := l.fileWriter.Write(b.Bytes())
	l.fileSize += n

	bytesbuffpool.Put(b)

	if n <= 0 {
		return
	}

	timestamp := convertTime(time.Now())

	if l.fileSize <= l.opts.maxFileSize && timestamp <= l.timestamp {
		return
	}

	l.Lock()
	defer l.Unlock()

	l.fflush()
	l.closeFile()
	err := l.create()
	if err != nil {
		fmt.Println("rotateWrite and create file err=", err)
		return
	}
	l.fileWriter.Reset(l.file)
}

func (l *FileSysOut) closeFile() {
	l.Lock()
	defer l.Unlock()

	if l.file != nil {
		_ = l.file.Close()
	}
}

func (l *FileSysOut) Close() {
	l.Lock()
	defer l.Unlock()

	if l.ticker != nil {
		l.ticker.Stop()
	}

	l.fflush()
	l.closeFile()

	close(l.Bucket)
	l.close <- struct{}{}
}
