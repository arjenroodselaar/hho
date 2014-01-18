
package main

import "fmt"

func a(x ...interface{}) []string {
  return []string{"howdy"}
}

func main() {
  fmt.Println(a("boners"))
}
