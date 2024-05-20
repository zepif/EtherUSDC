package main

import (
    "fmt"
    "reflect"

    "github.com/zepif/EtherUSDC/internal/data"
)

func main() {
    var tx data.Transaction
    t := reflect.TypeOf(tx)

    fmt.Printf("Структура %s имеет следующие поля:\n", t.Name())

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("- %s (%s)\n", field.Name, field.Type.Name())
    }
}
