package models

type Device struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    IP     string `json:"ip"`
    Type   string `json:"type"`
    Port   int    `json:"port"`
    Status string `json:"status"`
}