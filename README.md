# [Site](http://cw.krasovsky.me/)
![cw.krasovsky.me](https://i.imgur.com/lDsJw1Y.gif)

# CW Craft API 0.2 documentation

### items
* `/api/items`
* `/api/items/:id` ([example](https://cw.krasovsky.me/items/a32))
* `/api/items?name={url_encoded_name}`
* `/api/items?type={url_encoded_type}`
Returns object or array of objects with such structure:
```golang
ID       string
Name     string
Stats    stats
Type     string
ManaCost int
Recipe   map[string]int
```
Where `stats` is:
```golang
Attack  int
Defense int
Mana    int
```

### resources
* `/api/resources`
* `/api/resources/:id` ([example](https://cw.krasovsky.me/resources/25))
* `/api/resources?name={url_encoded_name}`
Returns object or array of objects with such structure:
```golang
item     item
Commands []command
```

### basics
* `/api/basics/:id` ([example](https://cw.krasovsky.me/basics/a32))
Returns array of objects with such structure:
```golang
Item   item
Basics []basic
```
Where `basic` is:
```golang
Name   string
Amount int
```

### commands
* `/api/commands/:id` ([example](https://cw.krasovsky.me/commands/a32))
Returns object with such structure:
```golang
Item          item
Commands      []command
TotalManaCost int
```
Where `command` is:
```golang
ID              string
Name            string
Amount          int
CommandManaCost int
```

## Feel free to contribute!
