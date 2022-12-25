// Message 
var Message;

// Main object
var GomliSearch = {
    RawTargetCalls: [],     // This properity can be removed, just to ease the process of copy and pasting instructions
                            // ie. "La/b;->a(Ljava/lang/String;)Ljava/lang/String;",
    Calls: []               // Calls is mandatory and must be an array of Call{Class:"",Method:""}

};

// GetGomli returns a JSON string of the GomliSearch object
function GetGomli() {
    GomliSearch.RawTargetCalls.forEach((item) => {  // Populate Calls array
        callObj = {
            Class: item.split("->", -1)[0],
            Method: item.split("->", -1)[1],
        }
        GomliSearch.Calls.push(callObj)
    })

    return JSON.stringify(GomliSearch)
}

// Compare will return true if the current instruction requires to be printed or transformed
function Compare() {
    var parsedJSON = JSON.parse(base64.decode(Message))

    for (var i = 0; i < GomliSearch.Calls.length; i++) {
        if (parsedJSON.OpCode == 0x71 &&
            parsedJSON.Verbs[3] == GomliSearch.Calls[i].Class &&
            parsedJSON.Verbs[4] == GomliSearch.Calls[i].Method) {
            return true
        }
    }
    return false
}

// Actual transform function
function Transform() {
    var payload = base64.decode(Message) // Base64 encoded payload
    
    return "decoded"
}
