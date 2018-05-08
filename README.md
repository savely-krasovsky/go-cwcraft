# [CW Craft Site 1.4](http://cw.krasovsky.me/)
![cw.krasovsky.me](https://i.imgur.com/Q6pUWhw.gif)

# CW Craft API 1.4 documentation

### Equipment
* `/api/equipment`
* `/api/equipment/:id` ([example](https://cw.krasovsky.me/api/equipment/a32))
* `/api/equipment?name={url_encoded_name}`
* `/api/equipment?type={url_encoded_type}`

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

### Alchemy
* `/api/alchemy`
* `/api/alchemy/:id` ([example](https://cw.krasovsky.me/api/alchemy/p03))
* `/api/alchemy?name={url_encoded_name}`
* `/api/alchemy?type={url_encoded_type}`

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
* `/api/resources/:id` ([example](https://cw.krasovsky.me/api/resources/25))
* `/api/resources?name={url_encoded_name}`

Returns object or array of objects with such structure:
```golang
ID        string
Name      string
ManaCost  int
Composite bool
Recipe    map[string]int
```

### basics
* `/api/basics/:type/:id` ([example](https://cw.krasovsky.me/api/basics/equipment/a32))

Where `type` could be:
* equipment
* alchemy

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
* `/api/commands/:type/:id` ([example](https://cw.krasovsky.me/api/commands/equipment/a32))

Where `type` could be:
* equipment
* alchemy

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

### shops
* `/api/shops`

Returns object with such structure:
```golang
Link        string
Name        string
OwnerName   string
OwnerCastle string
Kind        string
Mana        int
Offers      []OfferItem
```
Where `OfferItem` is:
```golang
Item  string
Price int
Mana  int
```

## Feel free to contribute!
