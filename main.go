package main

import (
  "bytes"
  "flag"
  "fmt"
  "os"
  "path"
  "strconv"
  "strings"
  "time"
)

var srcDir, dstDir string
var miteruFile string
var interval int64
const gap = 60 * 5

func init() {
  flag.StringVar(&srcDir, "src", "", "src")
  flag.StringVar(&dstDir, "dst", "", "dst")
  flag.StringVar(&miteruFile, "file", "miteru-kun", "koko ni nanika kaku")
  flag.Int64Var(&interval, "interval", 60*60*24, "kore yori chiisai")

  flag.Parse()
}

func isDirectory(path string) bool {
  info, err := os.Stat(path)
  if os.IsNotExist(err) {
    return false
  }
  mode := info.Mode()
  return mode.IsDir()
}

func readLast(path string) int64 {
  file, err := os.Open(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
    os.Exit(-1)
  }
  defer file.Close()

  buf := make([]byte, 16) // magic
  _, err = file.Read(buf)
  if err != nil {
    fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
    os.Exit(-1)
  }
  last, err := strconv.Atoi(strings.Split(bytes.NewBuffer(buf).String(), "\n")[0])
  if err != nil {
    fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
    os.Exit(-1)
  }
  return int64(last)
}

func main() {
  if !isDirectory(srcDir) {
    fmt.Fprintf(os.Stderr, "not directory\nsrc: %v\n", srcDir)
    os.Exit(-1)
  }
  if !isDirectory(dstDir) {
    fmt.Fprintf(os.Stderr, "not directory\ndst: %v\n", dstDir)
    os.Exit(-1)
  }

  dst := path.Join(dstDir, miteruFile)
  if _, err := os.Stat(dst); os.IsNotExist(err) {
    fmt.Fprintf(os.Stderr, "%v not found(maybe first time?).\n", dst)
    os.Exit(-1)
  }

  srcLast := readLast(path.Join(srcDir, miteruFile))
  dstLast := readLast(path.Join(dstDir, miteruFile))

  now := time.Now().Unix()
  if now - srcLast > interval - gap {
    fmt.Fprintf(os.Stderr, "miteru-kun seems to fail for over %v(seconds)\nlast: %v\n", interval, srcLast)
    os.Exit(-1)
  }

  if srcLast - dstLast > interval + gap {
    fmt.Fprintf(os.Stderr, "backup seems to fail for over %v(seconds)\nsrcLast: %v\ndstLast: %v\n", interval, srcLast, dstLast)
    os.Exit(-1)
  }

  fmt.Printf("seems good. update srcLast to %v(from %v)\n", now, srcLast)
  file, err := os.OpenFile(path.Join(srcDir, miteruFile), os.O_RDWR, 0644)
  if err != nil {
    fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
    os.Exit(-1)
  }
  defer file.Close()
  _, err = file.Write([]byte(strconv.FormatInt(now, 10)))
  if err != nil {
    fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
    os.Exit(-1)
  }

  fmt.Println("success")
}
