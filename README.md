# payloadGO
an library to add payload to file

## Usage

```
import (
  "github.com/thomas10-10/payloadGo"
)
```

### payloadGo.PUT(path string,payload string)
Adds a payload on the file given in path, if the path is "" then the payload will be added on the binary that has been launched. If a payload already exists then it will be replaced

### payloadGo.GET(path string).DATA
If the path is "" then return the payload of the binary that has been launched

### payloadGo.DELETE(path string)
delete the payload of the file given in path, if the path is "" then delete the payload of the binary that has been launched

