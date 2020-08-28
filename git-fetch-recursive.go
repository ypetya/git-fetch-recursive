package main

import (
"os"
"os/exec"
"sync"
"fmt"
w "./walkdir"
)

func main() {
  if len(os.Args) == 2 {
    // walk -> git fetch
    ch, done := make(chan string), make(chan bool)
    go collect(ch,done)

    w.WalkDir(os.Args[1],ch, exclude)
    close(ch)
    <-done
  } else {
    //help()
  }
}

func exclude(relative string,absolute string) bool {
  if relative[0] == '.' || relative == "node_modules" ||
    relative == "build" || relative == "dist" || relative == "src" ||
    relative == "test" || relative == "Pictures" {
    return true
  }
  return false
}


func collect(dir <-chan string, done chan<- bool) {
  var count int
  waitAll := sync.WaitGroup{}

  for d := range dir {
    f, err := os.Open(d+"/.git")
    if err == nil {
      f.Close()
      waitAll.Add(1)
      go runGitFetch(d, &waitAll)
      count++
    }
  }
  waitAll.Wait()
  fmt.Println("\nTotal git repositories:", count)
  done<-true
}

func runGitFetch(dir string, w *sync.WaitGroup) {

  cmd:=exec.Command("git","fetch","--all")
  cmd.Dir = dir
  if cmd.Run() == nil {
    fmt.Print("+")
  } else {
    fmt.Print("-")
  }
  (*w).Done()
}
