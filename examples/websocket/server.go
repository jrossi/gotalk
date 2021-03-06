package main
import (
  "fmt"
  "net/http"
  "github.com/rsms/gotalk"
)

type GreetIn struct {
  Name string `json:"name"`
}
type GreetOut struct {
  Greeting string `json:"greeting"`
}

func onAccept(s gotalk.Sock) {
  s.Notify("hello", "world")
  go func(){
    // Send a request & read result via JSON-encoded go values.
    greeting := GreetOut{}
    if err := s.Request("greet", GreetIn{"Rasmus"}, &greeting); err != nil {
      panic(err.Error())
    }
    fmt.Printf("greet: %+v\n", greeting)
  }()
}

func main() {
  gotalk.Handle("greet", func(in GreetIn) (GreetOut, error) {
    println("in greet handler: in.Name=", in.Name)
    return GreetOut{"Hello " + in.Name}, nil
  })

  gotalk.HandleBufferRequest("echó", func(s gotalk.Sock, op string, in []byte) ([]byte, error) {
    return in, nil
  })

  http.Handle("/gotalk", gotalk.WebSocketHandler(nil, onAccept))
  http.Handle("/", http.FileServer(http.Dir(".")))
  err := http.ListenAndServe(":1234", nil)
  if err != nil {
    panic(err.Error())
  }
}
