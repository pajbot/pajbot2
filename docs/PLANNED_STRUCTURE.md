```
/
    cmd/
        bot/
        web/

    bot/
        filter.go - defines a filter interface
        bot.go - defines a connect/disconnect interface

    irc/
        connection.go

    web/
        api/
            client.go
            command.go
            phrase.go

    pkg/
        phrases/
```

```go
// How modules could look like
package bot

// this is the generic "bot module"
type Module interface {
    Register() error
    OnMessage(user, message string) error
}

// maybe the bot is made up of two basic IO operations
type Bot interface {
    Sender
    Receiver

    Modules() []Module
}

type Sender interface {
    Send(string) error
}

type Receiver interface {
    Messages() chan<- string
}


// these would be in a "modules" folder? or a "chat"
type Chat struct {
    phrases []phrases
    db *sql.DB
    sender Sender
}

func New(db *sql.DB, sender Sender) Chat {
    return Chat{
        db: db,
        sender: sender,
    }
}

// this is where you load them all up?
func (m Chat) Register() error {
    return nil
}

func (m Chat) OnMessage(user, message string) error {
    m.sender.Send("some ban message")
    return nil
}
```
