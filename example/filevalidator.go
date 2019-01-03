package main

import (
  "time"
)

func FileValidator(key string, otime time.Time) bool {

  fi, err := os.Stat(key)
  if err != nil {
    // Does not exists anymore, invalid
    return false
  }
  mtime := fi.ModTime()
  if mtime.After(otime) {
    // file is newer, invalid
    return false
  }
  // All ok, valid
  return true
}