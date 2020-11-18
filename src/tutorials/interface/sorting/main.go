package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// 自定义排序接口要实现 Len, Less, Swap
type StringSlice []string

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// 这里对每个元素用指针，是因为性能更好
// 因为每个元素频繁比较，可以更省内存和CPU数据复制的时间
var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

// 将结果打印成表格
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "-----", "-----", "-----", "-----")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}

type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Album < x[j].Album }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	sort.Sort(byArtist(tracks))
	// 反序排序，不需要新建一个 byReverseArtist 的新类型，直接调用 reverse 即可
	sort.Sort(sort.Reverse(byArtist(tracks)))
	// 如果要以其他列排序，那么就得新加 sort 类，比如 type byTitle []*Track
	printTracks(tracks)
}
